package providers

import (
	"juraji.nl/chat-quest/core/database"
)

type LlmModelInstance struct {
	ProviderId   int
	ProviderType ProviderType
	BaseUrl      string
	ApiKey       string
	ModelId      string
}

func llmModelInstanceScanner(scanner database.RowScanner, dest *LlmModelInstance) error {
	return scanner.Scan(
		&dest.ProviderId,
		&dest.ProviderType,
		&dest.BaseUrl,
		&dest.ApiKey,
		&dest.ModelId,
	)
}

func GetLlmModelInstanceById(llmModelId int) (*LlmModelInstance, error) {
	query := `SELECT
                cp.id as provider_id,
                cp.provider_type AS provider_type,
                cp.base_url AS base_url,
                cp.api_key AS api_key,
                lm.model_id AS model_id
            FROM llm_models lm
                JOIN connection_profiles cp on cp.id = lm.connection_profile_id
                WHERE lm.id = ?`
	args := []any{llmModelId}

	return database.QueryForRecord(query, args, llmModelInstanceScanner)
}
