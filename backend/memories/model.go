package memories

import (
	"juraji.nl/chat-quest/cq"
	"juraji.nl/chat-quest/database"
	"juraji.nl/chat-quest/util"
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

func GetMemoriesByWorldId(cq *cq.ChatQuestContext, worldId int) ([]*Memory, error) {
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

	return database.QueryForList(cq.DB(), query, args, memoryScanner)
}

func GetMemoriesByWorldAndCharacterIdWithEmbeddings(
	cq *cq.ChatQuestContext,
	worldId int,
	characterId int,
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

func UpdateMemory(cq *cq.ChatQuestContext, id int, memory *Memory) error {
	query := `UPDATE memories SET content = ?, embedding = ?, embedding_model_id = ? WHERE id = ?`
	args := []any{memory.Content, memory.Embedding, memory.EmbeddingModelId, id}

	err := database.UpdateRecord(cq.DB(), query, args)
	if err != nil {
		return err
	}

	MemoryUpdatedSignal.Emit(cq.Context(), memory)
	return nil
}

func DeleteMemory(cq *cq.ChatQuestContext, id int) error {
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
