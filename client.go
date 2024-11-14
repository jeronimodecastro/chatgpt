// client.go
package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

type ClientOption func(*Client)

// Configurações do cliente
func WithBaseURL(url string) ClientOption {
	return func(c *Client) {
		c.baseURL = url
	}
}

func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.httpClient.Timeout = timeout
	}
}

func LoadEnv() error {
	// Carrega o arquivo .env se existir
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("erro ao carregar .env: %w", err)
	}
	return nil
}

func GetOpenAIKey() (string, error) {
	key := os.Getenv("OPENAI_API_KEY")
	if key == "" {
		return "", fmt.Errorf("OPENAI_API_KEY não está definida")
	}
	return key, nil
}

// NewClient cria uma nova instância do cliente OpenAI
func NewClient(apiKey string, opts ...ClientOption) (*Client, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("API key não pode estar vazia")
	}

	client := &Client{
		apiKey:  apiKey,
		baseURL: "https://api.openai.com/v1",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	// Aplica as opções customizadas
	for _, opt := range opts {
		opt(client)
	}

	return client, nil
}

// types.go
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	Temperature float64       `json:"temperature,omitempty"`
}

type ChatCompletionResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Choices []struct {
		Message ChatMessage `json:"message"`
	} `json:"choices"`
}

// chat.go
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

	req, err := c.newRequest("POST", "/v1/chat/completions", bytes.NewBuffer(jsonData))
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

// request.go
func (c *Client) newRequest(method, path string, body *bytes.Buffer) (*http.Request, error) {
	req, err := http.NewRequest(method, c.baseURL+path, body)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	return req, nil
}

func (c *Client) sendRequest(req *http.Request, v interface{}) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("erro ao enviar requisição: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API retornou status code inesperado: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		return fmt.Errorf("erro ao decodificar resposta: %w", err)
	}

	return nil
}
