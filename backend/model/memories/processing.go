package memories

import (
	"go.uber.org/zap"
	"juraji.nl/chat-quest/core/database"
	"juraji.nl/chat-quest/core/log"
	"juraji.nl/chat-quest/core/providers"
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

// GenerateEmbeddingForContent uses the embedding model from the memory preferences to
// embed the content. It returns the generated providers.Embeddings and the model ID used.
func GenerateEmbeddingForContent(content string) (providers.Embeddings, int, bool) {
	logger := log.Get()

	prefs, ok := GetMemoryPreferences()
	if !ok {
		logger.Warn("Could not get memory preferences")
		return nil, 0, false
	}
	err := prefs.Validate()
	if err != nil {
		logger.Error("Error validating preferences", zap.Error(err))
		return nil, 0, false
	}

	modelId := *prefs.EmbeddingModelID
	modelInstance, ok := providers.GetLlmModelInstanceById(modelId)
	if !ok {
		logger.Warn("Error getting embedding model instance",
			zap.Intp("modelId", prefs.EmbeddingModelID))
		return nil, 0, false
	}

	embeddings, err := providers.GenerateEmbeddings(modelInstance, content)
	if err != nil {
		logger.Error("Error generating embeddings", zap.Error(err))
		return nil, 0, false
	}

	return embeddings, modelId, true
}
