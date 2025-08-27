package preferences

import (
	"github.com/pkg/errors"
	"juraji.nl/chat-quest/core/database"
)

type Preferences struct {
	// Chat
	ChatModelId       *int `json:"chatModelId"`
	ChatInstructionId *int `json:"chatInstructionId"`
	// Embedding
	EmbeddingModelId *int `json:"embeddingModelId"`
	// Memories
	MemoriesModelId       *int    `json:"memoriesModelId"`
	MemoriesInstructionId *int    `json:"memoriesInstructionId"`
	MemoryMinP            float32 `json:"memoryMinP"`
	MemoryTriggerAfter    int     `json:"memoryTriggerAfter"`
	MemoryWindowSize      int     `json:"memoryWindowSize"`
}

func (p *Preferences) Validate() []string {
	if p == nil {
		return []string{"preferences is nil"}
	}

	var errs []string
	if p.ChatModelId == nil {
		errs = append(errs, "chat model not set")
	}
	if p.ChatInstructionId == nil {
		errs = append(errs, "chat instruction not set")
	}

	if p.EmbeddingModelId == nil {
		errs = append(errs, "embedding model not set")
	}

	if p.MemoriesModelId == nil {
		errs = append(errs, "memories model not set")
	}
	if p.MemoriesInstructionId == nil {
		errs = append(errs, "memories instruction not set")
	}

	return errs
}

func preferencesScanner(scanner database.RowScanner, dest *Preferences) error {
	var idSink int
	return scanner.Scan(
		&idSink,
		&dest.ChatModelId,
		&dest.ChatInstructionId,
		&dest.EmbeddingModelId,
		&dest.MemoriesModelId,
		&dest.MemoriesInstructionId,
		&dest.MemoryMinP,
		&dest.MemoryTriggerAfter,
		&dest.MemoryWindowSize,
	)
}

func GetPreferences(validate bool) (*Preferences, error) {
	query := "SELECT * FROM preferences WHERE id = 0"
	prefs, err := database.QueryForRecord(query, nil, preferencesScanner)

	if err != nil {
		return nil, err
	}
	if validate {
		if errs := prefs.Validate(); len(errs) > 0 {
			return nil, errors.Errorf("preferences invalid: %v", errs)
		}
	}

	return prefs, nil
}

func UpdatePreferences(prefs *Preferences) error {
	query := `UPDATE preferences
             SET chat_model_id = ?,
                 chat_instruction_id = ?,
                 embedding_model_id = ?,
                 memories_model_id = ?,
                 memories_instruction_id = ?,
                 memory_min_p = ?,
                 memory_trigger_after = ?,
                 memory_window_size = ?
             WHERE id = 0`
	args := []any{
		prefs.ChatModelId,
		prefs.ChatInstructionId,
		prefs.EmbeddingModelId,
		prefs.MemoriesModelId,
		prefs.MemoriesInstructionId,
		prefs.MemoryMinP,
		prefs.MemoryTriggerAfter,
		prefs.MemoryWindowSize,
	}

	if err := database.UpdateRecord(query, args); err != nil {
		return err
	}
	PreferencesUpdatedSignal.EmitBG(prefs)
	return nil
}
