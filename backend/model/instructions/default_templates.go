package instructions

import (
	"embed"
	"strings"

	"github.com/pkg/errors"
	"juraji.nl/chat-quest/core/util"
)

//go:embed templates/*.tmpl
var templatesFs embed.FS

func newDefaultChatInstruction() (*Instruction, error) {
	instruction := &Instruction{
		Name:             "Default Chat",
		Type:             ChatInstruction,
		Temperature:      1.3,
		MaxTokens:        300,
		TopP:             0.95,
		PresencePenalty:  1.1,
		FrequencyPenalty: 1.1,
		Stream:           true,
		StopSequences:    nil,

		SystemPrompt: "templates/default_chat__system_prompt.tmpl",
		WorldSetup:   "templates/default_chat__world_setup.tmpl",
		Instruction:  "templates/default_chat__instruction.tmpl",
	}

	return reifyInstructionTemplates(instruction)
}

func newMultiCharResponseChatInstruction() (*Instruction, error) {
	instruction := &Instruction{
		Name:             "Multi-Character Response (experimental)",
		Type:             ChatInstruction,
		Temperature:      1.3,
		MaxTokens:        300,
		TopP:             0.95,
		PresencePenalty:  1.1,
		FrequencyPenalty: 1.1,
		Stream:           true,
		StopSequences:    nil,

		SystemPrompt: "templates/multi_char_response_chat__system_prompt.tmpl",
		WorldSetup:   "templates/multi_char_response_chat__world_setup.tmpl",
		Instruction:  "templates/multi_char_response_chat__instruction.tmpl",
	}

	return reifyInstructionTemplates(instruction)
}

func newDefaultMemoryInstruction() (*Instruction, error) {
	instruction := &Instruction{
		Name:             "Default Memories",
		Type:             MemoriesInstruction,
		Temperature:      0.9,
		MaxTokens:        300,
		TopP:             0.95,
		PresencePenalty:  1.1,
		FrequencyPenalty: 1.1,
		Stream:           false,
		StopSequences:    nil,

		SystemPrompt: "templates/default_memories__system_prompt.tmpl",
		WorldSetup:   "templates/default_memories__world_setup.tmpl",
		Instruction:  "templates/default_memories__instruction.tmpl",
	}

	return reifyInstructionTemplates(instruction)
}

func reifyInstructionTemplates(instruction *Instruction) (*Instruction, error) {
	props := []*string{
		&instruction.SystemPrompt,
		&instruction.WorldSetup,
		&instruction.Instruction,
	}

	for _, prop := range props {
		tpl, err := util.ReadFileAsString(templatesFs, *prop)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to load template from '%v'", prop)
		}

		if strings.HasPrefix("{{- /*gotype:", tpl) {
			// Remove IDE tpl vars declaration (first line of file)
			tpl = strings.SplitN(tpl, "\n", 2)[1]
		}

		*prop = tpl
	}

	return instruction, nil
}
