package providers

import (
	"fmt"
	"juraji.nl/chat-quest/core/util"
	"strings"
)

type Provider interface {
	getAvailableModelIds() ([]string, error)
	generateEmbeddings(input string, modelId string) (util.Embeddings, error)
	generateChatResponse(request *ChatGenerateRequest) <-chan ChatGenerateResponse
}

type ChatGenerateRequest struct {
	Messages      []ChatRequestMessage
	ModelId       string
	MaxTokens     int
	Temperature   float32
	TopP          float32
	Stream        bool
	StopSequences []string
}

type ChatMessageRole string

const (
	RoleSystem    ChatMessageRole = "SYSTEM"
	RoleUser      ChatMessageRole = "USER"
	RoleAssistant ChatMessageRole = "ASSISTANT"
)

type ChatRequestMessage struct {
	Role    ChatMessageRole
	Content string
}

type ChatGenerateResponse struct {
	Content string
	Error   error
}

func newProvider(providerType ProviderType, baseUrl string, apiKey string) Provider {
	switch providerType {
	case ProviderOpenAi:
		return newOpenAiProvider(baseUrl, apiKey)
	default:
		panic(fmt.Sprintf("unknown provider type: %s", providerType))
	}
}

func (p *ConnectionProfile) GetAvailableModels() ([]*LlmModel, error) {
	provider := newProvider(p.ProviderType, p.BaseUrl, p.ApiKey)
	models, err := provider.getAvailableModelIds()
	if err != nil {
		return nil, fmt.Errorf("failed to get models for profile %s (id %d): %w", p.Name, p.ID, err)
	}

	llmModels := make([]*LlmModel, len(models))
	for i, model := range models {
		llmModels[i] = DefaultLlmModel(p.ID, model)
	}

	return llmModels, nil
}

func (lm *LlmModelInstance) GenerateEmbeddings(input string) (util.Embeddings, error) {
	provider := newProvider(lm.ProviderType, lm.BaseUrl, lm.ApiKey)
	embedding, err := provider.generateEmbeddings(input, lm.ModelId)
	if err != nil {
		return nil, fmt.Errorf("failed to generate embeddings for %s (%s): %w", lm.ModelId, lm.ProviderType, err)
	}

	return embedding, nil
}

func (lm *LlmModelInstance) GenerateChatResponse(
	messages []ChatRequestMessage,
	overrideTemperature *float32,
) <-chan ChatGenerateResponse {
	provider := newProvider(lm.ProviderType, lm.BaseUrl, lm.ApiKey)

	var temperature float32
	if overrideTemperature != nil {
		temperature = *overrideTemperature
	} else {
		temperature = lm.Temperature
	}

	var stopSequences []string
	if lm.StopSequences == nil || *lm.StopSequences == "" {
		stopSequences = nil
	} else {
		stopSequences = strings.Split(*lm.StopSequences, ",")
		for i := range stopSequences {
			stopSequences[i] = strings.TrimSpace(stopSequences[i])
		}
	}

	request := &ChatGenerateRequest{
		Messages:      messages,
		ModelId:       lm.ModelId,
		MaxTokens:     lm.MaxTokens,
		Temperature:   temperature,
		TopP:          lm.TopP,
		Stream:        lm.Stream,
		StopSequences: stopSequences,
	}

	return provider.generateChatResponse(request)
}
