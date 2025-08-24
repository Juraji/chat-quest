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

func SetProcessedStateForMessages(messages []chatsessions.ChatMessage) bool {
	query := `UPDATE chat_messages SET processed_by_memories = TRUE WHERE id = ?;`
	err := database.Transactional(func(ctx *database.TxContext) error {
		for _, message := range messages {
			args := []any{message.ID}
			err := ctx.UpdateRecord(query, args)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		log.Get().Error("Error updating processed state for memorized messages", zap.Error(err))
		return false
	}

	return true
}
