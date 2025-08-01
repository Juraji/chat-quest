package ai

import (
	"fmt"
	"juraji.nl/chat-quest/model"
)

func NewProvider(profile model.ConnectionProfile) (Provider, error) {
	switch profile.ProviderType {
	case "OPEN_AI":
		return &openAIProvider{ConnectionProfile: profile}, nil
	default:
		return nil, fmt.Errorf("unknown provider type: %s", profile.ProviderType)
	}
}
