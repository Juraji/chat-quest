package processing

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"juraji.nl/chat-quest/core/log"
	c "juraji.nl/chat-quest/model/characters"
	i "juraji.nl/chat-quest/model/instructions"
	p "juraji.nl/chat-quest/model/preferences"
)

func ExportCharacterAsText(
	ctx context.Context,
	characterID int,
	instructionID int,
) (string, error) {
	var err error
	logger := log.Get().With(
		zap.Int("characterID", characterID),
		zap.Int("instructionID", instructionID))

	if contextCheckPoint(ctx, logger) {
		return "", nil
	}

	logger.Info("Exporting character as text...")
	character, err := c.CharacterById(characterID)
	if err != nil {
		logger.Error("Failed to get character", zap.Error(err))
		return "", errors.WithMessage(err, "failed to fetch character")
	}

	prefs, err := p.GetPreferences(true)
	if err != nil {
		logger.Error("Failed to get preferences", zap.Error(err))
		return "", errors.WithMessage(err, "failed to fetch preferences")
	}

	instruction, err := i.InstructionById(instructionID)
	if err != nil {
		logger.Error("Failed to get instruction", zap.Error(err))
		return "", errors.WithMessage(err, "failed to fetch instruction")
	}

	if contextCheckPoint(ctx, logger) {
		return "", nil
	}

	instructionVars := NewTemplateCharacter(character, prefs, nil)
	if err = instruction.ApplyTemplates(instructionVars); err != nil {
		logger.Error("Failed to apply templates", zap.Error(err))
		return "", errors.WithMessage(err, "failed to apply templates")
	}

	logger.Info("Character exported completed...")
	return instruction.Instruction, nil
}
