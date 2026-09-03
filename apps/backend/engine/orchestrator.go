package engine

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/lucassoaresfr/winthor-sql-agent.git/client"
	"github.com/lucassoaresfr/winthor-sql-agent.git/config"
	"github.com/lucassoaresfr/winthor-sql-agent.git/tools"
	"google.golang.org/genai"
)

type ChatMessage struct {
	Role    string `json:"role"`    // "user" ou "assistant"
	Content string `json:"content"` // Conteúdo da mensagem
}

type Orchestrator struct {
	GeminiClient    *genai.Client
	APIClient       *client.APIClient
	BrasilAPIClient *client.BrasilAPIClient
	Config          *config.Config
	ModelMgr        *ModelManager // Injetado para gerenciar o Circuit Breaker de modelos no Redis
}

// ErrTodosModelosEmCooldown retornado quando nenhum modelo do fallback consegue responder
var ErrTodosModelosEmCooldown = errors.New("a API do Gemini está enfrentando alta demanda no momento. Por favor, tente novamente em alguns instantes")

func NewOrchestrator(
	geminiClient *genai.Client,
	apiClient *client.APIClient,
	brasilAPIClient *client.BrasilAPIClient,
	cfg *config.Config,
	modelMgr *ModelManager,
) *Orchestrator {
	return &Orchestrator{
		GeminiClient:    geminiClient,
		APIClient:       apiClient,
		BrasilAPIClient: brasilAPIClient,
		Config:          cfg,
		ModelMgr:        modelMgr,
	}
}

