package chat_sessions

import (
	"go.uber.org/zap"
	"juraji.nl/chat-quest/core/database"
	"juraji.nl/chat-quest/core/log"
	"time"
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

func GetAllChatMessages(sessionId int) ([]ChatMessage, bool) {
	query := "SELECT * FROM chat_messages WHERE chat_session_id=?"
	args := []any{sessionId}
	list, err := database.QueryForList(query, args, ChatMessageScanner)
	if err != nil {
		log.Get().Error("Error fetching chat session messages",
			zap.Int("sessionId", sessionId),
			zap.Error(err))
		return nil, false
	}

	return list, true
}

func GetUnarchivedChatMessages(sessionId int) ([]ChatMessage, error) {
	query := "SELECT * FROM chat_messages WHERE chat_session_id=? AND is_archived=FALSE"
	args := []any{sessionId}
	return database.QueryForList(query, args, ChatMessageScanner)
}

func CreateChatMessage(sessionId int, chatMessage *ChatMessage) bool {
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
	if err != nil {
		log.Get().Error("Error creating chat session message",
			zap.Int("sessionId", sessionId),
			zap.Error(err))
		return false
	}

	ChatMessageCreatedSignal.EmitBG(chatMessage)
	return true
}

func UpdateChatMessage(sessionId int, id int, chatMessage *ChatMessage) bool {
	query := `UPDATE chat_messages
            SET content = ?,
                is_generating = ?
            WHERE chat_session_id = ?
              AND id = ?`
	args := []any{chatMessage.Content, chatMessage.IsGenerating, sessionId, id}

	err := database.UpdateRecord(query, args)
	if err != nil {
		log.Get().Error("Error updating chat session message",
			zap.Int("sessionId", sessionId),
			zap.Int("id", id),
			zap.Error(err))
		return false
	}

	ChatMessageUpdatedSignal.EmitBG(chatMessage)
	return true
}

func SetMessageArchived(sessionId int, id int) {
	query := `UPDATE chat_messages SET is_archived = TRUE WHERE chat_session_id = ? AND id = ?`
	args := []any{sessionId, id}
	err := database.UpdateRecord(query, args)
	if err != nil {
		log.Get().Fatal("Error archiving chat session message")
	}
}

func DeleteChatMessagesFrom(sessionId int, id int) bool {
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
	if err != nil {
		log.Get().Error("Error deleting chat session message",
			zap.Int("sessionId", sessionId),
			zap.Int("startingAtId", id),
			zap.Error(err))
		return false
	}

	ChatMessageDeletedSignal.EmitAllBG(deletedIds)
	return true
}
