package providers

import (
	"go.uber.org/zap"
	"juraji.nl/chat-quest/core/database"
	"juraji.nl/chat-quest/core/log"
	"juraji.nl/chat-quest/core/util"
)

type LlmModel struct {
	ID                  int     `json:"id"`
	ConnectionProfileId int     `json:"profileId"`
	ModelId             string  `json:"modelId"`
	Temperature         float32 `json:"temperature"`
	MaxTokens           int     `json:"maxTokens"`
	TopP                float32 `json:"topP"`
	Stream              bool    `json:"stream"`
	StopSequences       *string `json:"stopSequences"`
	Disabled            bool    `json:"disabled"`
}

type LlmModelView struct {
	ID                    int    `json:"id"`
	ModelId               string `json:"modelId"`
	ConnectionProfileId   int    `json:"profileId"`
	ConnectionProfileName string `json:"profileName"`
}

func llmModelScanner(scanner database.RowScanner, dest *LlmModel) error {
	return scanner.Scan(
		&dest.ID,
		&dest.ConnectionProfileId,
		&dest.ModelId,
		&dest.Temperature,
		&dest.MaxTokens,
		&dest.TopP,
		&dest.Stream,
		&dest.StopSequences,
		&dest.Disabled,
	)
}

func llModelViewScanner(scanner database.RowScanner, dest *LlmModelView) error {
	return scanner.Scan(
		&dest.ID,
		&dest.ModelId,
		&dest.ConnectionProfileId,
		&dest.ConnectionProfileName,
	)
}

func LlmModelsByConnectionProfileId(profileId int) ([]LlmModel, error) {
	query := "SELECT * FROM llm_models WHERE connection_profile_id = ?"
	args := []any{profileId}
	return database.QueryForList(query, args, llmModelScanner)
}

func createLlmModel(ctx *database.TxContext, profileId int, llmModel *LlmModel) error {
	llmModel.ConnectionProfileId = profileId

	query := `INSERT INTO llm_models
            (connection_profile_id, model_id, temperature, max_tokens, top_p, stream, stop_sequences, disabled)
            VALUES (?, ?, ?, ?, ?, ?, ?, ?) RETURNING id`
	args := []any{
		llmModel.ConnectionProfileId,
		llmModel.ModelId,
		llmModel.Temperature,
		llmModel.MaxTokens,
		llmModel.TopP,
		llmModel.Stream,
		llmModel.StopSequences,
		llmModel.Disabled,
	}

	return ctx.InsertRecord(query, args, &llmModel.ID)
}

func UpdateLlmModel(id int, llmModel *LlmModel) error {
	query := `UPDATE llm_models
              SET temperature = ?,
                  max_tokens = ?,
                  top_p = ?,
                  stream = ?,
                  stop_sequences = ?,
                  disabled = ?
              WHERE id = ?`
	args := []any{
		llmModel.Temperature,
		llmModel.MaxTokens,
		llmModel.TopP,
		llmModel.Stream,
		llmModel.StopSequences,
		llmModel.Disabled,
		id,
	}

	err := database.UpdateRecord(query, args)

	if err == nil {
		LlmModelUpdatedSignal.EmitBG(llmModel)
	}

	return err
}

func DeleteLlmModelById(id int) error {
	return database.Transactional(func(ctx *database.TxContext) error {
		return deleteLlmModelById(ctx, id)
	})
}

func deleteLlmModelById(ctx *database.TxContext, id int) error {
	query := "DELETE FROM llm_models WHERE id = ?"
	args := []any{id}

	err := ctx.DeleteRecord(query, args)

	if err == nil {
		LlmModelDeletedSignal.EmitBG(id)
	}

	return err
}

func MergeLlmModels(profileId int, newModels []*LlmModel) error {
	// New model id set
	newModelIdSet := util.NewSetFrom(newModels, func(t *LlmModel) string {
		return t.ModelId
	})

	// Existing model id set
	existingModels, err := LlmModelsByConnectionProfileId(profileId)
	if err != nil {
		return err
	}

	existingModelIdSet := util.NewSetFrom(existingModels, func(t LlmModel) string {
		return t.ModelId
	})

	var createdModels []*LlmModel
	var deletedModelIds []int
	logger := log.Get().With(zap.Int("profileId", profileId))

	err = database.Transactional(func(ctx *database.TxContext) error {
		// Add new models
		for _, newModel := range newModels {
			if existingModelIdSet.NotContains(newModel.ModelId) {
				if err := createLlmModel(ctx, profileId, newModel); err != nil {
					logger.Error("Error saving new llm model",
						zap.String("modelId", newModel.ModelId), zap.Error(err))
					return err
				}
				createdModels = append(createdModels, newModel)
			}
		}

		// Remove models not in new set
		for _, existingModel := range existingModels {
			if newModelIdSet.NotContains(existingModel.ModelId) {
				if err := deleteLlmModelById(ctx, existingModel.ID); err != nil {
					logger.Error("Error deleting existing llm model",
						zap.Int("id", existingModel.ID), zap.Error(err))
					return err
				}
				deletedModelIds = append(deletedModelIds, existingModel.ID)
			}
		}

		return nil
	})

	if err == nil {
		LlmModelCreatedSignal.EmitAllBG(createdModels)
		LlmModelDeletedSignal.EmitAllBG(deletedModelIds)
	}

	return err
}

func GetAllLlmModelViews() ([]LlmModelView, error) {
	query := `SELECT lm.id       AS model_id,
                   lm.model_id AS model_model_id,
                   p.id       AS profile_id,
                   p.name     AS profile_name
            FROM llm_models lm
                     JOIN connection_profiles p on p.id = lm.connection_profile_id
                     WHERE lm.disabled = ?`
	return database.QueryForList(query, []any{false}, llModelViewScanner)
}

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
