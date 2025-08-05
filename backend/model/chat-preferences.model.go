package model

import "database/sql"

type ChatPreferences struct {
	ChatModelID           *int64 `json:"chatModelId"`
	ChatInstructionID     *int64 `json:"chatInstructionId"`
	MemoriesModelID       *int64 `json:"memoriesModelId"`
	MemoriesInstructionID *int64 `json:"memoriesInstructionId"`
	EmbeddingModelID      *int64 `json:"embeddingModelId"`
}

func chatPreferencesScanner(scanner rowScanner, dest *ChatPreferences) error {
	return scanner.Scan(
		dest.ChatModelID,
		dest.ChatInstructionID,
		dest.MemoriesModelID,
		dest.MemoriesInstructionID,
		dest.EmbeddingModelID,
	)
}

func GetChatPreferences(db *sql.DB) (*ChatPreferences, error) {
	query := "SELECT * FROM chat_preferences WHERE id = 0"
	return queryForRecord(db, query, nil, chatPreferencesScanner)
}

func UpdateChatPreferences(db *sql.DB, prefs *ChatPreferences) error {
	query := `UPDATE chat_preferences
            SET chat_model_id = $1,
                chat_instruction_id = $2,
                memories_model_id = $3,
                memories_instruction_id = $4,
                embedding_model_id = $5
            WHERE id = 0`
	args := []any{
		prefs.ChatModelID,
		prefs.ChatInstructionID,
		prefs.MemoriesModelID,
		prefs.MemoriesInstructionID,
		prefs.EmbeddingModelID}
	return updateRecord(db, query, args)
}
