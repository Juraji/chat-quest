package providers

import (
	"fmt"
	"juraji.nl/chat-quest/cq"
	"juraji.nl/chat-quest/database"
	"juraji.nl/chat-quest/util"
	"strings"
)

type ProviderType string

const (
	ProviderOpenAi ProviderType = "OPEN_AI"
)

func (p ProviderType) IsValid() bool {
	switch p {
	case ProviderOpenAi:
		return true
	default:
		return false
	}
}

type ConnectionProfile struct {
	ID           int64        `json:"id"`
	Name         string       `json:"name"`
	ProviderType ProviderType `json:"providerType"`
	BaseUrl      string       `json:"baseUrl"`
	ApiKey       string       `json:"apiKey"`
}

type LlmModel struct {
	ID                  int64   `json:"id"`
	ConnectionProfileId int64   `json:"profileId"`
	ModelId             string  `json:"modelId"`
	Temperature         float64 `json:"temperature"`
	MaxTokens           int64   `json:"maxTokens"`
	TopP                float64 `json:"topP"`
	Stream              bool    `json:"stream"`
	StopSequences       *string `json:"stopSequences"`
	Disabled            bool    `json:"disabled"`
}

type LlmModelView struct {
	ID                    int64  `json:"id"`
	ModelId               string `json:"modelId"`
	ConnectionProfileId   int64  `json:"profileId"`
	ConnectionProfileName string `json:"profileName"`
}

func (lm *LlmModel) GetStopSequences() []string {
	if lm.StopSequences == nil || *lm.StopSequences == "" {
		return nil
	}

	sequences := strings.Split(*lm.StopSequences, ",")
	for i := range sequences {
		sequences[i] = strings.TrimSpace(sequences[i])
	}
	return sequences
}

