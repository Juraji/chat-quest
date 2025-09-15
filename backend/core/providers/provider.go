package providers

import (
	"context"
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
	getAvailableModelIds(ctx context.Context) ([]*LlmModel, error)

	// generateEmbeddings creates vector embeddings from the given input text using the specified model.
	// The function should be thread-safe and handle its own locking internally if needed.
	// Returns the generated embeddings as a slice of floats or an error if generation failed.
	generateEmbeddings(ctx context.Context, input string, modelId string) (Embedding, error)

	// generateChatResponse creates a channel that will stream chat responses based on the provided request.
	// The function should be thread-safe and handle its own locking internally if needed.
	// Returns a receive-only channel (<-chan) that will yield ChatGenerateResponse objects as they become available.
	generateChatResponse(ctx context.Context, messages []ChatRequestMessage, modelId string, params LlmParameters) <-chan ChatGenerateResponse
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

// GetAvailableModels retrieves the list of available models for a given connection profile.
func GetAvailableModels(profile *ConnectionProfile) ([]*LlmModel, error) {
	ctx := context.Background()
	provider := newProvider(profile.ProviderType, profile.BaseUrl, profile.ApiKey)
	models, err := provider.getAvailableModelIds(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get models for profile %s (id %d): %w", profile.Name, profile.ID, err)
	}

	return models, nil
}

// GenerateEmbeddings creates vector embeddings from the given input text using a specified LLM model.
// This function uses a provider-specific lock to ensure singular and thread-safe access during embedding generation.
func GenerateEmbeddings(llm *LlmModelInstance, input string, cleanInput bool) (Embedding, error) {
	providerLock := getProviderLock(llm.ProviderId)
	providerLock.Lock()
	defer providerLock.Unlock()
	ctx := context.Background()

	if cleanInput {
		var builder strings.Builder
		const apos = '\''

		for _, char := range input {
			if unicode.IsLetter(char) || unicode.IsNumber(char) || unicode.IsSpace(char) || char == apos {
				builder.WriteRune(unicode.ToLower(char))
			}
		}

		input = strings.TrimSpace(builder.String())
	}

	provider := newProvider(llm.ProviderType, llm.BaseUrl, llm.ApiKey)
	embedding, err := provider.generateEmbeddings(ctx, input, llm.ModelId)
	if err != nil {
		return nil, fmt.Errorf("failed to generate embeddings for %s (%s): %w", llm.ModelId, llm.ProviderType, err)
	}

	return embedding.Normalize(), nil
}

// GenerateChatResponse creates a channel that will stream chat responses based on the provided messages and model configuration.
// This function uses a provider-specific lock to ensure singular and thread-safe access during response generation.
func GenerateChatResponse(
	ctx context.Context,
	llm *LlmModelInstance,
	messages []ChatRequestMessage,
	params LlmParameters,
) <-chan ChatGenerateResponse {
	providerLock := getProviderLock(llm.ProviderId)
	providerLock.Lock()
	//defer providerLock.Unlock()

	provider := newProvider(llm.ProviderType, llm.BaseUrl, llm.ApiKey)
	responseChan := make(chan ChatGenerateResponse)

	go func() {
		defer close(responseChan)
		defer providerLock.Unlock()

		res := provider.generateChatResponse(ctx, messages, llm.ModelId, params)

		for msg := range res {
			responseChan <- msg
		}
	}()

	return responseChan
}
