package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) CreateChatCompletion(question string) (string, error) {
	requestBody := ChatCompletionRequest{
		Model: "gpt-3.5-turbo",
		Messages: []ChatMessage{
			{
				Role:    "user",
				Content: question,
			},
		},
		Temperature: 0.7,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("erro ao serializar requisição: %w", err)
	}

	req, err := c.newRequest("POST", "/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	var response ChatCompletionResponse
	if err := c.sendRequest(req, &response); err != nil {
		return "", err
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("nenhuma resposta gerada pelo modelo")
	}

	return response.Choices[0].Message.Content, nil
}

func (c *Client) newRequest(method, path string, body *bytes.Buffer) (*http.Request, error) {
	req, err := http.NewRequest(method, c.baseURL+path, body)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	return req, nil
}

// func (c *Client) sendRequest(req *http.Request, v interface{}) error {
// 	resp, err := c.httpClient.Do(req)
// 	if err != nil {
// 		return fmt.Errorf("erro ao enviar requisição: %w", err)
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		return fmt.Errorf("API retornou status code inesperado: %d", resp.StatusCode)
// 	}

// 	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
// 		return fmt.Errorf("erro ao decodificar resposta: %w", err)
// 	}

// 	return nil
// }

type APIError struct {
	Code    int
	Message string
}

func (c *Client) sendRequest(req *http.Request, v interface{}) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("erro ao enviar requisição: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return json.NewDecoder(resp.Body).Decode(v)
	case http.StatusUnauthorized:
		var details struct {
			Error struct {
				Message string `json:"message"`
			} `json:"error"`
		}
		json.NewDecoder(resp.Body).Decode(&details)

		switch details.Error.Message {
		case "Invalid Authentication":
			return &APIError{401, "Chave API inválida ou organização incorreta"}
		case "Incorrect API key provided":
			return &APIError{401, "Chave API incorreta. Verifique ou gere uma nova"}
		default:
			return &APIError{401, "Você precisa ser membro de uma organização"}
		}

	case http.StatusForbidden:
		return &APIError{403, "País ou região não suportada"}

	case http.StatusTooManyRequests:
		var details struct {
			Error struct {
				Message string `json:"message"`
			} `json:"error"`
		}
		json.NewDecoder(resp.Body).Decode(&details)

		if strings.Contains(details.Error.Message, "quota") {
			return &APIError{429, "Cota excedida. Verifique seu plano"}
		}
		return &APIError{429, "Limite de requisições atingido. Aguarde"}

	case http.StatusInternalServerError:
		return &APIError{500, "Erro no servidor. Tente novamente em breve"}

	case http.StatusServiceUnavailable:
		return &APIError{503, "Servidor sobrecarregado. Tente novamente depois"}

	default:
		return fmt.Errorf("status code não esperado: %d", resp.StatusCode)
	}
}
