package memories

import (
	"go.uber.org/zap"
	"juraji.nl/chat-quest/core/database"
	"juraji.nl/chat-quest/core/log"
	chatsessions "juraji.nl/chat-quest/model/chat-sessions"
)

func GetUnprocessedMessagesForSession(sessionID int) ([]chatsessions.ChatMessage, bool) {
	query := `SELECT *
              FROM chat_messages m
              WHERE m.chat_session_id = ?
                AND m.processed_by_memories = FALSE;`
	args := []any{sessionID}

	list, err := database.QueryForList(query, args, chatsessions.ChatMessageScanner)
	if err != nil {
		log.Get().Error("Error fetching chat session messages for memories",
			zap.Int("sessionID", sessionID),
			zap.Error(err))
		return nil, false
	}

	return list, true
}
