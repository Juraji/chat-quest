package model

import (
	"database/sql"
	"juraji.nl/chat-quest/database"
)

type ChatPreferences struct {
	ChatModelID           *int64  `json:"chatModelId"`
	ChatInstructionID     *int64  `json:"chatInstructionId"`
	MemoriesModelID       *int64  `json:"memoriesModelId"`
	MemoriesInstructionID *int64  `json:"memoriesInstructionId"`
	EmbeddingModelID      *int64  `json:"embeddingModelId"`
	MemoryTopP            float64 `json:"memoryTopP"`
	MemoryTriggerAfter    int64   `json:"memoryTriggerAfter"`
	MemoryWindowSize      int64   `json:"memoryWindowSize"`
}

func chatPreferencesScanner(scanner database.RowScanner, dest *ChatPreferences) error {
	return scanner.Scan(
		dest.ChatModelID,
		dest.ChatInstructionID,
		dest.MemoriesModelID,
		dest.MemoriesInstructionID,
		dest.EmbeddingModelID,
		dest.MemoryTopP,
		dest.MemoryTriggerAfter,
		dest.MemoryWindowSize,
	)
}

func GetChatPreferences(db *sql.DB) (*ChatPreferences, error) {
	query := "SELECT * FROM chat_preferences WHERE id = 0"
	return database.QueryForRecord(db, query, nil, chatPreferencesScanner)
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
		prefs.EmbeddingModelID,
		prefs.MemoryTopP,
		prefs.MemoryTriggerAfter,
		prefs.MemoryWindowSize,
	}
	return database.UpdateRecord(db, query, args)
}
