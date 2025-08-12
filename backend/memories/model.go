package memories

import (
	"juraji.nl/chat-quest/cq"
	"juraji.nl/chat-quest/database"
	"juraji.nl/chat-quest/providers"
	"time"
)

type Memory struct {
	ID               int        `json:"id"`
	WorldId          int64      `json:"worldId"`
	ChatSessionId    int64      `json:"chatSessionId"`
	CharacterId      int64      `json:"characterId"`
	CreatedAt        *time.Time `json:"createdAt"`
	Content          string     `json:"content"`
	Embedding        providers.Embeddings
	EmbeddingModelId *int64
}

func (m *Memory) CosineSimilarity(other providers.Embeddings) (float32, error) {
	return m.Embedding.CosineSimilarity(other)
}

type MemoryPreferences struct {
	MemoriesModelID       *int64  `json:"memoriesModelId"`
	MemoriesInstructionID *int64  `json:"memoriesInstructionId"`
	EmbeddingModelID      *int64  `json:"embeddingModelId"`
	MemoryMinP            float64 `json:"memoryMinP"`
	MemoryTriggerAfter    int64   `json:"memoryTriggerAfter"`
	MemoryWindowSize      int64   `json:"memoryWindowSize"`
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

func GetMemoriesByWorldId(cq *cq.ChatQuestContext, worldId int64) ([]*Memory, error) {
	query := `SELECT id,
                   world_id,
                   chat_session_id,
                   character_id,
                   created_at,
                   content
            FROM memories
            WHERE world_id = ?`
	args := []interface{}{worldId}

	return database.QueryForList(cq.DB(), query, args, memoryScanner)
}

func GetMemoriesByWorldAndCharacterId(
	cq *cq.ChatQuestContext,
	worldId int64,
	characterId int64,
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

	return database.QueryForList(cq.DB(), query, args, memoryScanner)
}

func GetMemoriesByWorldAndCharacterIdWithEmbeddings(
	cq *cq.ChatQuestContext,
	worldId int64,
	characterId int64,
) ([]*Memory, error) {
	query := `SELECT * FROM memories m WHERE world_id = ? AND character_id = ?`
	args := []interface{}{worldId, characterId}

	return database.QueryForList(cq.DB(), query, args, memoryScanner)
}

func CreateMemory(cq *cq.ChatQuestContext, memory *Memory) error {
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

	err := database.InsertRecord(cq.DB(), query, args, &memory.ID)
	if err != nil {
		return err
	}

	MemoryCreatedSignal.Emit(cq.Context(), memory)
	return nil
}

func UpdateMemory(cq *cq.ChatQuestContext, id int64, memory *Memory) error {
	query := `UPDATE memories SET content = ?, embedding = ?, embedding_model_id = ? WHERE id = ?`
	args := []any{memory.Content, memory.Embedding, memory.EmbeddingModelId, id}

	err := database.UpdateRecord(cq.DB(), query, args)
	if err != nil {
		return err
	}

	MemoryUpdatedSignal.Emit(cq.Context(), memory)
	return nil
}

func DeleteMemory(cq *cq.ChatQuestContext, id int64) error {
	query := `DELETE FROM memories WHERE id = ?`
	args := []any{id}

	err := database.DeleteRecord(cq.DB(), query, args)
	if err != nil {
		return err
	}

	defer MemoryDeletedSignal.Emit(cq.Context(), id)
	return nil
}

func GetMemoryPreferences(cq *cq.ChatQuestContext) (*MemoryPreferences, error) {
	query := `SELECT memories_model_id,
                   memories_instruction_id,
                   embedding_model_id,
                   memory_min_p,
                   memory_trigger_after,
                   memory_window_size
            FROM memory_preferences
            WHERE id = 0`
	return database.QueryForRecord(cq.DB(), query, nil, memoryPreferencesScanner)
}

func UpdateMemoryPreferences(cq *cq.ChatQuestContext, prefs *MemoryPreferences) error {
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

	err := database.UpdateRecord(cq.DB(), query, args)
	if err != nil {
		return err
	}

	MemoryPreferencesUpdatedSignal.Emit(cq.Context(), prefs)
	return nil
}
