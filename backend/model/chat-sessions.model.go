package model

import (
	"database/sql"
	"time"
)

type ChatSession struct {
	ID         int64      `json:"id"`
	WorldID    int64      `json:"worldId"`
	CreatedAt  *time.Time `json:"createdAt"`
	Name       string     `json:"name"`
	ScenarioID int64      `json:"scenarioId"`

	ChatModelID       *int64 `json:"chatModelId"`
	ChatInstructionID *int64 `json:"chatInstructionId"`

	EnableMemories        bool   `json:"enableMemories"`
	MemoriesModelID       *int64 `json:"memoriesModelId"`
	MemoriesInstructionID *int64 `json:"memoriesInstructionId"`
}

type ChatMessage struct {
	ID            int64      `json:"id"`
	ChatSessionID int64      `json:"chatSessionId"`
	CreatedAt     *time.Time `json:"createdAt"`
	IsUser        bool       `json:"isUser"`
	CharacterID   *int64     `json:"characterId"`
	Content       string     `json:"content"`
}

func chatSessionScanner(scanner rowScanner, dest *ChatSession) error {
	return scanner.Scan(
		&dest.ID,
		&dest.WorldID,
		&dest.CreatedAt,
		&dest.Name,
		&dest.ScenarioID,
		&dest.ChatModelID,
		&dest.ChatInstructionID,
		&dest.EnableMemories,
		&dest.MemoriesModelID,
		&dest.MemoriesInstructionID,
	)
}

func chatMessageScanner(scanner rowScanner, dest *ChatMessage) error {
	return scanner.Scan(
		&dest.ID,
		&dest.ChatSessionID,
		&dest.CreatedAt,
		&dest.IsUser,
		&dest.CharacterID,
		&dest.Content,
	)
}

func GetAllChatSessionsByWorldId(db *sql.DB, worldId int64) ([]*ChatSession, error) {
	query := "SELECT * FROM chat_sessions WHERE world_id=$1"
	args := []any{worldId}

	return queryForList(db, query, args, chatSessionScanner)
}

func GetChatSessionById(db *sql.DB, worldId int64, id int64) (*ChatSession, error) {
	query := "SELECT * FROM chat_sessions WHERE world_id=$1 AND id=$2"
	args := []any{worldId, id}
	return queryForRecord(db, query, args, chatSessionScanner)
}

func CreateChatSession(db *sql.DB, worldId int64, chatSession *ChatSession) error {
	chatSession.WorldID = worldId
	chatSession.CreatedAt = nil

	query := `INSERT INTO chat_sessions (world_id, created_at, name, scenario_id, chat_model_id, chat_instruction_id, enable_memories, memories_model_id, memories_instruction_id)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`
	args := []any{
		chatSession.WorldID,
		chatSession.CreatedAt,
		chatSession.Name,
		chatSession.ScenarioID,
		chatSession.ChatModelID,
		chatSession.ChatInstructionID,
		chatSession.EnableMemories,
		chatSession.MemoriesModelID,
		chatSession.MemoriesInstructionID,
	}
	scanFunc := func(scanner rowScanner) error {
		return scanner.Scan(&chatSession.ID)
	}

	return insertRecord(db, query, args, scanFunc)
}

func UpdateChatSession(db *sql.DB, worldId int64, id int64, chatSession *ChatSession) error {
	query := `UPDATE chat_sessions
            SET name = $3,
                scenario_id = $5,
                chat_model_id = $5,
                chat_instruction_id = $6,
                enable_memories = $7,
                memories_model_id = $8,
                memories_instruction_id = $9
            WHERE world_id = $1
              AND id = $2`
	args := []any{
		worldId,
		id,
		chatSession.Name,
		chatSession.ScenarioID,
		chatSession.ChatModelID,
		chatSession.ChatInstructionID,
		chatSession.EnableMemories,
		chatSession.MemoriesModelID,
		chatSession.MemoriesInstructionID,
	}

	return updateRecord(db, query, args)
}

func DeleteChatSessionById(db *sql.DB, worldId int64, id int64) error {
	query := "DELETE FROM chat_sessions WHERE world_id=$1 AND id=$2"
	args := []any{worldId, id}
	return deleteRecord(db, query, args)
}

func GetChatMessages(db *sql.DB, sessionId int64) ([]*ChatMessage, error) {
	query := "SELECT * FROM chat_messages WHERE chat_session_id=$1"
	args := []any{sessionId}
	return queryForList(db, query, args, chatMessageScanner)
}

func CreateChatMessage(db *sql.DB, sessionId int64, chatMessage *ChatMessage) error {
	chatMessage.ChatSessionID = sessionId
	chatMessage.CreatedAt = nil

	query := `INSERT INTO chat_messages (chat_session_id, is_user, character_id, content)
            VALUES ($1, $2, $3, $4) RETURNING id`
	args := []any{
		chatMessage.ChatSessionID,
		chatMessage.IsUser,
		chatMessage.CharacterID,
		chatMessage.Content,
	}
	scanFunc := func(scanner rowScanner) error {
		return scanner.Scan(&chatMessage.ID)
	}

	return insertRecord(db, query, args, scanFunc)
}

func UpdateChatMessage(db *sql.DB, sessionId int64, id int64, chatMessage *ChatMessage) error {
	query := `UPDATE chat_messages
            SET content = $2
            WHERE chat_session_id = $1
              AND id = $3`
	args := []any{sessionId, id, chatMessage.Content}

	return updateRecord(db, query, args)
}

func DeleteChatMessagesFrom(db *sql.DB, sessionId int64, id int64) error {
	query := `DELETE FROM chat_messages
            WHERE chat_session_id = $1
              AND id >= $2`
	args := []any{sessionId, id}
	return deleteRecord(db, query, args)
}
