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
	IsSystem      bool       `json:"isSystem"`
	IsGenerating  bool       `json:"isGenerating"`
	CharacterID   *int       `json:"characterId"`
	Content       string     `json:"content"`

	// Managed by memories
	MemoryID *int `json:"memoryId"`
}

func chatMessageScanner(scanner database.RowScanner, dest *ChatMessage) error {
	return scanner.Scan(
		&dest.ID,
		&dest.ChatSessionID,
		&dest.CreatedAt,
		&dest.IsUser,
		&dest.IsSystem,
		&dest.IsGenerating,
		&dest.CharacterID,
		&dest.Content,
		&dest.MemoryID,
	)
}

func NewChatMessage(isUser bool, isSystem bool, isGenerating bool, characterId *int, content string) *ChatMessage {
	return &ChatMessage{
		ID:            0,
		ChatSessionID: 0,
		CreatedAt:     nil,
		IsUser:        isUser,
		IsSystem:      isSystem,
		IsGenerating:  isGenerating,
		CharacterID:   characterId,
		Content:       content,
		MemoryID:      nil,
	}
}

func GetChatMessages(sessionId int) ([]ChatMessage, bool) {
	query := "SELECT * FROM chat_messages WHERE chat_session_id=?"
	args := []any{sessionId}
	list, err := database.QueryForList(query, args, chatMessageScanner)
	if err != nil {
		log.Get().Error("Error fetching chat session messages",
			zap.Int("sessionId", sessionId),
			zap.Error(err))
		return nil, false
	}

	return list, true
}

func GetChatMessagesPreceding(sessionId int, messageId int) ([]ChatMessage, bool) {
	query := "SELECT * FROM chat_messages WHERE chat_session_id=? and id < ?"
	args := []any{sessionId, messageId}
	list, err := database.QueryForList(query, args, chatMessageScanner)
	if err != nil {
		log.Get().Error("Error fetching preceding chat session messages",
			zap.Int("sessionId", sessionId),
			zap.Int("precedingId", messageId),
			zap.Error(err))
		return nil, false
	}

	return list, true
}

func CreateChatMessage(sessionId int, chatMessage *ChatMessage) bool {
	chatMessage.ChatSessionID = sessionId
	chatMessage.CreatedAt = nil

	query := `INSERT INTO chat_messages (chat_session_id, is_user, is_system, is_generating, character_id, content)
            VALUES (?, ?, ?, ?, ?, ?) RETURNING id, created_at`
	args := []any{
		chatMessage.ChatSessionID,
		chatMessage.IsUser,
		chatMessage.IsSystem,
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
