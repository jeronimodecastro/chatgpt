# OpenAI Client

A Go client for interacting with OpenAI's API with secure credential management.

## Prerequisites

- Go 1.21 or higher
- OpenAI API key
- (Optional) Google Cloud account for Secret Manager

## Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/openai-client
cd openai-client

# Install dependencies
go mod download
go get github.com/joho/godotenv
go get cloud.google.com/go/secretmanager/apiv1/secretmanagerpb@v1.14.2
go get gopkg.in/yaml.v2
go mod tidy
```

## Configuration

1. Create `.env` file:
```bash
cp .env.example .env
```

2. Add your OpenAI API key to `.env`:
```
OPENAI_API_KEY=your_api_key_here
GO_ENV=development
```

## Project Structure

```
├── cmd
│   └── main.go
├── internal
│   ├── config
│   │   ├── config.go
│   │   └── secret_manager.go
│   └── openai
│       ├── client.go
│       ├── types.go
│       └── request.go
├── .env
├── .env.example
├── .gitignore
└── go.mod
```

## Usage

Run the application:
```bash
go run cmd/main.go
```

Example code:
```go
client, err := openai.NewClient(
    apiKey,
    openai.WithTimeout(20 * time.Second),
)

response, err := client.CreateChatCompletion("Hello, how are you?")
```

## Production Setup

For production environments:

1. Set up Google Cloud Secret Manager
2. Update environment variable:
```
GO_ENV=production
```

## Error Handling

The client includes comprehensive error handling for:
- API key validation
- HTTP request failures
- Response parsing
- Empty responses

## Contributing

1. Fork the repository
2. Create your feature branch
3. Submit a pull request

## Security

- Never commit `.env` file
- Use Secret Manager for production
- Rotate API keys regularly
- Monitor API usage

## License

MIT License