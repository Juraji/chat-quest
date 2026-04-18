package chat_sessions

import (
	"slices"
	"time"

	"juraji.nl/chat-quest/core/database"
)

type ChatMessage struct {
	ID            int        `json:"id"`
	ChatSessionID int        `json:"chatSessionId"`
	CreatedAt     *time.Time `json:"createdAt"`
	IsUser        bool       `json:"isUser"`
	IsGenerating  bool       `json:"isGenerating"`
	CharacterID   *int       `json:"characterId"`
	Content       string     `json:"content"`
	Reasoning     string     `json:"reasoning"`
}

func ChatMessageScanner(scanner database.RowScanner, dest *ChatMessage) error {
	return scanner.Scan(
		&dest.ID,
		&dest.ChatSessionID,
		&dest.CreatedAt,
		&dest.IsUser,
		&dest.IsGenerating,
		&dest.CharacterID,
		&dest.Content,
		&dest.Reasoning,
	)
}

func NewChatMessage(isUser bool, isGenerating bool, characterId *int, content string) *ChatMessage {
	return &ChatMessage{
		IsUser:       isUser,
		IsGenerating: isGenerating,
		CharacterID:  characterId,
		Content:      content,
	}
}

func GetAllChatMessages(sessionId int) ([]ChatMessage, error) {
	query := "SELECT * FROM chat_messages WHERE chat_session_id=? ORDER BY created_at"
	args := []any{sessionId}
	return database.QueryForList(query, args, ChatMessageScanner)
}

func GetTailChatMessages(sessionId int, limit int) ([]ChatMessage, error) {
	query := "SELECT * FROM chat_messages WHERE chat_session_id=? ORDER BY id DESC LIMIT ?"
	args := []any{sessionId, limit}
	list, err := database.QueryForList(query, args, ChatMessageScanner)

	if err != nil {
		return nil, err
	}

	// Reverse the slice to be in ASC order
	slices.Reverse(list)

	return list, nil
}

func GetChatSessionMessageCount(sessionId int) (int, error) {
	query := "SELECT COUNT(*) FROM chat_messages WHERE chat_session_id=?"
	args := []any{sessionId}
	res, err := database.QueryForRecord(query, args, database.IntScanner)
	if err != nil {
		return 0, err
	}
	return *res, nil
}

func GetMessageById(messageId int) (*ChatMessage, error) {
	query := "SELECT * FROM chat_messages WHERE id=?"
	args := []any{messageId}
	return database.QueryForRecord(query, args, ChatMessageScanner)
}

func GetMessagesInSessionBeforeId(sessionId int, messageId int, limit int) ([]ChatMessage, error) {
	query := "SELECT * FROM chat_messages WHERE chat_session_id=? AND id<? ORDER BY id DESC LIMIT ?"
	args := []any{sessionId, messageId, limit}
	list, err := database.QueryForList(query, args, ChatMessageScanner)

	if err != nil {
		return nil, err
	}

	// Reverse the slice to be in ASC order
	slices.Reverse(list)

	return list, nil
}

func GetMessagesInSessionAfterId(sessionId int, messageId int, limit int) ([]ChatMessage, error) {
	query := "SELECT * FROM chat_messages WHERE chat_session_id=? AND id>? ORDER BY id LIMIT ?"
	args := []any{sessionId, messageId, limit}

	return database.QueryForList(query, args, ChatMessageScanner)
}

func CreateChatMessage(sessionId int, chatMessage *ChatMessage) error {
	chatMessage.ChatSessionID = sessionId
	chatMessage.CreatedAt = nil
	chatMessage.IsUser = chatMessage.CharacterID == nil

	query := `INSERT INTO chat_messages (chat_session_id, is_user, is_generating, character_id, content, reasoning)
            VALUES (?, ?, ?, ?, ?, ?) RETURNING id, created_at`
	args := []any{
		chatMessage.ChatSessionID,
		chatMessage.IsUser,
		chatMessage.IsGenerating,
		chatMessage.CharacterID,
		chatMessage.Content,
		chatMessage.Reasoning,
	}

	err := database.InsertRecord(query, args, &chatMessage.ID, &chatMessage.CreatedAt)

	if err == nil {
		ChatMessageCreatedSignal.EmitBG(chatMessage)
	}

	return err
}

func UpdateChatMessage(sessionId int, id int, chatMessage *ChatMessage) error {
	chatMessage.ChatSessionID = sessionId
	chatMessage.IsUser = chatMessage.CharacterID == nil

	query := `UPDATE chat_messages
            SET is_user = ?,
                is_generating = ?,
                character_id = ?,
                content = ?,
                reasoning = ?
            WHERE chat_session_id = ?
              AND id = ?`
	args := []any{
		chatMessage.IsUser,
		chatMessage.IsGenerating,
		chatMessage.CharacterID,
		chatMessage.Content,
		chatMessage.Reasoning,
		sessionId,
		id}

	err := database.UpdateRecord(query, args)

	if err == nil {
		ChatMessageUpdatedSignal.EmitBG(chatMessage)
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