func (o *Orchestrator) ProcessarPergunta(c context.Context, historico []ChatMessage) (string, error) {
	if len(historico) == 0 {
		return "", fmt.Errorf("histórico de mensagens não pode estar vazio")
	}

	baseCtx := context.WithoutCancel(c)

	genConfig := &genai.GenerateContentConfig{
		SystemInstruction: &genai.Content{
			Parts: []*genai.Part{
				genai.NewPartFromText(GetSystemPrompt()),
			},
		},
		Tools: []*genai.Tool{
			{
				FunctionDeclarations: []*genai.FunctionDeclaration{
					tools.ToolConsultarClientes,
					tools.ToolConsultarProdutos,
					tools.ToolConsultarPromocoes,
					tools.ToolConsultarPedidos,
					tools.ToolConsultarItensPedido,
					tools.ToolConsultarCNPJExterno,
				},
			},
		},
	}

	contents := converterHistoricoParaGenAI(historico)

	// Permite até 5 rodadas de chamadas de função em sequência
	maxIteracoes := 5

	for iter := 0; iter < maxIteracoes; iter++ {
		var resp *genai.GenerateContentResponse
		var err error
		var modeloUsado string

		// 1. Chamada ao Gemini com Fallback entre modelos
		for _, modelName := range ModelosFallback {
			if o.ModelMgr.EstaEmCooldown(baseCtx, modelName) {
				log.Printf("[Orchestrator] Modelo %s em cooldown no Redis. Pulando...", modelName)
				continue
			}

			// Timeout expandido para suportar retornos complexos de Function Calling
			reqCtx, reqCancel := context.WithTimeout(baseCtx, 45*time.Second)
			resp, err = o.GeminiClient.Models.GenerateContent(reqCtx, modelName, contents, genConfig)
			reqCancel()

			if err == nil && resp != nil {
				modeloUsado = modelName
				_ = modeloUsado
				break
			}

			// Se houver falha de quota, timeout ou instabilidade (503)
			if ehErroDeQuotaOuInstabilidade(err) {
				log.Printf("[Orchestrator] Quota/Instabilidade no modelo %s: %v. Registrando cooldown...", modelName, err)
				_ = o.ModelMgr.RegistrarCooldownBloqueio(baseCtx, modelName, 2) // Cooldown de 2 min
			} else {
				log.Printf("[Orchestrator] Erro no modelo %s: %v", modelName, err)
			}

			time.Sleep(300 * time.Millisecond)
		}

		// Se após percorrer todos os modelos do fallback nenhum responder
		if resp == nil {
			log.Printf("[Orchestrator] Todos os modelos do fallback falharam na iteração %d. Último erro: %v", iter, err)
			return "", ErrTodosModelosEmCooldown
		}

		// Procura por chamadas de função na resposta do candidato
		var functionCalls []*genai.FunctionCall
		var modelContent *genai.Content

		for _, candidate := range resp.Candidates {
			if candidate.Content == nil {
				continue
			}
			modelContent = candidate.Content
			for _, part := range candidate.Content.Parts {
				if part.FunctionCall != nil {
					functionCalls = append(functionCalls, part.FunctionCall)
				}
			}
		}

		// Caso NÃO HAJA NENHUMA FunctionCall, o modelo devolveu a resposta final
		if len(functionCalls) == 0 {
			return extrairTextoResposta(resp), nil
		}

		// Adiciona a intenção/resposta do modelo ao histórico da conversa
		contents = append(contents, modelContent)

		// 2. Executa as chamadas de função retornadas
		var responseParts []*genai.Part

		for _, fnCall := range functionCalls {
			var jsonBytes []byte

			switch fnCall.Name {
			case "consultar_clientes_api":
				jsonBytes, err = o.executarChamadaAPIClientes(fnCall.Args)
			case "consultar_produtos_api":
				jsonBytes, err = o.executarChamadaAPIProdutos(fnCall.Args)
			case "consultar_promocoes_api":
				jsonBytes, err = o.executarChamadaAPIPromocoes(fnCall.Args)
			case "consultar_pedidos_api":
				jsonBytes, err = o.executarChamadaAPIPedidos(fnCall.Args)
			case "consultar_itens_pedido_api":
				jsonBytes, err = o.executarChamadaAPIItensPedido(fnCall.Args)
			case "consultar_cnpj_externo":
				cnpj, _ := fnCall.Args["cnpj"].(string)
				jsonBytes, err = o.BrasilAPIClient.ConsultarCNPJ(baseCtx, cnpj)
			default:
				err = fmt.Errorf("ferramenta desconhecida chamada: %s", fnCall.Name)
			}

			if err != nil {
				log.Printf("[Orchestrator] Erro ao executar a ferramenta %s: %v", fnCall.Name, err)
				jsonBytes = []byte(fmt.Sprintf(`{"error": "%s"}`, err.Error()))
			}

			var apiResponse interface{}
			if err := json.Unmarshal(jsonBytes, &apiResponse); err != nil {
				apiResponse = map[string]interface{}{"raw": string(jsonBytes)}
			}

			// Prepara a resposta da função
			partResp := genai.NewPartFromFunctionResponse(fnCall.Name, map[string]interface{}{
				"dados": apiResponse,
			})
			responseParts = append(responseParts, partResp)
		}

		// Adiciona o bloco com as respostas das funções ao histórico para a próxima iteração
		contents = append(contents, &genai.Content{
			Role:  "user",
			Parts: responseParts,
		})
	}

	return "", fmt.Errorf("excedido o limite máximo de iterações de chamadas de função")
}

// Auxiliar para detectar erros de quota, rate limit, timeout (504) ou instabilidade (503)
func ehErroDeQuotaOuInstabilidade(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "429") ||
		strings.Contains(msg, "resourceexhausted") ||
		strings.Contains(msg, "quota") ||
		strings.Contains(msg, "rate limit") ||
		strings.Contains(msg, "503") ||
		strings.Contains(msg, "504") ||
		strings.Contains(msg, "deadline exceeded") ||
		strings.Contains(msg, "unavailable") ||
		strings.Contains(msg, "high demand")
}

func converterHistoricoParaGenAI(historico []ChatMessage) []*genai.Content {
	contents := make([]*genai.Content, 0, len(historico))

	for _, msg := range historico {
		role := msg.Role
		if role == "assistant" {
			role = "model"
		}

		contents = append(contents, &genai.Content{
			Role: role,
			Parts: []*genai.Part{
				genai.NewPartFromText(msg.Content),
			},
		})
	}

	return contents
}

