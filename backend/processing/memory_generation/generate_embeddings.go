package memory_generation

import (
	"context"
	"go.uber.org/zap"
	"juraji.nl/chat-quest/core/log"
	p "juraji.nl/chat-quest/core/providers"
	m "juraji.nl/chat-quest/model/memories"
	"juraji.nl/chat-quest/model/preferences"
)

func GenerateEmbeddings(ctx context.Context, memory *m.Memory) {
	if memory == nil {
		return
	}

	memoryId := memory.ID
	memoryContent := memory.Content
	logger := log.Get().With(
		zap.Int("memoryId", memoryId),
		zap.String("content", memoryContent))

	if ctx.Err() != nil {
		logger.Debug("Cancelled by context")
		return
	}

	prefs, err := preferences.GetPreferences(true)
	if err != nil {
		logger.Error("Error getting preferences", zap.Error(err))
		return
	}

	modelId := *prefs.EmbeddingModelId
	modelInstance, err := p.GetLlmModelInstanceById(modelId)
	if err != nil {
		logger.Warn("Error getting embedding model instance",
			zap.Int("modelId", modelId), zap.Error(err))
		return
	}

	embeddings, err := p.GenerateEmbeddings(modelInstance, memoryContent)
	if err != nil {
		logger.Error("Error generating embeddings", zap.Error(err))
		return
	}

	err = m.SetMemoryEmbedding(memoryId, embeddings, modelId)
	if err != nil {
		logger.Error("Error setting memory embeddings", zap.Error(err))
		return
	}

	logger.Debug("Successfully generated embeddings")
}
