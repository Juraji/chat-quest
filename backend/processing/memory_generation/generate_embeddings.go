package memory_generation

import (
	"context"
	"go.uber.org/zap"
	"juraji.nl/chat-quest/core/log"
	m "juraji.nl/chat-quest/model/memories"
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

	embeddings, modelId, ok := m.GenerateEmbeddingForContent(memoryContent)
	if !ok {
		logger.Warn("Failed to generate embeddings")
		return
	}

	ok = m.SetMemoryEmbedding(memoryId, embeddings, modelId)
	if !ok {
		logger.Warn("Error setting memory embeddings")
		return
	}

	logger.Debug("Successfully generated embeddings")
}
