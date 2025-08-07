package providers

import (
	"database/sql"
	"fmt"
	"juraji.nl/chat-quest/database"
	"juraji.nl/chat-quest/util"
)

type ConnectionProfile struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	ProviderType string `json:"providerType"`
	BaseUrl      string `json:"baseUrl"`
	ApiKey       string `json:"apiKey"`
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

func AllConnectionProfiles(db *sql.DB) ([]*ConnectionProfile, error) {
	query := "SELECT * FROM connection_profiles"
	return database.QueryForList(db, query, nil, connectionProfileScanner)
}

func ConnectionProfileById(db *sql.DB, id int64) (*ConnectionProfile, error) {
	query := "SELECT * FROM connection_profiles WHERE id = $1"
	args := []any{id}
	return database.QueryForRecord(db, query, args, connectionProfileScanner)
}

func CreateConnectionProfile(db *sql.DB, profile *ConnectionProfile, llmModels []*LlmModel) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer database.RollBackOnErr(tx, err)

	query := "INSERT INTO connection_profiles (name, provider_type, base_url, api_key) VALUES ($1, $2, $3, $4) RETURNING id"
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

	return tx.Commit()
}

func UpdateConnectionProfile(db *sql.DB, id int64, profile *ConnectionProfile) error {
	query := `UPDATE connection_profiles
            SET name = $1,
                provider_type = $2,
                base_url = $3,
                api_key = $4
            WHERE id = $5`
	args := []any{profile.Name, profile.ProviderType, profile.BaseUrl, profile.ApiKey, id}

	return database.UpdateRecord(db, query, args)
}

func DeleteConnectionProfileById(db *sql.DB, id int64) error {
	query := "DELETE FROM connection_profiles WHERE id = $1"
	args := []any{id}

	return database.DeleteRecord(db, query, args)
}

func LlmModelsByConnectionProfileId(db *sql.DB, profileId int64) ([]*LlmModel, error) {
	query := "SELECT * FROM llm_models WHERE connection_profile_id = $1"
	args := []any{profileId}
	return database.QueryForList(db, query, args, llmModelScanner)
}

func CreateLlmModel(db *sql.DB, profileId int64, llmModel *LlmModel) error {
	return createLlmModel(db, profileId, llmModel)
}
func createLlmModel(db database.QueryExecutor, profileId int64, llmModel *LlmModel) error {
	llmModel.ConnectionProfileId = profileId

	query := `INSERT INTO llm_models
            (connection_profile_id, model_id, temperature, max_tokens, top_p, stream, stop_sequences, disabled)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
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

func UpdateLlmModel(db *sql.DB, id int64, llmModel *LlmModel) error {
	query := `UPDATE llm_models
              SET temperature = $1,
                  max_tokens = $2,
                  top_p = $3,
                  stream = $4,
                  stop_sequences = $5,
                  disabled = $6
              WHERE id = $7`
	args := []any{
		llmModel.Temperature,
		llmModel.MaxTokens,
		llmModel.TopP,
		llmModel.Stream,
		llmModel.StopSequences,
		llmModel.Disabled,
		id,
	}

	return database.UpdateRecord(db, query, args)
}

func DeleteLlmModelById(db *sql.DB, id int64) error {
	return deleteLlmModelById(db, id)
}

func deleteLlmModelById(db database.QueryExecutor, id int64) error {
	query := "DELETE FROM llm_models WHERE id = $1"
	args := []any{id}

	return database.DeleteRecord(db, query, args)
}

func MergeLlmModels(db *sql.DB, profileId int64, newModels []*LlmModel) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer database.RollBackOnErr(tx, err)

	// New model id set
	newModelIdSet := util.NewSetFrom(newModels, func(t *LlmModel) string {
		return t.ModelId
	})

	// Existing model id set
	existingModels, err := LlmModelsByConnectionProfileId(db, profileId)
	if err != nil {
		return err
	}

	existingModelIdSet := util.NewSetFrom(existingModels, func(t *LlmModel) string {
		return t.ModelId
	})

	// Add new models
	for _, newModel := range newModels {
		if existingModelIdSet.NotContains(newModel.ModelId) {
			if err = createLlmModel(tx, profileId, newModel); err != nil {
				return err
			}
		}
	}

	// Remove models not in new set
	for _, existingModel := range existingModels {
		if newModelIdSet.NotContains(existingModel.ModelId) {
			if err = deleteLlmModelById(tx, existingModel.ID); err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

func GetAllLlmModelViews(db database.QueryExecutor) ([]*LlmModelView, error) {
	query := `SELECT lm.id       AS model_id,
                   lm.model_id AS model_model_id,
                   p.id       AS profile_id,
                   p.name     AS profile_name
            FROM llm_models lm
                     JOIN connection_profiles p on p.id = lm.connection_profile_id
                     WHERE lm.disabled = FALSE`
	return database.QueryForList(db, query, nil, llModelViewScanner)
}
