package providers

import (
	"fmt"
	"juraji.nl/chat-quest/util"
)

type Provider interface {
	getAvailableModels() ([]*LlmModel, error)
	generateEmbeddings(input string, modelId string) (util.Embeddings, error)
	generateChatResponse(request *ChatGenerateRequest) chan ChatGenerateResponse
}

type ChatGenerateRequest struct {
	Messages      []*ChatGenerateRequestMessage
	ModelId       string
	MaxTokens     int
	Temperature   float32
	TopP          float32
	Stream        bool
	StopSequences []string
}

type ChatRequestMessageRole string

const (
	RoleSystem    ChatRequestMessageRole = "SYSTEM"
	RoleUser      ChatRequestMessageRole = "USER"
	RoleAssistant ChatRequestMessageRole = "ASSISTANT"
)

type ChatGenerateRequestMessage struct {
	Role    ChatRequestMessageRole
	Content string
}

type ChatGenerateResponse struct {
	Content string
	Error   error
}

func newProvider(profile *ConnectionProfile) Provider {
	switch profile.ProviderType {
	case ProviderOpenAi:
		return newOpenAiProvider(profile)
	default:
		panic(fmt.Sprintf("unknown provider type: %s", profile.ProviderType))
	}
}

func (p *ConnectionProfile) GetAvailableModels() ([]*LlmModel, error) {
	provider := newProvider(p)
	models, err := provider.getAvailableModels()
	if err != nil {
		return nil, fmt.Errorf("failed to get models for profile %s (id %d): %w", p.Name, p.ID, err)
	}

	return models, nil
}

func (p *ConnectionProfile) GenerateEmbeddings(input string, llmModel LlmModel) (util.Embeddings, error) {
	provider := newProvider(p)
	embedding, err := provider.generateEmbeddings(input, llmModel.ModelId)
	if err != nil {
		return nil, fmt.Errorf("failed to generate embeddings for %s (id %d): %w", p.Name, p.ID, err)
	}

	return embedding, nil
}

func (p *ConnectionProfile) GenerateChatResponse(
	messages []*ChatGenerateRequestMessage,
	llmModel LlmModel,
	overrideTemperature *float32,
) chan ChatGenerateResponse {
	provider := newProvider(p)

	var temperature float32
	if overrideTemperature != nil {
		temperature = *overrideTemperature
	} else {
		temperature = llmModel.Temperature
	}

	request := &ChatGenerateRequest{
		Messages:      messages,
		ModelId:       llmModel.ModelId,
		MaxTokens:     llmModel.MaxTokens,
		Temperature:   temperature,
		TopP:          llmModel.TopP,
		Stream:        llmModel.Stream,
		StopSequences: llmModel.GetStopSequences(),
	}

	return provider.generateChatResponse(request)
}
