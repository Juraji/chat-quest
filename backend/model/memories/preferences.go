package memories

import (
	"errors"
	"go.uber.org/zap"
	"juraji.nl/chat-quest/core/database"
	"juraji.nl/chat-quest/core/log"
)

type MemoryPreferences struct {
	MemoriesModelID       *int    `json:"memoriesModelId"`
	MemoriesInstructionID *int    `json:"memoriesInstructionId"`
	EmbeddingModelID      *int    `json:"embeddingModelId"`
	MemoryMinP            float32 `json:"memoryMinP"`
	MemoryTriggerAfter    int     `json:"memoryTriggerAfter"`
	MemoryWindowSize      int     `json:"memoryWindowSize"`
}

func (p *MemoryPreferences) Validate() error {
	if p == nil {
		return errors.New("memory preferences is nil")
	}
	if p.MemoriesModelID == nil {
		return errors.New("memories model not set")
	}
	if p.MemoriesInstructionID == nil {
		return errors.New("memories instruction not set")
	}
	if p.EmbeddingModelID == nil {
		return errors.New("embedding model not set")
	}

	return nil
}

func memoryPreferencesScanner(scanner database.RowScanner, dest *MemoryPreferences) error {
	var idSink int
	return scanner.Scan(
		&idSink,
		&dest.MemoriesModelID,
		&dest.MemoriesInstructionID,
		&dest.EmbeddingModelID,
		&dest.MemoryMinP,
		&dest.MemoryTriggerAfter,
		&dest.MemoryWindowSize,
	)
}

func GetMemoryPreferences() (*MemoryPreferences, bool) {
	query := "SELECT * FROM memory_preferences WHERE id = 0"
	prefs, err := database.QueryForRecord(query, nil, memoryPreferencesScanner)
	if err != nil {
		log.Get().Error("Error fetching preferences", zap.Error(err))
		return nil, false
	}

	return prefs, true
}

func UpdateMemoryPreferences(prefs *MemoryPreferences) bool {
	query := `UPDATE memory_preferences
            SET memories_model_id = ?,
                memories_instruction_id = ?,
                embedding_model_id = ?,
                memory_min_p = ?,
                memory_trigger_after = ?,
                memory_window_size = ?
            WHERE id = 0`
	args := []any{
		prefs.MemoriesModelID,
		prefs.MemoriesInstructionID,
		prefs.EmbeddingModelID,
		prefs.MemoryMinP,
		prefs.MemoryTriggerAfter,
		prefs.MemoryWindowSize,
	}

	err := database.UpdateRecord(query, args)
	if err != nil {
		log.Get().Error("Error updating preferences", zap.Error(err))
		return false
	}

	MemoryPreferencesUpdatedSignal.EmitBG(prefs)
	return true
}
