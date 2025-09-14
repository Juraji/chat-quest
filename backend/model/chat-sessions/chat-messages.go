package chat_sessions

import (
	"time"

	"juraji.nl/chat-quest/core/database"
)

type ChatMessage struct {
	ID            int        `json:"id"`
	ChatSessionID int        `json:"chatSessionId"`
	CreatedAt     *time.Time `json:"createdAt"`
	IsUser        bool       `json:"isUser"`
	IsGenerating  bool       `json:"isGenerating"`
	IsArchived    bool       `json:"isArchived"`
	CharacterID   *int       `json:"characterId"`
	Content       string     `json:"content"`
}

func ChatMessageScanner(scanner database.RowScanner, dest *ChatMessage) error {
	return scanner.Scan(
		&dest.ID,
		&dest.ChatSessionID,
		&dest.CreatedAt,
		&dest.IsUser,
		&dest.IsGenerating,
		&dest.IsArchived,
		&dest.CharacterID,
		&dest.Content,
	)
}

func NewChatMessage(isUser bool, isGenerating bool, characterId *int, content string) *ChatMessage {
	return &ChatMessage{
		ID:            0,
		ChatSessionID: 0,
		CreatedAt:     nil,
		IsUser:        isUser,
		IsGenerating:  isGenerating,
		IsArchived:    false,
		CharacterID:   characterId,
		Content:       content,
	}
}

func GetAllChatMessages(sessionId int) ([]ChatMessage, error) {
	query := "SELECT * FROM chat_messages WHERE chat_session_id=?"
	args := []any{sessionId}
	return database.QueryForList(query, args, ChatMessageScanner)
}

func GetUnarchivedChatMessages(sessionId int) ([]ChatMessage, error) {
	query := "SELECT * FROM chat_messages WHERE chat_session_id=? AND is_archived=FALSE"
	args := []any{sessionId}
	return database.QueryForList(query, args, ChatMessageScanner)
}

func GetMessageById(messageId int) (*ChatMessage, error) {
	query := "SELECT * FROM chat_messages WHERE id=?"
	args := []any{messageId}
	return database.QueryForRecord(query, args, ChatMessageScanner)
}

func CreateChatMessage(sessionId int, chatMessage *ChatMessage) error {
	chatMessage.ChatSessionID = sessionId
	chatMessage.CreatedAt = nil

	query := `INSERT INTO chat_messages (chat_session_id, is_user, is_generating, character_id, content)
            VALUES (?, ?, ?, ?, ?) RETURNING id, created_at`
	args := []any{
		chatMessage.ChatSessionID,
		chatMessage.IsUser,
		chatMessage.IsGenerating,
		chatMessage.CharacterID,
		chatMessage.Content,
	}

	err := database.InsertRecord(query, args, &chatMessage.ID, &chatMessage.CreatedAt)

	if err == nil {
		ChatMessageCreatedSignal.EmitBG(chatMessage)
	}

	return err
}

func UpdateChatMessage(sessionId int, id int, chatMessage *ChatMessage) error {
	query := `UPDATE chat_messages
            SET content = ?,
                is_generating = ?
            WHERE chat_session_id = ?
              AND id = ?`
	args := []any{chatMessage.Content, chatMessage.IsGenerating, sessionId, id}

	err := database.UpdateRecord(query, args)

	if err == nil {
		ChatMessageUpdatedSignal.EmitBG(chatMessage)
	}

	return err
}

func SetMessageArchived(sessionId int, id int) error {
	query := `UPDATE chat_messages SET is_archived = TRUE WHERE chat_session_id = ? AND id = ?`
	args := []any{sessionId, id}
	err := database.UpdateRecord(query, args)

	if err == nil {
		ChatMessageArchivedSignal.EmitBG(id)
	}

	return err
}

func DeleteChatMessagesFrom(sessionId int, id int) error {
	//language=SQL
	query := `DELETE
            FROM chat_messages
            WHERE chat_session_id = ?
              AND id >= ?
            RETURNING id`
	args := []any{sessionId, id}

	deletedIds, err := database.QueryForList(query, args, func(scanner database.RowScanner, dest *int) error {
		return scanner.Scan(dest)
	})

	if err == nil {
		ChatMessageDeletedSignal.EmitAllBG(deletedIds)
	}

	return err
}
