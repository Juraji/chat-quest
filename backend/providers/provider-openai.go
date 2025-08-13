package providers

import (
	"context"
	"errors"
	"fmt"
	openai "github.com/sashabaranov/go-openai"
	"io"
	"juraji.nl/chat-quest/util"
	"math"
)

type openAIProvider struct {
	connectionProfile *ConnectionProfile
	client            *openai.Client
	ctx               context.Context
}

func newOpenAiProvider(profile *ConnectionProfile) *openAIProvider {
	config := openai.DefaultConfig(profile.ApiKey)
	config.BaseURL = profile.BaseUrl

	return &openAIProvider{
		connectionProfile: profile,
		client:            openai.NewClientWithConfig(config),
		ctx:               context.Background(),
	}
}

func (o *openAIProvider) getAvailableModels() ([]*LlmModel, error) {
	models, err := o.client.ListModels(o.ctx)
	if err != nil {
		return nil, fmt.Errorf("openAIProvider failed to list models: %w", err)
	}

	llmModels := make([]*LlmModel, len(models.Models))
	for i, model := range models.Models {
		llmModels[i] = defaultLlmModel(o.connectionProfile.ID, model.ID)
	}

	return llmModels, nil
}

func (o *openAIProvider) generateEmbeddings(input, modelID string) (util.Embeddings, error) {
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

func (o *openAIProvider) generateChatResponse(request *ChatGenerateRequest) chan ChatGenerateResponse {
	messages := make([]openai.ChatCompletionMessage, len(request.Messages))
	for i, msg := range request.Messages {
		var role string
		switch msg.Role {
		case RoleSystem:
			role = openai.ChatMessageRoleSystem
		case RoleUser:
			role = openai.ChatMessageRoleUser
		case RoleAssistant:
			role = openai.ChatMessageRoleAssistant
		default:
			responseChannel := make(chan ChatGenerateResponse)
			go func() {
				defer close(responseChannel)
				responseChannel <- ChatGenerateResponse{
					Error: fmt.Errorf("openAIProvider unknown role: %s", msg.Role),
				}
			}()
			return responseChannel
		}

		messages[i] = openai.ChatCompletionMessage{
			Role:    role,
			Content: msg.Content,
		}
	}

	completionRequest := openai.ChatCompletionRequest{
		Model:               request.ModelId,
		Messages:            messages,
		MaxTokens:           int(request.MaxTokens),
		MaxCompletionTokens: int(request.MaxTokens),
		Temperature:         float32(math.Max(math.SmallestNonzeroFloat32, request.Temperature)),
		TopP:                float32(request.TopP),
		Stream:              request.Stream,
		Stop:                request.StopSequences,
		PresencePenalty:     0,
		FrequencyPenalty:    0,
		StreamOptions:       nil,
		Store:               false,
		Prediction:          nil,
	}

	responseChannel := make(chan ChatGenerateResponse, len(messages))
	go func() {
		defer close(responseChannel)

		if request.Stream {
			stream, err := o.client.CreateChatCompletionStream(o.ctx, completionRequest)
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
		} else {
			completion, err := o.client.CreateChatCompletion(o.ctx, completionRequest)
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
		}
	}()

	return responseChannel
}
