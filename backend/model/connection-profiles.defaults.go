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

func GetConnectionProfileDefaults(providerType string) (ConnectionProfile, error) {
	switch providerType {
	case "OPEN_AI":
		return openAiDefaults, nil
	case "X_AI":
		return xAiDefaults, nil
	default:
		return ConnectionProfile{}, fmt.Errorf("unknown providerType: %s", providerType)
	}
}
