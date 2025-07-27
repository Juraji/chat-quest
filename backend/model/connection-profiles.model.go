package model

import (
	"database/sql"
	"strings"
)

type ConnectionProfile struct {
	ID           int64  `json:"id"`
	ProviderType string `json:"providerType"`
	BaseUrl      string `json:"baseUrl"`
	ApiKey       string `json:"apiKey"`
}

type LlmModel struct {
	ID                  int64    `json:"id"`
	ConnectionProfileId string   `json:"connectionProfileId"`
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

func CreateConnectionProfile(db *sql.DB, profile *ConnectionProfile) error {
	query := "INSERT INTO connection_profiles (provider_type, base_url, api_key) VALUES ($1, $2, $3) RETURNING id"
	args := []any{profile.ProviderType, profile.BaseUrl, profile.ApiKey}
	scanFunc := func(scanner RowScanner) error {
		return scanner.Scan(&profile.ID)
	}

	return insertRecord(db, query, args, scanFunc)
}

func UpdateConnectionProfile(db *sql.DB, id int64, profile *ConnectionProfile) error {
	query := "UPDATE connection_profiles SET provider_type = $1, base_url = $2, api_key = $3 WHERE id = $4"
	args := []any{profile.ProviderType, profile.BaseUrl, profile.ApiKey, id}

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
	query := "INSERT INTO llm_models (model_id, connection_profile_id, temperature, max_tokens, top_p, stream, stop) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"
	args := []any{
		llmModel.ModelId,
		llmModel.ConnectionProfileId,
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
	query := "UPDATE llm_models SET temperature = $1, max_tokens = $2, top_p = $3, stream = $4, stop = $5 WHERE id = $6"
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
