package providers

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"unicode"
)

var (
	providerInstanceMapLock sync.Mutex
	providerInstances       = make(map[int]Provider)
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

// getProviderLock retrieves or creates a new instance of the specified provider type.
// This ensures thread-safe access to provider instances by preventing concurrent execution.
func getProvider(providerId int, providerType ProviderType, baseUrl string, apiKey string) Provider {
	providerInstanceMapLock.Lock()
	defer providerInstanceMapLock.Unlock()

	p, exists := providerInstances[providerId]
	if !exists {
		switch providerType {
		case ProviderOpenAi:
			p = newOpenAiProvider(baseUrl, apiKey)
		default:
			panic(fmt.Sprintf("unknown provider type: %s", providerType))
		}

		providerInstances[providerId] = p
	}
	return p
}

// GetAvailableModels retrieves the list of available models for a given connection profile.
func GetAvailableModels(profile *ConnectionProfile) ([]*LlmModel, error) {
	ctx := context.Background()
	provider := getProvider(profile.ID, profile.ProviderType, profile.BaseUrl, profile.ApiKey)
	models, err := provider.getAvailableModelIds(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get models for profile %s (id %d): %w", profile.Name, profile.ID, err)
	}

	return models, nil
}

// GenerateEmbeddings creates vector embeddings from the given input text using a specified LLM model.
func GenerateEmbeddings(llm *LlmModelInstance, input string, cleanInput bool) (Embedding, error) {
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

	provider := getProvider(llm.ProviderId, llm.ProviderType, llm.BaseUrl, llm.ApiKey)
	embedding, err := provider.generateEmbeddings(ctx, input, llm.ModelId)
	if err != nil {
		return nil, fmt.Errorf("failed to generate embeddings for %s (%s): %w", llm.ModelId, llm.ProviderType, err)
	}

	return embedding.Normalize(), nil
}

// GenerateChatResponse creates a channel that will stream chat responses based on the provided messages and model configuration.
func GenerateChatResponse(
	ctx context.Context,
	llm *LlmModelInstance,
	messages []ChatRequestMessage,
	params LlmParameters,
) <-chan ChatGenerateResponse {
	provider := getProvider(llm.ProviderId, llm.ProviderType, llm.BaseUrl, llm.ApiKey)
	return provider.generateChatResponse(ctx, messages, llm.ModelId, params)
}