func connectionProfileScanner(scanner database.RowScanner, dest *ConnectionProfile) error {
	return scanner.Scan(
		&dest.ID,
		&dest.Name,
		&dest.ProviderType,
		&dest.BaseUrl,
		&dest.ApiKey,
	)
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

func AllConnectionProfiles(cq *cq.ChatQuestContext) ([]*ConnectionProfile, error) {
	query := "SELECT * FROM connection_profiles"
	return database.QueryForList(cq.DB(), query, nil, connectionProfileScanner)
}

func ConnectionProfileById(cq *cq.ChatQuestContext, id int64) (*ConnectionProfile, error) {
	query := "SELECT * FROM connection_profiles WHERE id = ?"
	args := []any{id}
	return database.QueryForRecord(cq.DB(), query, args, connectionProfileScanner)
}

func CreateConnectionProfile(cq *cq.ChatQuestContext, profile *ConnectionProfile, llmModels []*LlmModel) error {
	tx, err := cq.DB().Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer database.RollBackOnErr(cq, tx, err)

	query := "INSERT INTO connection_profiles (name, provider_type, base_url, api_key) VALUES (?, ?, ?, ?) RETURNING id"
	args := []any{profile.Name, profile.ProviderType, profile.BaseUrl, profile.ApiKey}

	if err = database.InsertRecord(tx, query, args, &profile.ID); err != nil {
		return err
	}

	for _, llmModel := range llmModels {
		err = createLlmModel(tx, profile.ID, llmModel)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	ConnectionProfileCreatedSignal.Emit(cq.Context(), profile)
	util.EmitAll(cq, LlmModelCreatedSignal, llmModels)
	return nil
}

func UpdateConnectionProfile(cq *cq.ChatQuestContext, id int64, profile *ConnectionProfile) error {
	query := `UPDATE connection_profiles
            SET name = ?,
                provider_type = ?,
                base_url = ?,
                api_key = ?
            WHERE id = ?`
	args := []any{profile.Name, profile.ProviderType, profile.BaseUrl, profile.ApiKey, id}

	err := database.UpdateRecord(cq.DB(), query, args)
	if err != nil {
		return err
	}

	ConnectionProfileUpdatedSignal.Emit(cq.Context(), profile)
	return nil
}

func DeleteConnectionProfileById(cq *cq.ChatQuestContext, id int64) error {
	query := "DELETE FROM connection_profiles WHERE id = ?"
	args := []any{id}

	err := database.DeleteRecord(cq.DB(), query, args)
	if err != nil {
		return err
	}

	ConnectionProfileDeletedSignal.Emit(cq.Context(), id)
	return nil
}

func LlmModelsByConnectionProfileId(cq *cq.ChatQuestContext, profileId int64) ([]*LlmModel, error) {
	query := "SELECT * FROM llm_models WHERE connection_profile_id = ?"
	args := []any{profileId}
	return database.QueryForList(cq.DB(), query, args, llmModelScanner)
}

func CreateLlmModel(cq *cq.ChatQuestContext, profileId int64, llmModel *LlmModel) error {
	err := createLlmModel(cq.DB(), profileId, llmModel)
	if err != nil {
		return err
	}

	LlmModelCreatedSignal.Emit(cq.Context(), llmModel)
	return nil
}

func createLlmModel(db database.QueryExecutor, profileId int64, llmModel *LlmModel) error {
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

	return database.InsertRecord(db, query, args, &llmModel.ID)
}

func UpdateLlmModel(cq *cq.ChatQuestContext, id int64, llmModel *LlmModel) error {
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

	err := database.UpdateRecord(cq.DB(), query, args)
	if err != nil {
		return err
	}

	LlmModelUpdatedSignal.Emit(cq.Context(), llmModel)
	return nil
}

func DeleteLlmModelById(cq *cq.ChatQuestContext, id int64) error {
	err := deleteLlmModelById(cq.DB(), id)
	if err != nil {
		return err
	}

	LlmModelDeletedSignal.Emit(cq.Context(), id)
	return nil
}

func deleteLlmModelById(db database.QueryExecutor, id int64) error {
	query := "DELETE FROM llm_models WHERE id = ?"
	args := []any{id}

	return database.DeleteRecord(db, query, args)
}

func MergeLlmModels(cq *cq.ChatQuestContext, profileId int64, newModels []*LlmModel) error {
	tx, err := cq.DB().Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer database.RollBackOnErr(cq, tx, err)

	// New model id set
	newModelIdSet := util.NewSetFrom(newModels, func(t *LlmModel) string {
		return t.ModelId
	})

	// Existing model id set
	existingModels, err := LlmModelsByConnectionProfileId(cq, profileId)
	if err != nil {
		return err
	}

	existingModelIdSet := util.NewSetFrom(existingModels, func(t *LlmModel) string {
		return t.ModelId
	})

	// Add new models
	var createdModels []*LlmModel
	for _, newModel := range newModels {
		if existingModelIdSet.NotContains(newModel.ModelId) {
			if err = createLlmModel(tx, profileId, newModel); err != nil {
				return err
			}
			createdModels = append(createdModels, newModel)
		}
	}

	// Remove models not in new set
	var deletedModelIds []int64
	for _, existingModel := range existingModels {
		if newModelIdSet.NotContains(existingModel.ModelId) {
			if err = deleteLlmModelById(tx, existingModel.ID); err != nil {
				return err
			}
			deletedModelIds = append(deletedModelIds, existingModel.ID)
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	util.EmitAll(cq, LlmModelCreatedSignal, createdModels)
	util.EmitAll(cq, LlmModelDeletedSignal, deletedModelIds)
	return nil
}

func GetAllLlmModelViews(cq *cq.ChatQuestContext) ([]*LlmModelView, error) {
	query := `SELECT lm.id       AS model_id,
                   lm.model_id AS model_model_id,
                   p.id       AS profile_id,
                   p.name     AS profile_name
            FROM llm_models lm
                     JOIN connection_profiles p on p.id = lm.connection_profile_id
                     WHERE lm.disabled = ?`
	return database.QueryForList(cq.DB(), query, []any{false}, llModelViewScanner)
}
