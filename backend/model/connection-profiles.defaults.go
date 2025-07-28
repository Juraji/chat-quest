package model

import "fmt"

var openAiDefaults = ConnectionProfile{
	ID:           0,
	ProviderType: "OPEN_AI",
	BaseUrl:      "https://api.openai.com/v1",
	ApiKey:       "",
}

var xAiDefaults = ConnectionProfile{
	ID:           0,
	ProviderType: "X_AI",
	BaseUrl:      "https://api.x.ai/v1",
	ApiKey:       "",
}

var lmStudioDefaults = ConnectionProfile{
	ID:           0,
	ProviderType: "LM_STUDIO",
	BaseUrl:      "https://localhost:1234/v1",
	ApiKey:       "lm-studio",
}

func GetConnectionProfileDefaults(providerType string) (ConnectionProfile, error) {
	switch providerType {
	case "OPEN_AI":
		return openAiDefaults, nil
	case "X_AI":
		return xAiDefaults, nil
	case "LM_STUDIO":
		return lmStudioDefaults, nil
	default:
		return ConnectionProfile{}, fmt.Errorf("unknown providerType: %s", providerType)
	}
}
