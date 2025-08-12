package chat_sessions

import (
	"fmt"
	"juraji.nl/chat-quest/characters"
	"juraji.nl/chat-quest/cq"
	"juraji.nl/chat-quest/database"
	"juraji.nl/chat-quest/util"
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

type ChatParticipant struct {
	ChatSessionID int64 `json:"chatSessionId"`
	CharacterID   int64 `json:"characterId"`
}

func chatSessionScanner(scanner database.RowScanner, dest *ChatSession) error {
	return scanner.Scan(
		&dest.ID,
		&dest.WorldID,
		&dest.CreatedAt,
		&dest.Name,
		&dest.ScenarioID,
		&dest.EnableMemories,
	)
}

func chatMessageScanner(scanner database.RowScanner, dest *ChatMessage) error {
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

func GetAllChatSessionsByWorldId(cq *cq.ChatQuestContext, worldId int64) ([]*ChatSession, error) {
	query := "SELECT * FROM chat_sessions WHERE world_id=?"
	args := []any{worldId}

	return database.QueryForList(cq.DB(), query, args, chatSessionScanner)
}

func GetChatSessionById(cq *cq.ChatQuestContext, worldId int64, id int64) (*ChatSession, error) {
	query := "SELECT * FROM chat_sessions WHERE world_id=? AND id=?"
	args := []any{worldId, id}
	return database.QueryForRecord(cq.DB(), query, args, chatSessionScanner)
}

func CreateChatSession(cq *cq.ChatQuestContext, worldId int64, chatSession *ChatSession, characterIds []int64) error {
	chatSession.WorldID = worldId
	chatSession.CreatedAt = nil

	tx, err := cq.DB().Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer database.RollBackOnErr(tx, err)

	query := `INSERT INTO chat_sessions (world_id, name, scenario_id, enable_memories)
            VALUES (?, ?, ?, ?) RETURNING id, created_at`
	args := []any{
		chatSession.WorldID,
		chatSession.Name,
		chatSession.ScenarioID,
		chatSession.EnableMemories,
	}

	err = database.InsertRecord(tx, query, args, &chatSession.ID, &chatSession.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert chat session: %w", err)
	}

	for _, characterId := range characterIds {
		err = addChatSessionParticipant(tx, chatSession.ID, characterId)
		if err != nil {
			return fmt.Errorf("failed to insert chat participant (%d -> %d):  %w", chatSession.ID, characterId, err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	ChatSessionCreatedSignal.Emit(cq.Context(), chatSession)
	for _, characterId := range characterIds {
		participant := ChatParticipant{chatSession.ID, characterId}
		ChatParticipantAddedSignal.Emit(cq.Context(), &participant)
	}

	return nil
}

func UpdateChatSession(cq *cq.ChatQuestContext, worldId int64, id int64, chatSession *ChatSession) error {
	query := `UPDATE chat_sessions
            SET name = ?,
                scenario_id = ?,
                enable_memories = ?
            WHERE world_id = ?
              AND id = ?`
	args := []any{
		chatSession.Name,
		chatSession.ScenarioID,
		chatSession.EnableMemories,
		worldId,
		id,
	}

	err := database.UpdateRecord(cq.DB(), query, args)
	if err != nil {
		return err
	}

	ChatSessionUpdatedSignal.Emit(cq.Context(), chatSession)
	return nil
}

func DeleteChatSessionById(cq *cq.ChatQuestContext, worldId int64, id int64) error {
	query := "DELETE FROM chat_sessions WHERE world_id=? AND id=?"
	args := []any{worldId, id}

	err := database.DeleteRecord(cq.DB(), query, args)
	if err != nil {
		return err
	}

	ChatSessionDeletedSignal.Emit(cq.Context(), worldId)
	return nil
}

func GetChatMessages(cq *cq.ChatQuestContext, sessionId int64) ([]*ChatMessage, error) {
	query := "SELECT * FROM chat_messages WHERE chat_session_id=?"
	args := []any{sessionId}
	return database.QueryForList(cq.DB(), query, args, chatMessageScanner)
}

func CreateChatMessage(cq *cq.ChatQuestContext, sessionId int64, chatMessage *ChatMessage) error {
	chatMessage.ChatSessionID = sessionId
	chatMessage.CreatedAt = nil

	query := `INSERT INTO chat_messages (chat_session_id, is_user, character_id, content)
            VALUES (?, ?, ?, ?) RETURNING id, created_at`
	args := []any{
		chatMessage.ChatSessionID,
		chatMessage.IsUser,
		chatMessage.CharacterID,
		chatMessage.Content,
	}

	err := database.InsertRecord(cq.DB(), query, args, &chatMessage.ID, &chatMessage.CreatedAt)
	if err != nil {
		return err
	}

	ChatMessageCreatedSignal.Emit(cq.Context(), chatMessage)

	return nil
}

func UpdateChatMessage(cq *cq.ChatQuestContext, sessionId int64, id int64, chatMessage *ChatMessage) error {
	query := `UPDATE chat_messages
            SET content = ?
            WHERE chat_session_id = ?
              AND id = ?`
	args := []any{chatMessage.Content, sessionId, id}

	err := database.UpdateRecord(cq.DB(), query, args)
	if err != nil {
		return err
	}

	ChatMessageUpdatedSignal.Emit(cq.Context(), chatMessage)
	return nil
}

func DeleteChatMessagesFrom(cq *cq.ChatQuestContext, sessionId int64, id int64) error {
	//language=SQL
	query := `DELETE
            FROM chat_messages
            WHERE chat_session_id = ?
              AND id >= ?
            RETURNING id`
	args := []any{sessionId, id}

	deletedIds, err := database.QueryForList(cq.DB(), query, args, func(scanner database.RowScanner, dest *int64) error {
		return scanner.Scan(dest)
	})
	if err != nil {
		return err
	}

	util.EmitAllNonNil(cq, ChatMessageDeletedSignal, deletedIds)
	return nil
}

func GetChatSessionParticipants(cq *cq.ChatQuestContext, chatSessionId int64) ([]*characters.Character, error) {
	query := `SELECT c.* FROM chat_participants cp
                JOIN characters c ON cp.character_id = c.id
            WHERE cp.chat_session_id = ?`
	args := []any{chatSessionId}
	return database.QueryForList(cq.DB(), query, args, characters.CharacterScanner)
}

func AddChatSessionParticipant(cq *cq.ChatQuestContext, chatSessionId int64, characterId int64) error {
	err := addChatSessionParticipant(cq.DB(), chatSessionId, characterId)
	if err != nil {
		return err
	}

	participant := ChatParticipant{chatSessionId, characterId}
	ChatParticipantAddedSignal.Emit(cq.Context(), &participant)
	return nil
}

func addChatSessionParticipant(db database.QueryExecutor, chatSessionId int64, characterId int64) error {
	query := `INSERT INTO chat_participants (chat_session_id, character_id) VALUES (?, ?)`
	args := []any{chatSessionId, characterId}
	return database.InsertRecord(db, query, args)
}

func RemoveChatSessionParticipant(cq *cq.ChatQuestContext, chatSessionId int64, characterId int64) error {
	query := `DELETE FROM chat_participants WHERE chat_session_id = ? AND character_id = ?`
	args := []any{chatSessionId, characterId}

	err := database.DeleteRecord(cq.DB(), query, args)
	if err != nil {
		return err
	}

	participant := ChatParticipant{chatSessionId, characterId}
	ChatParticipantRemovedSignal.Emit(cq.Context(), &participant)
	return nil
}
