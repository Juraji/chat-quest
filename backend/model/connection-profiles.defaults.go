package model

func GetConnectionProfileDefaults() []*ConnectionProfile {
	return []*ConnectionProfile{
		{
			ID:           0,
			Name:         "Open AI",
			ProviderType: "OPEN_AI",
			BaseUrl:      "https://api.openai.com/v1",
			ApiKey:       "",
		},
		{
			ID:           0,
			Name:         "X AI",
			ProviderType: "OPEN_AI",
			BaseUrl:      "https://api.x.ai/v1",
			ApiKey:       "",
		},
		{
			ID:           0,
			Name:         "LM Studio",
			ProviderType: "OPEN_AI",
			BaseUrl:      "https://localhost:1234/v1",
			ApiKey:       "lm-studio",
		},
	}
}

func DefaultLlmModel(ConnectionProfileId int64, ModelId string, opts ...func(*LlmModel)) *LlmModel {
	model := LlmModel{
		ConnectionProfileId: ConnectionProfileId,
		ModelId:             ModelId,
		Temperature:         1.0,
		MaxTokens:           256,
		TopP:                0.95,
		Stream:              false,
		Stop:                []string{},
	}

	for _, opt := range opts {
		opt(&model)
	}

	return &model
}
