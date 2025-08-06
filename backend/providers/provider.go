package providers

import (
	"fmt"
)

type Provider interface {
	getAvailableModels() ([]*LlmModel, error)
	generateEmbeddings(input string, modelId string) (Embeddings, error)
}

func newProvider(profile *ConnectionProfile) (Provider, error) {
	switch profile.ProviderType {
	case "OPEN_AI":
		return &openAIProvider{connectionProfile: profile}, nil
	default:
		return nil, fmt.Errorf("unknown provider type: %s", profile.ProviderType)
	}
}

func (p *ConnectionProfile) GetAvailableModels() ([]*LlmModel, error) {
	provider, err := newProvider(p)
	if err != nil {
		return nil, fmt.Errorf("failed to create provider for profile %s (id %d): %w", p.Name, p.ID, err)
	}

	models, err := provider.getAvailableModels()
	if err != nil {
		return nil, fmt.Errorf("failed to get models for profile %s (id %d): %w", p.Name, p.ID, err)
	}

	return models, nil
}

func (p *ConnectionProfile) GenerateEmbeddings(input string, llmModel LlmModel) (Embeddings, error) {
	provider, err := newProvider(p)
	if err != nil {
		return nil, fmt.Errorf("failed to create provider for profile %s (id %d): %w", p.Name, p.ID, err)
	}

	embedding, err := provider.generateEmbeddings(input, llmModel.ModelId)
	if err != nil {
		return nil, fmt.Errorf("failed to generate embeddings for %s (id %d): %w", p.Name, p.ID, err)
	}

	return embedding, nil
}
