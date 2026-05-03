package processing

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"juraji.nl/chat-quest/core/log"
	p "juraji.nl/chat-quest/core/providers"
	c "juraji.nl/chat-quest/model/characters"
	i "juraji.nl/chat-quest/model/instructions"
	w "juraji.nl/chat-quest/model/worlds"
)

const charactersResponseFormat = `{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "required": [
    "age",
    "pronouns",
    "appearance",
    "personality",
    "history"
  ],
  "properties": {
    "age": {"type": "integer"},
    "pronouns": {"type": "string"},
    "appearance": {"type": "string"},
    "personality": {"type": "string"},
    "history": {"type": "string"}
  }
}`

type CharacterBuilderRequest struct {
	Character     *c.Character `json:"character"`
	Description   string       `json:"description"`
	WorldId       *int         `json:"worldId"`
	InstructionID int          `json:"instructionId"`
	LlmModelId    int          `json:"llmModelId"`
}

func BuildCharacter(
	ctx context.Context,
	request *CharacterBuilderRequest,
) (*c.Character, error) {
	logger := log.Get()

	// Cancellation
	ctx, cleanup := setupCancelBySystem(ctx, logger, "GenerateMemories")
	defer cleanup()

	var err error
	character := request.Character

	var world *w.World
	if request.WorldId != nil {
		if world, err = w.WorldById(*request.WorldId); err != nil {
			logger.Error("Error fetching world", zap.Error(err))
			return nil, errors.Wrap(err, "error fetching world")
		}
	}

	instruction, err := i.InstructionById(request.InstructionID)
	if err != nil {
		logger.Error("Failed to get instruction", zap.Error(err))
		return nil, errors.Wrap(err, "failed to fetch instruction")
	}

	modelInstance, err := p.GetLlmModelInstanceById(request.LlmModelId)
	if err != nil {
		logger.Error("Could not fetch model", zap.Error(err))
		return nil, errors.Wrap(err, "could not fetch model")
	}

	templateVars := NewCharacterBuilderVars(character, world, request.Description)
	if err = instruction.ApplyTemplates(templateVars); err != nil {
		logger.Error("Error applying instruction templates", zap.Error(err))
		return nil, errors.Wrap(err, "error applying instruction templates")
	}

	logInstructionsToFile(logger, instruction, nil)

	requestMessages := createChatRequestMessages(nil, instruction)
	llmParameters := instruction.AsLlmParameters()
	llmParameters.ResponseFormat = new(charactersResponseFormat)

	chatResponseChan := p.GenerateChatResponse(ctx, modelInstance, requestMessages, llmParameters)
	var rawResponse string

responseLoop:
	for {
		select {
		case r, hasNext := <-chatResponseChan:
			if r.Error != nil {
				logger.Error("Error in repsonse",
					zap.String("generated", rawResponse),
					zap.Error(r.Error))
				return nil, errors.Wrap(err, "error in response")
			}

			rawResponse = rawResponse + r.Content
			if !hasNext {
				// Done
				break responseLoop
			}
		case <-ctx.Done():
			logger.Debug("Cancelled by context")
			return nil, nil
		}
	}

	err = json.Unmarshal([]byte(rawResponse), &character)
	if err != nil {
		logger.Error("Could not unmarshal response",
			zap.String("response", rawResponse),
			zap.Error(err))
	}

	return character, nil
}
