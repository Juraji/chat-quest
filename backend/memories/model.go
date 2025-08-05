package memories

import (
	"database/sql"
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

func (m *Memory) CosineSimilarity(other providers.Embeddings) (float64, error) {
	return m.Embedding.CosineSimilarity(other)
}

func memoryScanner(scanner database.RowScanner, dest *Memory) error {
	return scanner.Scan(
		&dest.ID,
		&dest.ChatSessionId,
		&dest.CharacterId,
		&dest.CreatedAt,
		&dest.Content,
		&dest.Embedding,
		&dest.EmbeddingModelId,
	)
}

func GetMemoriesByWorldId(db *sql.DB, worldId int64) ([]*Memory, error) {
	query := `SELECT id,
                   world_id,
                   chat_session_id,
                   character_id,
                   created_at,
                   content
            FROM memories
            WHERE world_id = $1`
	args := []interface{}{worldId}

	return database.QueryForList(db, query, args, memoryScanner)
}

func GetMemoriesByWorldAndCharacterId(
	db *sql.DB,
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
            WHERE world_id = $1 AND character_id = $2`
	args := []interface{}{worldId, characterId}

	return database.QueryForList(db, query, args, memoryScanner)
}

func GetMemoriesByWorldAndCharacterIdWithEmbeddings(
	db *sql.DB,
	worldId int64,
	characterId int64,
) ([]*Memory, error) {
	query := `SELECT * FROM memories m WHERE world_id = $1 AND character_id = $2`
	args := []interface{}{worldId, characterId}

	return database.QueryForList(db, query, args, memoryScanner)
}

func CreateMemory(db *sql.DB, memory *Memory) error {
	query := `INSERT INTO memories (world_id, chat_session_id, character_id, created_at, content, embedding, embedding_model_id)
            VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	args := []any{
		memory.WorldId,
		memory.ChatSessionId,
		memory.CharacterId,
		memory.CreatedAt,
		memory.Content,
		memory.Embedding,
		memory.EmbeddingModelId,
	}
	scanFunc := func(scanner database.RowScanner) error {
		return scanner.Scan(&memory.ID)
	}
	return database.InsertRecord(db, query, args, scanFunc)
}

func UpdateMemory(db *sql.DB, id int64, memory *Memory) error {
	query := `UPDATE memories SET content = $2, embedding = $3, embedding_model_id = $4 WHERE id = $1`
	args := []any{id, memory.Content, memory.Embedding, memory.EmbeddingModelId}
	return database.UpdateRecord(db, query, args)
}

func DeleteMemory(db *sql.DB, id int64) error {
	query := `DELETE FROM memories WHERE id = $1`
	args := []any{id}
	return database.DeleteRecord(db, query, args)
}
