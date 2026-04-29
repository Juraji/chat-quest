package processing

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"juraji.nl/chat-quest/core/log"
	p "juraji.nl/chat-quest/core/providers"
	m "juraji.nl/chat-quest/model/memories"
	"juraji.nl/chat-quest/model/preferences"
)

func RegenerateEmbeddingsOnPrefsUpdate(ctx context.Context, prefs *preferences.Preferences) error {
	logger := log.Get().With(zap.Intp("embeddingModelId", prefs.EmbeddingModelId))

	logger.Info("Preferences updated, checking memory embeddings...")

	memories, err := m.GetMemoriesNotMatchingEmbeddingModelId(*prefs.EmbeddingModelId)
	if err != nil {
		logger.Error("Error fetching memories to regenerate", zap.Error(err))
		return errors.Wrap(err, "error fetching memories to regenerate")
	}
	if len(memories) == 0 {
		logger.Info("No memories to regenerate")
		return nil
	}

	logger.Info("Updating memories...", zap.Int("memoryCount", len(memories)))

	if contextCheckPoint(ctx, logger) {
		return nil
	}

	for _, memory := range memories {
		err := GenerateEmbeddings(ctx, &memory)
		if err != nil {
			return err
		}
		if contextCheckPoint(ctx, logger) {
			return nil
		}
	}

	logger.Info("Embedding updated")
	return nil
}

func GenerateEmbeddings(ctx context.Context, memory *m.Memory) error {
	if memory == nil {
		return nil
	}

	memoryId := memory.ID
	memoryContent := memory.Content
	logger := log.Get().With(zap.Int("memoryId", memoryId))

	logger.Info("Generating embeddings for memory")

	if contextCheckPoint(ctx, logger) {
		return nil
	}

	prefs, err := preferences.GetPreferences(true)
	if err != nil {
		logger.Error("Error getting preferences", zap.Error(err))
		return errors.Wrap(err, "error getting preferences")
	}

	modelId := *prefs.EmbeddingModelId
	modelInstance, err := p.GetLlmModelInstanceById(modelId)
	if err != nil {
		logger.Warn("Error getting embedding model instance",
			zap.Int("modelId", modelId), zap.Error(err))
		return errors.Wrap(err, "error getting embedding model instance")
	}

	embeddings, err := p.GenerateEmbeddings(modelInstance, memoryContent)
	if err != nil {
		logger.Error("Error generating embeddings", zap.Error(err))
		return errors.Wrap(err, "error generating embeddings")
	}

	err = m.SetMemoryEmbedding(memoryId, embeddings, modelId)
	if err != nil {
		logger.Error("Error setting memory embeddings", zap.Error(err))
		return errors.Wrap(err, "error setting memory embeddings")
	}

	logger.Debug("Successfully generated embeddings")
	return nil
}
