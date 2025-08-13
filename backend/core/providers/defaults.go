package providers

func DefaultLlmModel(ConnectionProfileId int, ModelId string, opts ...func(*LlmModel)) *LlmModel {
	model := LlmModel{
		ConnectionProfileId: ConnectionProfileId,
		ModelId:             ModelId,
		Temperature:         1.0,
		MaxTokens:           300,
		TopP:                0.95,
		Stream:              false,
	}

	for _, opt := range opts {
		opt(&model)
	}

	return &model
}
