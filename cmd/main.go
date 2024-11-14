package main

import (
	"chatgpt/internal/config"
	"chatgpt/internal/openai"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "development"
	}

	var apiKey string
	var err error

	if env == "development" {
		if err := config.LoadEnv(); err != nil {
			log.Fatal(err)
		}

		apiKey, err = config.GetOpenAIKey()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		sm, err := config.NewSecretManager("seu-projeto-id")
		if err != nil {
			log.Fatal(err)
		}
		defer sm.Close()

		apiKey, err = sm.GetOpenAIKey()
		if err != nil {
			log.Fatal(err)
		}
	}

	client, err := openai.NewClient(
		apiKey,
		openai.WithTimeout(20*time.Second),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Pergunta fora do contexto.
	// response, err := client.CreateChatCompletion("Olá, como você está?")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// Pergunta com contexto.
	analyzer := openai.NewURLAnalyzer(client, "urls.txt")
	response, err := analyzer.Analyze("O que pode ser Alíquota de ICMS superior a definida para a operação interestadual ?")

	fmt.Println("Resposta:", response)
}
