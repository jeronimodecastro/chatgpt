package openai

import (
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

type ClientOption func(*Client)

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

	for _, opt := range opts {
		opt(client)
	}

	return client, nil
}
