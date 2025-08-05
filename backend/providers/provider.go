package providers

import (
	"fmt"
)

type Provider interface {
	getAvailableModels() ([]*LlmModel, error)
}

func newProvider(profile *ConnectionProfile) (Provider, error) {
	switch profile.ProviderType {
	case "OPEN_AI":
		return &openAIProvider{ConnectionProfile: profile}, nil
	default:
		return nil, fmt.Errorf("unknown provider type: %s", profile.ProviderType)
	}
}

func (profile *ConnectionProfile) getAvailableModels() ([]*LlmModel, error) {
	provider, err := newProvider(profile)
	if err != nil {
		return nil, fmt.Errorf("failed to create provider for profile %s (id %d): %w", profile.Name, profile.ID, err)
	}

	models, err := provider.getAvailableModels()
	if err != nil {
		return nil, fmt.Errorf("failed to get models for profile %s (id %d): %w", profile.Name, profile.ID, err)
	}

	return models, nil
}
