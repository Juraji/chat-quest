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
	CharacterId      int        `json:"characterId"`
	CreatedAt        *time.Time `json:"createdAt"`
	Content          string     `json:"content"`
	Embedding        providers.Embeddings
	EmbeddingModelId *int
}

func (m *Memory) CosineSimilarity(other providers.Embeddings) (float32, error) {
	return m.Embedding.CosineSimilarity(other)
}

func memoryScanner(scanner database.RowScanner, dest *Memory) error {
	return scanner.Scan(
		&dest.ID,
		&dest.WorldId,
		&dest.CharacterId,
		&dest.CreatedAt,
		&dest.Content,
		&dest.Embedding,
		&dest.EmbeddingModelId,
	)
}

func GetMemoriesByWorldId(worldId int) ([]Memory, bool) {
	query := `SELECT id,
                   world_id,
                   character_id,
                   created_at,
                   content
            FROM memories
            WHERE world_id = ?`
	args := []any{worldId}
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
                   character_id,
                   created_at,
                   content
            FROM memories
            WHERE world_id = ? AND character_id = ?`
	args := []any{worldId, characterId}

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
	query := `SELECT *
              FROM memories m
              WHERE world_id = ?
                AND embedding IS NOT NULL
                And embedding IS NOT NULL
                AND (character_id IS NULL OR character_id = ?)`
	args := []any{worldId, characterId}

	list, err := database.QueryForList(query, args, memoryScanner)
	if err != nil {
		log.Get().Error("Error fetching memories (with embeddings) for character",
			zap.Int("characterId", characterId), zap.Error(err))
		return nil, false
	}

	return list, true
}

func CreateMemory(worldId int, memory *Memory) bool {
	memory.WorldId = worldId

	query := `INSERT INTO memories (world_id, character_id, content)
            VALUES (?, ?, ?) RETURNING id`
	args := []any{
		memory.WorldId,
		memory.CharacterId,
		memory.Content,
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
	query := `UPDATE memories SET content = ? WHERE id = ?`
	args := []any{memory.Content, id}

	err := database.UpdateRecord(query, args)
	if err != nil {
		log.Get().Error("Error updating memory",
			zap.Int("id", id), zap.Error(err))
		return false
	}

	MemoryUpdatedSignal.EmitBG(memory)
	return true
}

func SetMemoryEmbedding(id int, embeddings providers.Embeddings, embeddingModelId int) bool {
	query := `UPDATE memories SET embedding = ?, embedding_model_id = ? WHERE id = ?`
	args := []any{embeddings, embeddingModelId, id}
	err := database.UpdateRecord(query, args)
	if err != nil {
		log.Get().Error("Error updating memory",
			zap.Int("id", id), zap.Error(err))
		return false
	}

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
