package providers

import (
	"juraji.nl/chat-quest/core/database"
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
	ID           int          `json:"id"`
	Name         string       `json:"name"`
	ProviderType ProviderType `json:"providerType"`
	BaseUrl      string       `json:"baseUrl"`
	ApiKey       string       `json:"apiKey"`
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

func AllConnectionProfiles() ([]ConnectionProfile, error) {
	query := "SELECT * FROM connection_profiles"
	return database.QueryForList(query, nil, connectionProfileScanner)
}

func ConnectionProfileById(id int) (*ConnectionProfile, error) {
	query := "SELECT * FROM connection_profiles WHERE id = ?"
	args := []any{id}
	return database.QueryForRecord(query, args, connectionProfileScanner)
}

func CreateConnectionProfile(profile *ConnectionProfile, llmModels []*LlmModel) error {
	err := database.Transactional(func(ctx *database.TxContext) error {
		query := "INSERT INTO connection_profiles (name, provider_type, base_url, api_key) VALUES (?, ?, ?, ?) RETURNING id"
		args := []any{profile.Name, profile.ProviderType, profile.BaseUrl, profile.ApiKey}

		if err := ctx.InsertRecord(query, args, &profile.ID); err != nil {
			return err
		}

		for _, llmModel := range llmModels {
			err := createLlmModel(ctx, profile.ID, llmModel)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err == nil {
		ConnectionProfileCreatedSignal.EmitBG(profile)
		LlmModelCreatedSignal.EmitAllBG(llmModels)
	}
	return err
}

func UpdateConnectionProfile(id int, profile *ConnectionProfile) error {
	query := `UPDATE connection_profiles
            SET name = ?,
                provider_type = ?,
                base_url = ?,
                api_key = ?
            WHERE id = ?`
	args := []any{profile.Name, profile.ProviderType, profile.BaseUrl, profile.ApiKey, id}

	err := database.UpdateRecord(query, args)

	if err == nil {
		ConnectionProfileUpdatedSignal.EmitBG(profile)
	}

	return err
}

func DeleteConnectionProfileById(id int) error {
	query := "DELETE FROM connection_profiles WHERE id = ?"
	args := []any{id}

	err := database.DeleteRecord(query, args)

	if err == nil {
		ConnectionProfileDeletedSignal.EmitBG(id)
	}

	return err
}
