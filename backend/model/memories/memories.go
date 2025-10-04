package memories

import (
	"time"

	"juraji.nl/chat-quest/core/database"
	"juraji.nl/chat-quest/core/providers"
)

type Memory struct {
	ID               int                 `json:"id"`
	WorldId          int                 `json:"worldId"`
	CharacterId      *int                `json:"characterId"`
	CreatedAt        *time.Time          `json:"createdAt"`
	Content          string              `json:"content"`
	AlwaysInclude    bool                `json:"alwaysInclude"`
	Embedding        providers.Embedding `json:"-"`
	EmbeddingModelId *int                `json:"-"`
}

type GenerationRequest struct {
	BaseMessageId     int
	IncludeNPreceding int
}

func memoryScanner(scanner database.RowScanner, dest *Memory) error {
	return scanner.Scan(
		&dest.ID,
		&dest.WorldId,
		&dest.CharacterId,
		&dest.CreatedAt,
		&dest.Content,
		&dest.AlwaysInclude,
	)
}

func memoryWithEmbeddingsScanner(scanner database.RowScanner, dest *Memory) error {
	return scanner.Scan(
		&dest.ID,
		&dest.WorldId,
		&dest.CharacterId,
		&dest.CreatedAt,
		&dest.Content,
		&dest.AlwaysInclude,
		&dest.Embedding,
		&dest.EmbeddingModelId,
	)
}

func GetMemoriesByWorldId(worldId int) ([]Memory, error) {
	query := `SELECT id, world_id, character_id, created_at, content, always_include
            FROM memories
            WHERE world_id = ?`
	args := []any{worldId}
	return database.QueryForList(query, args, memoryScanner)
}

func GetMemoriesByWorldAndCharacterId(
	worldId int,
	characterId int,
) ([]Memory, error) {
	query := `SELECT id, world_id, character_id, created_at, content, always_include
				FROM memories
            	WHERE world_id = ? AND character_id = ?`
	args := []any{worldId, characterId}

	return database.QueryForList(query, args, memoryScanner)
}

func GetMemoriesByWorldAndCharacterIdWithEmbeddings(
	worldId int,
	characterId int,
	modelId int,
) ([]Memory, error) {
	query := `SELECT *
              FROM memories m
              WHERE world_id = ?
                AND embedding IS NOT NULL
                AND embedding_model_id = ?
                AND (character_id IS NULL OR character_id = ?)`
	args := []any{worldId, modelId, characterId}

	return database.QueryForList(query, args, memoryWithEmbeddingsScanner)
}

func GetMemoriesNotMatchingEmbeddingModelId(modelId int) ([]Memory, error) {
	query := `SELECT id, world_id, character_id, created_at, content, always_include
			  FROM memories
			  WHERE embedding_model_id IS NULL
			     OR embedding_model_id != ?`
	args := []any{modelId}
	return database.QueryForList(query, args, memoryScanner)
}

func CreateMemory(worldId int, memory *Memory) error {
	memory.WorldId = worldId

	query := `INSERT INTO memories (world_id, character_id, content, always_include)
            VALUES (?, ?, ?, ?) RETURNING id, created_at`
	args := []any{
		memory.WorldId,
		memory.CharacterId,
		memory.Content,
		memory.AlwaysInclude,
	}

	err := database.InsertRecord(query, args, &memory.ID, &memory.CreatedAt)

	if err == nil {
		MemoryCreatedSignal.EmitBG(memory)
	}

	return err
}

func UpdateMemory(id int, memory *Memory) error {
	query := `UPDATE memories
			  SET content = ?,
			      character_id = ?,
			      always_include = ?
			  WHERE id = ?`
	args := []any{memory.Content, memory.CharacterId, memory.AlwaysInclude, id}

	err := database.UpdateRecord(query, args)

	if err == nil {
		MemoryUpdatedSignal.EmitBG(memory)
	}

	return err
}

func SetMemoryEmbedding(id int, embeddings providers.Embedding, embeddingModelId int) error {
	query := `UPDATE memories SET embedding = ?, embedding_model_id = ? WHERE id = ?`
	args := []any{embeddings, embeddingModelId, id}
	return database.UpdateRecord(query, args)
}

func DeleteMemory(id int) error {
	query := `DELETE FROM memories WHERE id = ?`
	args := []any{id}

	_, err := database.DeleteRecord(query, args)

	if err == nil {
		MemoryDeletedSignal.EmitBG(id)
	}

	return err
}
