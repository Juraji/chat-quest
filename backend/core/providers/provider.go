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

func newProvider(profile *ConnectionProfile) Provider {
	switch profile.ProviderType {
	case ProviderOpenAi:
		return newOpenAiProvider(profile.BaseUrl, profile.ApiKey)
	default:
		panic(fmt.Sprintf("unknown provider type: %s", profile.ProviderType))
	}
}

func (p *ConnectionProfile) GetAvailableModels() ([]*LlmModel, error) {
	provider := newProvider(p)
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

func (p *ConnectionProfile) GenerateEmbeddings(input string, llmModel LlmModel) (util.Embeddings, error) {
	provider := newProvider(p)
	embedding, err := provider.generateEmbeddings(input, llmModel.ModelId)
	if err != nil {
		return nil, fmt.Errorf("failed to generate embeddings for %s (id %d): %w", p.Name, p.ID, err)
	}

	return embedding, nil
}

func (lm *LlmModel) GetStopSequences() []string {
	if lm.StopSequences == nil || *lm.StopSequences == "" {
		return nil
	}

	sequences := strings.Split(*lm.StopSequences, ",")
	for i := range sequences {
		sequences[i] = strings.TrimSpace(sequences[i])
	}
	return sequences
}

func (p *ConnectionProfile) GenerateChatResponse(
	messages []ChatRequestMessage,
	llmModel LlmModel,
	overrideTemperature *float32,
) <-chan ChatGenerateResponse {
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
