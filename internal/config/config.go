package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() error {
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
