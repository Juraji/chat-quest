package providers

import (
	"context"
	"errors"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"io"
)

type openAIProvider struct {
	client *openai.Client
	ctx    context.Context
}

func newOpenAiProvider(baseUrl string, apiKey string) *openAIProvider {
	config := openai.DefaultConfig(apiKey)
	config.BaseURL = baseUrl

	return &openAIProvider{
		client: openai.NewClientWithConfig(config),
		ctx:    context.Background(),
	}
}

func (o *openAIProvider) getAvailableModelIds() ([]string, error) {
	models, err := o.client.ListModels(o.ctx)
	if err != nil {
		return nil, fmt.Errorf("openAIProvider failed to list models: %w", err)
	}

	modelIds := make([]string, len(models.Models))
	for i, model := range models.Models {
		modelIds[i] = model.ID
	}

	return modelIds, nil
}

func (o *openAIProvider) generateEmbeddings(input, modelID string) (Embeddings, error) {
	request := openai.EmbeddingRequest{
		Input: input,
		Model: openai.EmbeddingModel(modelID),
	}

	embeddings, err := o.client.CreateEmbeddings(o.ctx, request)
	if err != nil {
		return nil, fmt.Errorf("openAIProvider failed to create embeddings: %w", err)
	}
	if len(embeddings.Data) == 0 {
		return nil, errors.New("openAIProvider no embeddings returned")
	}

	return embeddings.Data[0].Embedding, nil
}

func (o *openAIProvider) generateChatResponse(
	messages []ChatRequestMessage,
	modelId string,
	params LlmParameters,
) <-chan ChatGenerateResponse {
	oMessages := make([]openai.ChatCompletionMessage, len(messages))
	for i, msg := range messages {
		var openAiRole string

		switch msg.Role {
		case RoleSystem:
			openAiRole = openai.ChatMessageRoleSystem
		case RoleUser:
			openAiRole = openai.ChatMessageRoleUser
		case RoleAssistant:
			openAiRole = openai.ChatMessageRoleAssistant
		default:
			// Dev error, missing branch?
			panic(fmt.Errorf("developer error, invalid role '%s'", msg.Role))
		}

		oMessages[i] = openai.ChatCompletionMessage{
			Role:    openAiRole,
			Content: msg.Content,
		}
	}

	//goland:noinspection GoDeprecation for MaxTokens, we supply it for compat reasons
	completionRequest := openai.ChatCompletionRequest{
		Model:               modelId,
		Messages:            oMessages,
		MaxTokens:           params.MaxTokens,
		MaxCompletionTokens: params.MaxTokens,
		Temperature:         params.Temperature,
		TopP:                params.TopP,
		Stream:              params.Stream,
		Stop:                params.StopSequencesAsSlice(),
		PresencePenalty:     params.PresencePenalty,
		FrequencyPenalty:    params.FrequencyPenalty,
		StreamOptions:       nil,
		Store:               false,
		Prediction:          nil,
	}

	if params.Stream {
		return generateChatResponseStream(o.ctx, o.client, completionRequest)
	} else {
		return generateChatResponseSingle(o.ctx, o.client, completionRequest)
	}
}

func generateChatResponseSingle(
	ctx context.Context,
	client *openai.Client,
	completionRequest openai.ChatCompletionRequest,
) <-chan ChatGenerateResponse {
	responseChannel := make(chan ChatGenerateResponse, 1)
	go func() {
		defer close(responseChannel)

		completion, err := client.CreateChatCompletion(ctx, completionRequest)
		if err != nil {
			responseChannel <- ChatGenerateResponse{
				Error: fmt.Errorf("openAIProvider failed to create completion: %w", err),
			}
			return
		}

		if len(completion.Choices) == 0 {
			responseChannel <- ChatGenerateResponse{
				Error: errors.New("openAIProvider no chat completions returned"),
			}
			return
		}

		responseChannel <- ChatGenerateResponse{
			Content: completion.Choices[0].Message.Content,
		}
	}()
	return responseChannel
}

func generateChatResponseStream(
	ctx context.Context,
	client *openai.Client,
	completionRequest openai.ChatCompletionRequest,
) <-chan ChatGenerateResponse {
	responseChannel := make(chan ChatGenerateResponse)
	go func() {
		defer close(responseChannel)

		stream, err := client.CreateChatCompletionStream(ctx, completionRequest)
		if err != nil {
			responseChannel <- ChatGenerateResponse{
				Error: fmt.Errorf("openAIProvider failed to create chat completion stream: %w", err),
			}
			return
		}
		defer stream.Close()

		for {
			response, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				return
			}

			if err != nil {
				responseChannel <- ChatGenerateResponse{
					Error: fmt.Errorf("openAIProvider chat completion stream error: %w", err),
				}
				return
			}

			responseChannel <- ChatGenerateResponse{
				Content: response.Choices[0].Delta.Content,
			}
		}
	}()
	return responseChannel
}
