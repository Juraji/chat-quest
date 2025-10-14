package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
	"github.com/pkg/errors"
)

type openAIProvider struct {
	client openai.Client
	lock   *sync.Mutex
}

func newOpenAiProvider(baseUrl string, apiKey string) *openAIProvider {
	return &openAIProvider{
		client: openai.NewClient(
			option.WithBaseURL(baseUrl),
			option.WithAPIKey(apiKey),
		),
		lock: &sync.Mutex{},
	}
}

func (o *openAIProvider) getAvailableModelIds(ctx context.Context) ([]*LlmModel, error) {
	modelsIter := o.client.Models.ListAutoPaging(ctx)
	if err := modelsIter.Err(); err != nil {
		return nil, fmt.Errorf("openAIProvider failed to list models: %w", err)
	}

	var llmModels []*LlmModel

	for modelsIter.Next() {
		if err := modelsIter.Err(); err != nil {
			return nil, fmt.Errorf("openAIProvider failed to list models: %w", err)
		}

		model := modelsIter.Current()

		// OpenAI endpoints don't have a type, but we can generally infer from model id here.
		var t LlmModelType
		if strings.Contains(model.ID, "embedding-") {
			t = EmbeddingModel
		} else {
			t = ChatModel
		}

		llmModels = append(llmModels, &LlmModel{
			ModelId:   model.ID,
			ModelType: t,
		})
	}

	return llmModels, nil
}

func (o *openAIProvider) generateEmbeddings(ctx context.Context, input, modelID string) (Embedding, error) {
	o.lock.Lock()
	defer o.lock.Unlock()

	response, err := o.client.Embeddings.New(ctx, openai.EmbeddingNewParams{
		Input: openai.EmbeddingNewParamsInputUnion{
			OfString: openai.String(input),
		},
		Model:          modelID,
		EncodingFormat: openai.EmbeddingNewParamsEncodingFormatFloat,
	})
	if err != nil {
		return nil, fmt.Errorf("openAIProvider failed to create embeddings: %w", err)
	}
	if len(response.Data) == 0 {
		return nil, errors.New("openAIProvider no embeddings returned")
	}

	return response.Data[0].Embedding, nil
}

func (o *openAIProvider) generateChatResponse(ctx context.Context, messages []ChatRequestMessage, modelId string, params LlmParameters) <-chan ChatGenerateResponse {
	oMessages := make([]openai.ChatCompletionMessageParamUnion, len(messages))

	for i, msg := range messages {
		switch msg.Role {
		case RoleSystem:
			oMessages[i] = openai.SystemMessage(msg.Content)
		case RoleUser:
			oMessages[i] = openai.UserMessage(msg.Content)
		case RoleAssistant:
			oMessages[i] = openai.AssistantMessage(msg.Content)
		default:
			// Dev error, missing branch?
			panic(fmt.Errorf("developer error, invalid role '%s'", msg.Role))
		}
	}

	completionParams := openai.ChatCompletionNewParams{
		Messages:            oMessages,
		Model:               modelId,
		MaxTokens:           openai.Int(int64(params.MaxTokens)),
		MaxCompletionTokens: openai.Int(int64(params.MaxTokens)),
		Temperature:         openai.Float(float64(params.Temperature)),
		TopP:                openai.Float(float64(params.TopP)),
		Stop: openai.ChatCompletionNewParamsStopUnion{
			OfStringArray: params.StopSequencesAsSlice(),
		},
		PresencePenalty:  openai.Float(float64(params.PresencePenalty)),
		FrequencyPenalty: openai.Float(float64(params.FrequencyPenalty)),
	}

	if params.ResponseFormat != nil {
		var schema interface{}
		err := json.Unmarshal([]byte(*params.ResponseFormat), &schema)
		if err != nil {
			panic(errors.Wrap(err, "Error parsing JSON schema:"))
		}
		completionParams.ResponseFormat = openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{
				JSONSchema: openai.ResponseFormatJSONSchemaJSONSchemaParam{
					Name:        "ResponseFormat",
					Description: openai.String("Default response schema"),
					Strict:      openai.Bool(true),
					Schema:      schema,
				},
			},
		}
	}

	if params.Stream {
		// Include usage options in final chunk
		completionParams.StreamOptions = openai.ChatCompletionStreamOptionsParam{
			IncludeUsage: openai.Bool(true),
		}

		return o.generateChatResponseStream(ctx, completionParams)
	} else {
		return o.generateChatResponseSingle(ctx, completionParams)
	}
}

func (o *openAIProvider) generateChatResponseSingle(
	ctx context.Context,
	params openai.ChatCompletionNewParams,
) <-chan ChatGenerateResponse {
	responseChannel := make(chan ChatGenerateResponse, 1)
	go func() {
		o.lock.Lock()
		defer o.lock.Unlock()
		defer close(responseChannel)

		completion, err := o.client.Chat.Completions.New(ctx, params)
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
			Content:     completion.Choices[0].Message.Content,
			TotalTokens: int(completion.Usage.TotalTokens),
		}
	}()
	return responseChannel
}

func (o *openAIProvider) generateChatResponseStream(
	ctx context.Context,
	params openai.ChatCompletionNewParams,
) <-chan ChatGenerateResponse {
	responseChannel := make(chan ChatGenerateResponse)

	go func() {
		o.lock.Lock()
		defer o.lock.Unlock()
		defer close(responseChannel)

		stream := o.client.Chat.Completions.NewStreaming(ctx, params)

		for stream.Next() {
			if err := stream.Err(); err != nil {
				responseChannel <- ChatGenerateResponse{
					Error: fmt.Errorf("openAIProvider error during chat completion stream: %w", err),
				}
				return
			}

			if ctx.Err() != nil {
				// Context in error, probably canceled, stop streaming
				err := stream.Close()
				if err != nil {
					panic("openAIProvider failed to close chat completion stream")
				}
				return
			}

			chunk := stream.Current()

			if len(chunk.Choices) == 0 {
				if chunk.Usage.TotalTokens == 0 {
					responseChannel <- ChatGenerateResponse{
						Error: errors.New("openAIProvider returned empty chat completions chunk"),
					}

					err := stream.Close()
					if err != nil {
						panic("openAIProvider failed to close chat completion stream")
					}
					return
				} else {
					responseChannel <- ChatGenerateResponse{
						TotalTokens:      int(chunk.Usage.TotalTokens),
						CompletionTokens: int(chunk.Usage.CompletionTokens),
					}
				}

			} else {
				responseChannel <- ChatGenerateResponse{
					Content: chunk.Choices[0].Delta.Content,
				}
			}
		}
	}()

	return responseChannel
}
