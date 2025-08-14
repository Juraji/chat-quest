package memories

import (
	"juraji.nl/chat-quest/core/database"
	"juraji.nl/chat-quest/core/util"
	"time"
)

type Memory struct {
	ID               int        `json:"id"`
	WorldId          int        `json:"worldId"`
	ChatSessionId    int        `json:"chatSessionId"`
	CharacterId      int        `json:"characterId"`
	CreatedAt        *time.Time `json:"createdAt"`
	Content          string     `json:"content"`
	Embedding        util.Embeddings
	EmbeddingModelId *int
}

func (m *Memory) CosineSimilarity(other util.Embeddings) (float32, error) {
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

func GetMemoriesByWorldId(worldId int) ([]*Memory, error) {
	query := `SELECT id,
                   world_id,
                   chat_session_id,
                   character_id,
                   created_at,
                   content
            FROM memories
            WHERE world_id = ?`
	args := []interface{}{worldId}

	return database.QueryForList(database.GetDB(), query, args, memoryScanner)
}

func GetMemoriesByWorldAndCharacterId(
	worldId int,
	characterId int,
) ([]*Memory, error) {
	query := `SELECT id,
                   world_id,
                   chat_session_id,
                   character_id,
                   created_at,
                   content
            FROM memories
            WHERE world_id = ? AND character_id = ?`
	args := []interface{}{worldId, characterId}

	return database.QueryForList(database.GetDB(), query, args, memoryScanner)
}

func GetMemoriesByWorldAndCharacterIdWithEmbeddings(
	worldId int,
	characterId int,
) ([]*Memory, error) {
	query := `SELECT * FROM memories m WHERE world_id = ? AND character_id = ?`
	args := []interface{}{worldId, characterId}

	return database.QueryForList(database.GetDB(), query, args, memoryScanner)
}

func CreateMemory(memory *Memory) error {
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

	err := database.InsertRecord(database.GetDB(), query, args, &memory.ID)
	if err != nil {
		return err
	}

	util.Emit(MemoryCreatedSignal, memory)
	return nil
}

func UpdateMemory(id int, memory *Memory) error {
	query := `UPDATE memories SET content = ?, embedding = ?, embedding_model_id = ? WHERE id = ?`
	args := []any{memory.Content, memory.Embedding, memory.EmbeddingModelId, id}

	err := database.UpdateRecord(database.GetDB(), query, args)
	if err != nil {
		return err
	}

	util.Emit(MemoryUpdatedSignal, memory)
	return nil
}

func DeleteMemory(id int) error {
	query := `DELETE FROM memories WHERE id = ?`
	args := []any{id}

	err := database.DeleteRecord(database.GetDB(), query, args)
	if err != nil {
		return err
	}

	util.Emit(MemoryDeletedSignal, id)
	return nil
}

func GetMemoryPreferences() (*MemoryPreferences, error) {
	query := `SELECT memories_model_id,
                   memories_instruction_id,
                   embedding_model_id,
                   memory_min_p,
                   memory_trigger_after,
                   memory_window_size
            FROM memory_preferences
            WHERE id = 0`
	return database.QueryForRecord(database.GetDB(), query, nil, memoryPreferencesScanner)
}

func UpdateMemoryPreferences(prefs *MemoryPreferences) error {
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

	err := database.UpdateRecord(database.GetDB(), query, args)
	if err != nil {
		return err
	}

	util.Emit(MemoryPreferencesUpdatedSignal, prefs)
	return nil
}
