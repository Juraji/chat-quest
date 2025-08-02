package model

import (
	"database/sql"
	"fmt"
	"strings"
)

type ConnectionProfile struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	ProviderType string `json:"providerType"`
	BaseUrl      string `json:"baseUrl"`
	ApiKey       string `json:"apiKey"`
}

type LlmModel struct {
	ID                  int64    `json:"id"`
	ConnectionProfileId int64    `json:"connectionProfileId"`
	ModelId             string   `json:"modelId"`
	Temperature         float64  `json:"temperature"`
	MaxTokens           int64    `json:"maxTokens"`
	TopP                float64  `json:"topP"`
	Stream              bool     `json:"stream"`
	Stop                []string `json:"stop"`
}

func connectionProfileScanner(scanner RowScanner, dest *ConnectionProfile) error {
	return scanner.Scan(
		&dest.ID,
		&dest.Name,
		&dest.ProviderType,
		&dest.BaseUrl,
		&dest.ApiKey,
	)
}

func llmModelScanner(scanner RowScanner, dest *LlmModel) error {
	var stopString string
	err := scanner.Scan(
		&dest.ID,
		&dest.ConnectionProfileId,
		&dest.ModelId,
		&dest.Temperature,
		&dest.MaxTokens,
		&dest.TopP,
		&dest.Stream,
		&stopString,
	)
	if err != nil {
		return err
	}

	if stopString != "" {
		dest.Stop = strings.Split(stopString, ",")
	} else {
		dest.Stop = []string{}
	}
	return nil
}

func AllConnectionProfiles(db *sql.DB) ([]*ConnectionProfile, error) {
	query := "SELECT * FROM connection_profiles"
	return queryForList(db, query, nil, connectionProfileScanner)
}

func ConnectionProfileById(db *sql.DB, id int64) (*ConnectionProfile, error) {
	query := "SELECT * FROM connection_profiles WHERE id = $1"
	args := []any{id}
	return queryForRecord(db, query, args, connectionProfileScanner)
}

func CreateConnectionProfile(db *sql.DB, profile *ConnectionProfile, llmModels []*LlmModel) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func(tx *sql.Tx, err error) {
		if err != nil {
			_ = tx.Rollback()
		}
	}(tx, err)

	query := "INSERT INTO connection_profiles (name, provider_type, base_url, api_key) VALUES ($1, $2, $3, $4) RETURNING id"
	args := []any{profile.Name, profile.ProviderType, profile.BaseUrl, profile.ApiKey}
	scanFunc := func(scanner RowScanner) error {
		return scanner.Scan(&profile.ID)
	}

	err = insertRecord(db, query, args, scanFunc)
	if err != nil {
		return err
	}

	for _, llmModel := range llmModels {
		err = CreateLlmModel(db, llmModel)
		if err != nil {
			return err
		}
	}

	return nil
}

func UpdateConnectionProfile(db *sql.DB, id int64, profile *ConnectionProfile) error {
	query := `UPDATE connection_profiles
            SET name = $1,
                provider_type = $2,
                base_url = $3,
                api_key = $4
            WHERE id = $5`
	args := []any{profile.Name, profile.ProviderType, profile.BaseUrl, profile.ApiKey, id}

	return updateRecord(db, query, args)
}

func DeleteConnectionProfileById(db *sql.DB, id int64) error {
	query := "DELETE FROM connection_profiles WHERE id = $1"
	args := []any{id}

	return deleteRecord(db, query, args)
}

func LlmModelsByConnectionProfileId(db *sql.DB, profileId int64) ([]*LlmModel, error) {
	query := "SELECT * FROM llm_models WHERE connection_profile_id = $1"
	args := []any{profileId}
	return queryForList(db, query, args, llmModelScanner)
}

func CreateLlmModel(db *sql.DB, llmModel *LlmModel) error {
	query := `INSERT INTO llm_models (connection_profile_id, model_id, temperature, max_tokens, top_p, stream, stop)
            VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	args := []any{
		llmModel.ConnectionProfileId,
		llmModel.ModelId,
		llmModel.Temperature,
		llmModel.MaxTokens,
		llmModel.TopP,
		llmModel.Stream,
		strings.Join(llmModel.Stop, ","),
	}
	scanFunc := func(scanner RowScanner) error {
		return scanner.Scan(&llmModel.ID)
	}

	return insertRecord(db, query, args, scanFunc)
}

func UpdateLlmModel(db *sql.DB, id int64, llmModel *LlmModel) error {
	query := `UPDATE llm_models
              SET temperature = $1,
                  max_tokens = $2,
                  top_p = $3,
                  stream = $4,
                  stop = $5
              WHERE id = $6`
	args := []any{
		llmModel.Temperature,
		llmModel.MaxTokens,
		llmModel.TopP,
		llmModel.Stream,
		strings.Join(llmModel.Stop, ","),
		id,
	}

	return updateRecord(db, query, args)
}

func DeleteLlmModelById(db *sql.DB, id int64) error {
	query := "DELETE FROM llm_models WHERE id = $1"
	args := []any{id}

	return deleteRecord(db, query, args)
}
