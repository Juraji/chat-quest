package providers

import (
	"go.uber.org/zap"
	"juraji.nl/chat-quest/core/database"
	"juraji.nl/chat-quest/core/log"
)

type LlmModelInstance struct {
	ProviderId    int
	ProviderType  ProviderType
	BaseUrl       string
	ApiKey        string
	ModelId       string
	Temperature   float32
	MaxTokens     int
	TopP          float32
	Stream        bool
	StopSequences *string
}

func llmModelInstanceScanner(scanner database.RowScanner, dest *LlmModelInstance) error {
	return scanner.Scan(
		&dest.ProviderId,
		&dest.ProviderType,
		&dest.BaseUrl,
		&dest.ApiKey,
		&dest.ModelId,
		&dest.Temperature,
		&dest.MaxTokens,
		&dest.TopP,
		&dest.Stream,
		&dest.StopSequences,
	)
}

func GetLlmModelInstanceById(llmModelId int) (*LlmModelInstance, bool) {
	query := `SELECT
                cp.id as provider_id,
                cp.provider_type AS provider_type,
                cp.base_url AS base_url,
                cp.api_key AS api_key,
                lm.model_id AS model_id,
                lm.temperature AS temperature,
                lm.max_tokens AS max_tokens,
                lm.top_p AS top_p,
                lm.stream AS stream,
                lm.stop_sequences AS stop_sequences
            FROM llm_models lm
                JOIN connection_profiles cp on cp.id = lm.connection_profile_id
                WHERE lm.id = ?`
	args := []any{llmModelId}

	inst, err := database.QueryForRecord(query, args, llmModelInstanceScanner)
	if err != nil {
		log.Get().Error("Error querying for llm model", zap.Int("modelId", llmModelId), zap.Error(err))
		return nil, false
	}
	return inst, true
}
