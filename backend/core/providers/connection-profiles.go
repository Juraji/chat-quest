package providers

import (
	"go.uber.org/zap"
	"juraji.nl/chat-quest/core/database"
	"juraji.nl/chat-quest/core/log"
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

func AllConnectionProfiles() ([]ConnectionProfile, bool) {
	query := "SELECT * FROM connection_profiles"
	list, err := database.QueryForList(query, nil, connectionProfileScanner)
	if err != nil {
		log.Get().Error("Error querying for connection profiles", zap.Error(err))
		return list, false
	}

	return list, true
}

func ConnectionProfileById(id int) (*ConnectionProfile, bool) {
	query := "SELECT * FROM connection_profiles WHERE id = ?"
	args := []any{id}
	p, err := database.QueryForRecord(query, args, connectionProfileScanner)
	if err != nil {
		log.Get().Error("Error querying for connection profile", zap.Int("profileId", id), zap.Error(err))
		return nil, false
	}
	return p, true
}

func CreateConnectionProfile(profile *ConnectionProfile, llmModels []*LlmModel) bool {
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

	if err != nil {
		log.Get().Error("Error creating connection profile", zap.Error(err))
		return false
	}

	ConnectionProfileCreatedSignal.EmitBG(profile)
	LlmModelCreatedSignal.EmitAllBG(llmModels)
	return true
}

func UpdateConnectionProfile(id int, profile *ConnectionProfile) bool {
	query := `UPDATE connection_profiles
            SET name = ?,
                provider_type = ?,
                base_url = ?,
                api_key = ?
            WHERE id = ?`
	args := []any{profile.Name, profile.ProviderType, profile.BaseUrl, profile.ApiKey, id}

	err := database.UpdateRecord(query, args)
	if err != nil {
		log.Get().Error("Error updating connection profile", zap.Int("profileId", id), zap.Error(err))
		return false
	}

	ConnectionProfileUpdatedSignal.EmitBG(profile)
	return true
}

func DeleteConnectionProfileById(id int) bool {
	query := "DELETE FROM connection_profiles WHERE id = ?"
	args := []any{id}

	err := database.DeleteRecord(query, args)
	if err != nil {
		log.Get().Error("Error deleting connection profile", zap.Int("profileId", id), zap.Error(err))
		return false
	}

	ConnectionProfileDeletedSignal.EmitBG(id)
	return true
}
