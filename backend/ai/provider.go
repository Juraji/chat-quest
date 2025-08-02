package ai

import (
	"fmt"
	"juraji.nl/chat-quest/model"
)

type Provider interface {
	getAvailableModels() ([]*model.LlmModel, error)
	generateChatCompletions(llmModel model.LlmModel, messages []Message) (string, error)
}

type Message struct {
	Role    string
	Content string
}

func newProvider(profile model.ConnectionProfile) (Provider, error) {
	switch profile.ProviderType {
	case "OPEN_AI":
		return &openAIProvider{ConnectionProfile: profile}, nil
	default:
		return nil, fmt.Errorf("unknown provider type: %s", profile.ProviderType)
	}
}
