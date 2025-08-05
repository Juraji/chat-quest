package model

import (
  "database/sql"
  "time"
)

type ChatSession struct {
  ID             int64      `json:"id"`
  WorldID        int64      `json:"worldId"`
  CreatedAt      *time.Time `json:"createdAt"`
  Name           string     `json:"name"`
  ScenarioID     *int64     `json:"scenarioId"`
  EnableMemories bool       `json:"enableMemories"`
}

type ChatMessage struct {
  ID            int64      `json:"id"`
  ChatSessionID int64      `json:"chatSessionId"`
  CreatedAt     *time.Time `json:"createdAt"`
  IsUser        bool       `json:"isUser"`
  CharacterID   *int64     `json:"characterId"`
  Content       string     `json:"content"`

  // Managed by memories
  MemoryID *int64 `json:"memoryId"`
}

func chatSessionScanner(scanner rowScanner, dest *ChatSession) error {
  return scanner.Scan(
    &dest.ID,
    &dest.WorldID,
    &dest.CreatedAt,
    &dest.Name,
    &dest.ScenarioID,
    &dest.EnableMemories,
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
    &dest.MemoryID,
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

  query := `INSERT INTO chat_sessions (world_id, created_at, name, scenario_id, enable_memories)
            VALUES ($1, $2, $3, $4, $5) RETURNING id`
  args := []any{
    chatSession.WorldID,
    chatSession.CreatedAt,
    chatSession.Name,
    chatSession.ScenarioID,
    chatSession.EnableMemories,
  }
  scanFunc := func(scanner rowScanner) error {
    return scanner.Scan(&chatSession.ID)
  }

  return insertRecord(db, query, args, scanFunc)
}

func UpdateChatSession(db *sql.DB, worldId int64, id int64, chatSession *ChatSession) error {
  query := `UPDATE chat_sessions
            SET name = $3,
                scenario_id = $4,
                enable_memories = $5
            WHERE world_id = $1
              AND id = $2`
  args := []any{
    worldId,
    id,
    chatSession.Name,
    chatSession.ScenarioID,
    chatSession.EnableMemories,
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
