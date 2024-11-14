package config

import (
	"context"
	"fmt"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

type SecretManager struct {
	client    *secretmanager.Client
	projectID string
}

func NewSecretManager(projectID string) (*SecretManager, error) {
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create secretmanager client: %w", err)
	}

	return &SecretManager{
		client:    client,
		projectID: projectID,
	}, nil
}

func (sm *SecretManager) GetOpenAIKey() (string, error) {
	ctx := context.Background()
	name := fmt.Sprintf("projects/%s/secrets/OPENAI_API_KEY/versions/latest", sm.projectID)

	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	}

	result, err := sm.client.AccessSecretVersion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to access secret version: %w", err)
	}

	return string(result.Payload.Data), nil
}

func (sm *SecretManager) Close() error {
	return sm.client.Close()
}