func (o *Orchestrator) executarChamadaAPIClientes(args map[string]interface{}) ([]byte, error) {
	baseURL, err := url.Parse("/api/v1/tools/client")
	if err != nil {
		return nil, err
	}

	queryParams := url.Values{}
	for chave, valor := range args {
		if valor != nil {
			valorStr := formatarValorParametro(valor)
			if valorStr != "" {
				queryParams.Add(chave, valorStr)
			}
		}
	}

	baseURL.RawQuery = queryParams.Encode()
	return o.APIClient.ChamarRotaGet(baseURL.String())
}

func (o *Orchestrator) executarChamadaAPIProdutos(args map[string]interface{}) ([]byte, error) {
	baseURL, err := url.Parse("/api/v1/tools/prod")
	if err != nil {
		return nil, err
	}

	queryParams := url.Values{}
	for chave, valor := range args {
		if valor != nil {
			valorStr := formatarValorParametro(valor)
			if valorStr != "" {
				queryParams.Add(chave, valorStr)
			}
		}
	}

	baseURL.RawQuery = queryParams.Encode()
	return o.APIClient.ChamarRotaGet(baseURL.String())
}

func (o *Orchestrator) executarChamadaAPIPromocoes(args map[string]interface{}) ([]byte, error) {
	baseURL, err := url.Parse("/api/v1/tools/promotion")
	if err != nil {
		return nil, err
	}

	queryParams := url.Values{}
	for chave, valor := range args {
		if valor != nil {
			valorStr := formatarValorParametro(valor)
			if valorStr != "" {
				queryParams.Add(chave, valorStr)
			}
		}
	}

	baseURL.RawQuery = queryParams.Encode()
	return o.APIClient.ChamarRotaGet(baseURL.String())
}

func (o *Orchestrator) executarChamadaAPIPedidos(args map[string]interface{}) ([]byte, error) {
	baseURL, err := url.Parse("/api/v1/tools/orders")
	if err != nil {
		return nil, err
	}

	queryParams := url.Values{}
	for chave, valor := range args {
		if valor != nil {
			valorStr := formatarValorParametro(valor)
			if valorStr != "" {
				queryParams.Add(chave, valorStr)
			}
		}
	}

	baseURL.RawQuery = queryParams.Encode()
	return o.APIClient.ChamarRotaGet(baseURL.String())
}

func (o *Orchestrator) executarChamadaAPIItensPedido(args map[string]interface{}) ([]byte, error) {
	baseURL, err := url.Parse("/api/v1/tools/items")
	if err != nil {
		return nil, err
	}

	queryParams := url.Values{}
	for chave, valor := range args {
		if valor != nil {
			valorStr := formatarValorParametro(valor)
			if valorStr != "" {
				queryParams.Add(chave, valorStr)
			}
		}
	}

	baseURL.RawQuery = queryParams.Encode()
	return o.APIClient.ChamarRotaGet(baseURL.String())
}

func extrairTextoResposta(resp *genai.GenerateContentResponse) string {
	if resp == nil {
		return "Não foi possível gerar uma resposta legível"
	}

	for _, candidate := range resp.Candidates {
		if candidate.Content == nil {
			continue
		}
		for _, part := range candidate.Content.Parts {
			if part.Text != "" {
				return part.Text
			}
		}
	}

	return "Não foi possível gerar uma resposta legível"
}

func formatarValorParametro(valor interface{}) string {
	if valor == nil {
		return ""
	}
	switch v := valor.(type) {
	case float64:
		if v == float64(int64(v)) {
			return fmt.Sprintf("%d", int64(v))
		}
		return fmt.Sprintf("%.2f", v)
	case float32:
		if v == float32(int32(v)) {
			return fmt.Sprintf("%d", int32(v))
		}
		return fmt.Sprintf("%.2f", v)
	default:
		return fmt.Sprintf("%v", valor)
	}
}
