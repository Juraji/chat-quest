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

	prefs, err := preferences.GetPreferences()
	if err != nil {
		logger.Error("Failed to get preferences", zap.Error(err))
		return
	}
	if err := prefs.Validate(); err != nil {
		logger.Error("Error validating memory preferences", zap.Error(err))
		return
	}

	modelId := *prefs.EmbeddingModelId
	modelInstance, ok := p.GetLlmModelInstanceById(modelId)
	if !ok {
		logger.Warn("Error getting embedding model instance",
			zap.Int("modelId", modelId))
		return
	}

	embeddings, err := p.GenerateEmbeddings(modelInstance, memoryContent)
	if err != nil {
		logger.Error("Error generating embeddings", zap.Error(err))
		return
	}

	ok = m.SetMemoryEmbedding(memoryId, embeddings, modelId)
	if !ok {
		logger.Warn("Error setting memory embeddings")
		return
	}

	logger.Debug("Successfully generated embeddings")
}
