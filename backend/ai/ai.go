package ai

import (
	"fmt"
	"juraji.nl/chat-quest/model"
)

func GetAvailableModels(profile model.ConnectionProfile) ([]*model.LlmModel, error) {
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

func GenerateChatCompletions(
	profile model.ConnectionProfile,
	llmModel model.LlmModel,
	messages []Message,
) (string, error) {
	return "", fmt.Errorf("not yet implemented")
}
