package memories

import (
	"go.uber.org/zap"
	"juraji.nl/chat-quest/core/database"
	"juraji.nl/chat-quest/core/log"
	"juraji.nl/chat-quest/core/providers"
	"time"
)

type Memory struct {
	ID               int        `json:"id"`
	WorldId          int        `json:"worldId"`
	ChatSessionId    int        `json:"chatSessionId"`
	CharacterId      int        `json:"characterId"`
	CreatedAt        *time.Time `json:"createdAt"`
	Content          string     `json:"content"`
	Embedding        providers.Embeddings
	EmbeddingModelId *int
}

func (m *Memory) CosineSimilarity(other providers.Embeddings) (float32, error) {
	return m.Embedding.CosineSimilarity(other)
}

type MemoryPreferences struct {
	MemoriesModelID       *int    `json:"memoriesModelId"`
	MemoriesInstructionID *int    `json:"memoriesInstructionId"`
	EmbeddingModelID      *int    `json:"embeddingModelId"`
	MemoryMinP            float32 `json:"memoryMinP"`
	MemoryTriggerAfter    int     `json:"memoryTriggerAfter"`
	MemoryWindowSize      int     `json:"memoryWindowSize"`
}

func memoryScanner(scanner database.RowScanner, dest *Memory) error {
	return scanner.Scan(
		&dest.ID,
		&dest.WorldId,
		&dest.ChatSessionId,
		&dest.CharacterId,
		&dest.CreatedAt,
		&dest.Content,
		&dest.Embedding,
		&dest.EmbeddingModelId,
	)
}

func memoryPreferencesScanner(scanner database.RowScanner, dest *MemoryPreferences) error {
	return scanner.Scan(
		&dest.MemoriesModelID,
		&dest.MemoriesInstructionID,
		&dest.EmbeddingModelID,
		&dest.MemoryMinP,
		&dest.MemoryTriggerAfter,
		&dest.MemoryWindowSize,
	)
}

func GetMemoriesByWorldId(worldId int) ([]Memory, bool) {
	query := `SELECT id,
                   world_id,
                   chat_session_id,
                   character_id,
                   created_at,
                   content
            FROM memories
            WHERE world_id = ?`
	args := []interface{}{worldId}
	list, err := database.QueryForList(query, args, memoryScanner)
	if err != nil {
		log.Get().Error("Error fetching memories", zap.Error(err))
		return nil, false
	}

	return list, true
}

func GetMemoriesByWorldAndCharacterId(
	worldId int,
	characterId int,
) ([]Memory, bool) {
	query := `SELECT id,
                   world_id,
                   chat_session_id,
                   character_id,
                   created_at,
                   content
            FROM memories
            WHERE world_id = ? AND character_id = ?`
	args := []interface{}{worldId, characterId}

	list, err := database.QueryForList(query, args, memoryScanner)
	if err != nil {
		log.Get().Error("Error fetching memories for character",
			zap.Int("characterId", characterId), zap.Error(err))
		return nil, false
	}

	return list, true
}

func GetMemoriesByWorldAndCharacterIdWithEmbeddings(
	worldId int,
	characterId int,
) ([]Memory, bool) {
	query := `SELECT * FROM memories m WHERE world_id = ? AND (character_id IS NULL OR character_id = ?)`
	args := []interface{}{worldId, characterId}

	list, err := database.QueryForList(query, args, memoryScanner)
	if err != nil {
		log.Get().Error("Error fetching memories (with embeddings) for character",
			zap.Int("characterId", characterId), zap.Error(err))
		return nil, false
	}

	return list, true
}

func CreateMemory(memory *Memory) bool {
	query := `INSERT INTO memories (world_id, chat_session_id, character_id, created_at, content, embedding, embedding_model_id)
            VALUES (?, ?, ?, ?, ?, ?, ?) RETURNING id`
	args := []any{
		memory.WorldId,
		memory.ChatSessionId,
		memory.CharacterId,
		memory.CreatedAt,
		memory.Content,
		memory.Embedding,
		memory.EmbeddingModelId,
	}

	err := database.InsertRecord(query, args, &memory.ID)
	if err != nil {
		log.Get().Error("Error creating memory", zap.Error(err))
		return false
	}

	MemoryCreatedSignal.EmitBG(memory)
	return true
}

func UpdateMemory(id int, memory *Memory) bool {
	query := `UPDATE memories SET content = ?, embedding = ?, embedding_model_id = ? WHERE id = ?`
	args := []any{memory.Content, memory.Embedding, memory.EmbeddingModelId, id}

	err := database.UpdateRecord(query, args)
	if err != nil {
		log.Get().Error("Error updating memory",
			zap.Int("id", id), zap.Error(err))
		return false
	}

	MemoryUpdatedSignal.EmitBG(memory)
	return true
}

func DeleteMemory(id int) bool {
	query := `DELETE FROM memories WHERE id = ?`
	args := []any{id}

	err := database.DeleteRecord(query, args)
	if err != nil {
		log.Get().Error("Error deleting memory",
			zap.Int("id", id), zap.Error(err))
		return false
	}

	MemoryDeletedSignal.EmitBG(id)
	return true
}

func GetMemoryPreferences() (*MemoryPreferences, bool) {
	query := `SELECT memories_model_id,
                   memories_instruction_id,
                   embedding_model_id,
                   memory_min_p,
                   memory_trigger_after,
                   memory_window_size
            FROM memory_preferences
            WHERE id = 0`
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
