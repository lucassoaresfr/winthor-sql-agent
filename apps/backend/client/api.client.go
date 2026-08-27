package client

import (
	"fmt"

	"io"

	"net/http"

	"time"
)

type APIClient struct {
	BaseURL string

	AuthToken string

	HTTPClient *http.Client
}

func NewAPIClient(baseURL string, authToken string) *APIClient {

	return &APIClient{

		BaseURL: baseURL,

		AuthToken: authToken,

		HTTPClient: &http.Client{

			Timeout: 10 * time.Second,
		},
	}

}

// ChamarRotaGet executa a requisição HTTPS para a API principal e retorna o JSON puro em bytes

func (c *APIClient) ChamarRotaGet(endpoint string) ([]byte, error) {

	url := fmt.Sprintf("%s%s", c.BaseURL, endpoint)

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {

		return nil, err

	}

	req.Header.Set("Authorization", "Bearer "+c.AuthToken)

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)

	if err != nil {

		return nil, fmt.Errorf("falha ao conectar na API: %w", err)

	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {

		return nil, fmt.Errorf("API retornou status HTTP %d", resp.StatusCode)

	}

	return io.ReadAll(resp.Body)

}
