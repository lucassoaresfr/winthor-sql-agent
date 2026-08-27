package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"
)

type BrasilAPIClient struct {
	httpClient *http.Client
}

func NewBrasilAPIClient() *BrasilAPIClient {
	return &BrasilAPIClient{
		httpClient: &http.Client{
			Timeout: 8 * time.Second,
		},
	}
}

// ConsultarCNPJ realiza a chamada HTTP REST para a BrasilAPI v1
func (c *BrasilAPIClient) ConsultarCNPJ(ctx context.Context, cnpj string) ([]byte, error) {
	// Remove caracteres não numéricos (. / -)
	re := regexp.MustCompile(`\D`)
	cnpjLimpo := re.ReplaceAllString(cnpj, "")

	if len(cnpjLimpo) != 14 {
		return nil, fmt.Errorf("CNPJ inválido: deve conter 14 dígitos numéricos")
	}

	url := fmt.Sprintf("https://brasilapi.com.br/api/cnpj/v1/%s", cnpjLimpo)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("falha ao criar requisição: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro na chamada à BrasilAPI: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("falha ao ler resposta da BrasilAPI: %w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("CNPJ %s não foi encontrado na base da Receita Federal", cnpj)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("BrasilAPI retornou status %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}
