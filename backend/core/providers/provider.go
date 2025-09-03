package providers

import (
	"fmt"
	"strings"
	"sync"
	"unicode"
)

var (
	providerExecutionLocksLock sync.Mutex
	// providerExecutionLocks is a map that stores mutex locks for each provider instance to prevent concurrent execution.
	// The key is the provider ID and the value is a pointer to a sync.Mutex.
	providerExecutionLocks = make(map[int]*sync.Mutex)
)

// Provider defines an interface for interacting with different AI model providers.
// Implementations should handle their own specific API calls while maintaining
// consistent behavior in terms of locking, error handling, and response formats.
type Provider interface {
	// getAvailableModelIds returns a list of available model IDs that can be used with this provider.
	// The returned models are provider-specific identifiers for different AI models.
	// Returns an empty slice if no models are available or an error occurred.
	getAvailableModelIds() ([]string, error)

	// generateEmbeddings creates vector embeddings from the given input text using the specified model.
	// The function should be thread-safe and handle its own locking internally if needed.
	// Returns the generated embeddings as a slice of floats or an error if generation failed.
	generateEmbeddings(input string, modelId string) (Embeddings, error)

	// generateChatResponse creates a channel that will stream chat responses based on the provided request.
	// The function should be thread-safe and handle its own locking internally if needed.
	// Returns a receive-only channel (<-chan) that will yield ChatGenerateResponse objects as they become available.
	generateChatResponse(request *ChatGenerateRequest) <-chan ChatGenerateResponse
}

// newProvider creates a new instance of the specified provider type.
// The function panics if an unknown provider type is requested.
func newProvider(providerType ProviderType, baseUrl string, apiKey string) Provider {
	switch providerType {
	case ProviderOpenAi:
		return newOpenAiProvider(baseUrl, apiKey)
	default:
		panic(fmt.Sprintf("unknown provider type: %s", providerType))
	}
}

// getProviderLock retrieves or creates a mutex lock for the specified provider ID.
// This ensures thread-safe access to provider instances by preventing concurrent execution.
func getProviderLock(providerId int) *sync.Mutex {
	providerExecutionLocksLock.Lock()
	defer providerExecutionLocksLock.Unlock()

	mutex, exists := providerExecutionLocks[providerId]
	if !exists {
		mutex = &sync.Mutex{}
		providerExecutionLocks[providerId] = mutex
	}
	return mutex
}

func cleanTextForEmbedding(text string) string {
	var builder strings.Builder
	const apos = '\''

	for _, char := range text {
		// Check if the character is alphanumeric or space/apostrophe (adjust as needed)
		if unicode.IsLetter(char) || unicode.IsNumber(char) || unicode.IsSpace(char) || char == apos {
			builder.WriteRune(unicode.ToLower(char))
		}
	}

	return strings.TrimSpace(builder.String())
}

// GetAvailableModels retrieves the list of available models for a given connection profile.
func GetAvailableModels(profile *ConnectionProfile) ([]*LlmModel, error) {
	provider := newProvider(profile.ProviderType, profile.BaseUrl, profile.ApiKey)
	models, err := provider.getAvailableModelIds()
	if err != nil {
		return nil, fmt.Errorf("failed to get models for profile %s (id %d): %w", profile.Name, profile.ID, err)
	}

	llmModels := make([]*LlmModel, len(models))
	for i, model := range models {
		llmModels[i] = DefaultLlmModel(profile.ID, model)
	}

	return llmModels, nil
}

// GenerateEmbeddings creates vector embeddings from the given input text using a specified LLM model.
// This function uses a provider-specific lock to ensure thread-safe access during embedding generation.
func GenerateEmbeddings(llm *LlmModelInstance, input string, cleanInput bool) (Embeddings, error) {
	providerLock := getProviderLock(llm.ProviderId)
	providerLock.Lock()
	defer providerLock.Unlock()

	if cleanInput {
		input = cleanTextForEmbedding(input)
	}

	provider := newProvider(llm.ProviderType, llm.BaseUrl, llm.ApiKey)
	embedding, err := provider.generateEmbeddings(input, llm.ModelId)
	if err != nil {
		return nil, fmt.Errorf("failed to generate embeddings for %s (%s): %w", llm.ModelId, llm.ProviderType, err)
	}

	return embedding, nil
}

// GenerateChatResponse creates a channel that will stream chat responses based on the provided messages and model configuration.
// This function uses a provider-specific lock to ensure thread-safe access during response generation.
func GenerateChatResponse(
	llm *LlmModelInstance,
	messages []ChatRequestMessage,
	overrideTemperature *float32,
) <-chan ChatGenerateResponse {
	providerLock := getProviderLock(llm.ProviderId)
	providerLock.Lock()
	defer providerLock.Unlock()

	provider := newProvider(llm.ProviderType, llm.BaseUrl, llm.ApiKey)

	var temperature float32
	if overrideTemperature != nil {
		temperature = *overrideTemperature
	} else {
		temperature = llm.Temperature
	}

	var stopSequences []string
	if llm.StopSequences == nil || *llm.StopSequences == "" {
		stopSequences = nil
	} else {
		stopSequences = strings.Split(*llm.StopSequences, ",")
		for i := range stopSequences {
			stopSequences[i] = strings.TrimSpace(stopSequences[i])
		}
	}

	request := &ChatGenerateRequest{
		Messages:      messages,
		ModelId:       llm.ModelId,
		MaxTokens:     llm.MaxTokens,
		Temperature:   temperature,
		TopP:          llm.TopP,
		Stream:        llm.Stream,
		StopSequences: stopSequences,
	}

	return provider.generateChatResponse(request)
}
