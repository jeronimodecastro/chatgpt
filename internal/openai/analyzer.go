package openai

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/net/html"
)

func NewURLAnalyzer(client *Client, urlsFile string) *URLAnalyzer {
	return &URLAnalyzer{
		Client:   client,
		URLsFile: urlsFile,
	}
}

// func (a *URLAnalyzer) Analyze(question string) (string, error) {
// 	file, err := os.Open(a.URLsFile)
// 	if err != nil {
// 		return "", fmt.Errorf("erro ao abrir arquivo: %w", err)
// 	}
// 	defer file.Close()

// 	var content strings.Builder
// 	var loadedURLs []string
// 	scanner := bufio.NewScanner(file)

// 	// Coleta todas as URLs
// 	for scanner.Scan() {
// 		url := scanner.Text()
// 		if url != "" {
// 			loadedURLs = append(loadedURLs, url)
// 		}
// 	}

// 	if len(loadedURLs) == 0 {
// 		return "", fmt.Errorf("nenhuma URL encontrada no arquivo %s", a.URLsFile)
// 	}

// 	fmt.Printf("Carregadas %d URLs do arquivo\n", len(loadedURLs))

// 	// Processa cada URL
// 	for i, url := range loadedURLs {
// 		fmt.Printf("Processando URL %d/%d: %s\n", i+1, len(loadedURLs), url)

// 		resp, err := http.Get(url)
// 		if err != nil {
// 			fmt.Printf("Erro ao acessar %s: %v\n", url, err)
// 			continue
// 		}
// 		defer resp.Body.Close()

// 		if resp.StatusCode != http.StatusOK {
// 			fmt.Printf("URL %s retornou status code: %d\n", url, resp.StatusCode)
// 			continue
// 		}

// 		body, err := ioutil.ReadAll(resp.Body)
// 		if err != nil {
// 			fmt.Printf("Erro ao ler conteúdo de %s: %v\n", url, err)
// 			continue
// 		}
// 		content.WriteString(string(body))
// 		content.WriteString("\n")
// 	}

// 	if content.Len() == 0 {
// 		return "", fmt.Errorf("não foi possível obter conteúdo de nenhuma URL")
// 	}

//		prompt := fmt.Sprintf("Contexto:\n%s\n\nPergunta: %s", content.String(), question)
//		return a.Client.CreateChatCompletion(prompt)
//	}
func extractText(node *html.Node) string {
	if node.Type == html.TextNode {
		return node.Data
	}
	var text string
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		text += extractText(c)
	}
	return text
}

func cleanText(input string) string {
	// Remove espaços extras
	text := strings.Join(strings.Fields(input), " ")
	// Limita o tamanho (por exemplo, 1000 caracteres por URL)
	if len(text) > 1000 {
		text = text[:1000]
	}
	return text
}

func (a *URLAnalyzer) Analyze(question string) (string, error) {
	file, err := os.Open(a.URLsFile)
	if err != nil {
		return "", fmt.Errorf("erro ao abrir arquivo: %w", err)
	}
	defer file.Close()

	var content strings.Builder
	var loadedURLs []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		url := scanner.Text()
		if url != "" {
			loadedURLs = append(loadedURLs, url)
		}
	}

	if len(loadedURLs) == 0 {
		return "", fmt.Errorf("nenhuma URL encontrada no arquivo %s", a.URLsFile)
	}

	fmt.Printf("Carregadas %d URLs do arquivo\n", len(loadedURLs))

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	for i, url := range loadedURLs {
		fmt.Printf("Processando URL %d/%d: %s\n", i+1, len(loadedURLs), url)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Printf("Erro ao criar requisição para %s: %v\n", url, err)
			continue
		}

		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
		req.Header.Set("Accept", "text/html")

		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Erro ao acessar %s: %v\n", url, err)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("URL %s retornou status code: %d\n", url, resp.StatusCode)
			resp.Body.Close()
			continue
		}

		doc, err := html.Parse(resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Printf("Erro ao parsear HTML de %s: %v\n", url, err)
			continue
		}

		// Extrai e limpa o texto
		text := cleanText(extractText(doc))
		content.WriteString(fmt.Sprintf("\nConteúdo resumido da URL %s:\n", url))
		content.WriteString(text)
		content.WriteString("\n---\n")
	}

	if content.Len() == 0 {
		return "", fmt.Errorf("não foi possível obter conteúdo de nenhuma URL")
	}

	prompt := fmt.Sprintf("Analise o seguinte conteúdo das URLs e responda: %s\n\nConteúdo:\n%s",
		question, content.String())

	return a.Client.CreateChatCompletion(prompt)
}
