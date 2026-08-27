package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/lucassoaresfr/winthor-sql-agent.git/client"
	"github.com/lucassoaresfr/winthor-sql-agent.git/config"
	"github.com/lucassoaresfr/winthor-sql-agent.git/tools"
	"google.golang.org/genai"
)

var ModelosFallback = []string{
	"gemini-3.7-flash",
	"gemini-3.6-flash",
	"gemini-3.5-flash",
	"gemini-3.5-flash-lite",
	"gemini-3.1-pro-preview",
	"gemini-3.1-flash-lite",
	"gemini-3-flash-preview",
	"gemini-2.5-pro",
	"gemini-2.5-flash",
	"gemini-2.5-flash-lite",
	"gemini-pro-latest",
	"gemini-flash-latest",
	"gemini-flash-lite-latest",
	"gemma-4-31b-it",
	"gemma-4-26b-a4b-it",
	"gemini-omni-1.1-flash",
}

type ChatMessage struct {
	Role    string `json:"role"`    // "user" ou "assistant"
	Content string `json:"content"` // Conteúdo da mensagem
}

type Orchestrator struct {
	GeminiClient    *genai.Client
	APIClient       *client.APIClient
	BrasilAPIClient *client.BrasilAPIClient
	Config          *config.Config
}

func NewOrchestrator(geminiClient *genai.Client, apiClient *client.APIClient, BrasilAPIClient *client.BrasilAPIClient, cfg *config.Config) *Orchestrator {
	return &Orchestrator{
		GeminiClient:    geminiClient,
		APIClient:       apiClient,
		BrasilAPIClient: BrasilAPIClient,
		Config:          cfg,
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
				genai.NewPartFromText(SystemPrompt),
			},
		},
		Tools: []*genai.Tool{
			{
				FunctionDeclarations: []*genai.FunctionDeclaration{
					tools.ToolConsultarClientes,
					tools.ToolConsultarProdutos,
					tools.ToolConsultarCNPJExterno,
				},
			},
		},
	}

	contents := converterHistoricoParaGenAI(historico)

	var resp *genai.GenerateContentResponse
	var err error
	var modeloUsado string

	// 1. Primeira Chamada (Seleção de Tool) - Timeout ajustado para 20s
	for _, modelName := range ModelosFallback {
		reqCtx, reqCancel := context.WithTimeout(baseCtx, 20*time.Second)
		resp, err = o.GeminiClient.Models.GenerateContent(reqCtx, modelName, contents, genConfig)
		reqCancel()

		if err == nil {
			modeloUsado = modelName
			break
		}

		log.Printf("[Orchestrator] Falha/Timeout no modelo %s: %v. Tentando próximo...", modelName, err)
		time.Sleep(300 * time.Millisecond)
	}

	if err != nil || resp == nil {
		return "", fmt.Errorf("todos os modelos falharam na requisição inicial: %w", err)
	}

	// 2. Processamento de Function Calls
	for _, candidate := range resp.Candidates {
		if candidate.Content == nil {
			continue
		}

		for _, part := range candidate.Content.Parts {
			if part.FunctionCall != nil {
				fnCall := part.FunctionCall
				var jsonBytes []byte

				switch fnCall.Name {
				case "consultar_clientes_api":
					jsonBytes, err = o.executarChamadaAPIClientes(fnCall.Args)
				case "consultar_produtos_api":
					jsonBytes, err = o.executarChamadaAPIProdutos(fnCall.Args)
				case "consultar_cnpj_externo":
					cnpj, _ := fnCall.Args["cnpj"].(string)
					jsonBytes, err = o.BrasilAPIClient.ConsultarCNPJ(baseCtx, cnpj)
				default:
					return "", fmt.Errorf("ferramenta desconhecida chamada: %s", fnCall.Name)
				}

				if err != nil {
					return "", fmt.Errorf("erro ao executar a ferramenta %s: %w", fnCall.Name, err)
				}

				var apiResponse interface{}
				if err := json.Unmarshal(jsonBytes, &apiResponse); err != nil {
					apiResponse = map[string]interface{}{"raw": string(jsonBytes)}
				}

				promptComHistorico := append(contents,
					candidate.Content,
					&genai.Content{
						Role: "user",
						Parts: []*genai.Part{
							genai.NewPartFromFunctionResponse(fnCall.Name, map[string]interface{}{
								"dados": apiResponse,
							}),
						},
					},
				)

				// 3. Segunda Chamada (Síntese final)
				modelosSintese := []string{modeloUsado, "gemini-2.5-flash-lite"}
				var finalResp *genai.GenerateContentResponse

				for _, modelName := range modelosSintese {
					sinteseCtx, sinteseCancel := context.WithTimeout(baseCtx, 15*time.Second)
					finalResp, err = o.GeminiClient.Models.GenerateContent(sinteseCtx, modelName, promptComHistorico, genConfig)
					sinteseCancel()

					if err == nil {
						txt := extrairTextoResposta(finalResp)
						if txt != "Não foi possível gerar uma resposta legível" {
							return txt, nil
						}
					}
					log.Printf("[Orchestrator] Falha na síntese com %s: %v", modelName, err)
				}

				return "", fmt.Errorf("falha ao sintetizar resposta final: %w", err)
			}
		}
	}

	return extrairTextoResposta(resp), nil
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
	baseURL, err := url.Parse("/api/v1/client")
	if err != nil {
		return nil, err
	}

	queryParams := url.Values{}
	for chave, valor := range args {
		if valor != nil {
			valorStr := fmt.Sprintf("%v", valor)
			if valorStr != "" {
				queryParams.Add(chave, valorStr)
			}
		}
	}

	baseURL.RawQuery = queryParams.Encode()
	return o.APIClient.ChamarRotaGet(baseURL.String())
}

func (o *Orchestrator) executarChamadaAPIProdutos(args map[string]interface{}) ([]byte, error) {
	baseURL, err := url.Parse("/api/v1/prod")
	if err != nil {
		return nil, err
	}

	queryParams := url.Values{}
	for chave, valor := range args {
		if valor != nil {
			valorStr := fmt.Sprintf("%v", valor)
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
