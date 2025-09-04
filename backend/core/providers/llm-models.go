package providers

import (
	"go.uber.org/zap"
	"juraji.nl/chat-quest/core/database"
	"juraji.nl/chat-quest/core/log"
	"juraji.nl/chat-quest/core/util"
)

type LlmModelType string

const (
	UnknownModel   LlmModelType = "UNKNOWN"
	ChatModel      LlmModelType = "CHAT_MODEL"
	EmbeddingModel LlmModelType = "EMBEDDING_MODEL"
)

func (i LlmModelType) IsValid() bool {
	switch i {
	case UnknownModel:
		return true
	case ChatModel:
		return true
	case EmbeddingModel:
		return true
	default:
		return false
	}
}

type LlmModel struct {
	ID                  int          `json:"id"`
	ConnectionProfileId int          `json:"profileId"`
	ModelId             string       `json:"modelId"`
	ModelType           LlmModelType `json:"modelType"`
	Disabled            bool         `json:"disabled"`
}

type LlmModelView struct {
	ID                    int          `json:"id"`
	ModelId               string       `json:"modelId"`
	ModelType             LlmModelType `json:"modelType"`
	ConnectionProfileId   int          `json:"profileId"`
	ConnectionProfileName string       `json:"profileName"`
}

func llmModelScanner(scanner database.RowScanner, dest *LlmModel) error {
	return scanner.Scan(
		&dest.ID,
		&dest.ConnectionProfileId,
		&dest.ModelId,
		&dest.ModelType,
		&dest.Disabled,
	)
}

func llModelViewScanner(scanner database.RowScanner, dest *LlmModelView) error {
	return scanner.Scan(
		&dest.ID,
		&dest.ModelId,
		&dest.ModelType,
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
	if llmModel.ModelType == "" {
		llmModel.ModelType = UnknownModel
	}

	query := `INSERT INTO llm_models
            (connection_profile_id, model_id, model_type, disabled)
            VALUES (?, ?, ?, ?) RETURNING id`
	args := []any{
		llmModel.ConnectionProfileId,
		llmModel.ModelId,
		llmModel.ModelType,
		llmModel.Disabled,
	}

	return ctx.InsertRecord(query, args, &llmModel.ID)
}

func UpdateLlmModel(id int, llmModel *LlmModel) error {
	query := `UPDATE llm_models
              SET model_type= ?,
                  disabled = ?
              WHERE id = ?`
	args := []any{
		llmModel.ModelType,
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

	_, err := ctx.DeleteRecord(query, args)

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
                     lm.model_type AS model_type,
                     p.id       AS profile_id,
                     p.name     AS profile_name
              FROM llm_models lm
                       JOIN connection_profiles p on p.id = lm.connection_profile_id
                       WHERE lm.disabled = FALSE`
	return database.QueryForList(query, nil, llModelViewScanner)
}
