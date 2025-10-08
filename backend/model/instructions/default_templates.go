package instructions

import (
	"embed"
	"strings"

	"github.com/pkg/errors"
	"juraji.nl/chat-quest/core/util"
)

//go:embed templates/*.tmpl
var templatesFs embed.FS

var defaultInstructionTemplates = []*Instruction{
	{
		Name:             "Default Chat",
		Type:             ChatInstruction,
		Temperature:      1.1,
		MaxTokens:        300,
		TopP:             0.95,
		PresencePenalty:  1.1,
		FrequencyPenalty: 1.1,
		Stream:           true,
		StopSequences:    nil,

		ReasoningPrefix:   "<think>",
		ReasoningSuffix:   "</think>",
		CharacterIdPrefix: "<characterid>",
		CharacterIdSuffix: "</characterid>",

		SystemPrompt: util.StrAsPointer("templates/default_chat__system_prompt.tmpl"),
		WorldSetup:   util.StrAsPointer("templates/default_chat__world_setup.tmpl"),
		Instruction:  "templates/default_chat__instruction.tmpl",
	},
	{
		Name:             "Multi-Character Response (experimental)",
		Type:             ChatInstruction,
		Temperature:      1.1,
		MaxTokens:        300,
		TopP:             0.95,
		PresencePenalty:  1.1,
		FrequencyPenalty: 1.1,
		Stream:           true,
		StopSequences:    nil,

		ReasoningPrefix:   "<think>",
		ReasoningSuffix:   "</think>",
		CharacterIdPrefix: "<characterid>",
		CharacterIdSuffix: "</characterid>",

		SystemPrompt: util.StrAsPointer("templates/multi_char_response_chat__system_prompt.tmpl"),
		WorldSetup:   util.StrAsPointer("templates/multi_char_response_chat__world_setup.tmpl"),
		Instruction:  "templates/multi_char_response_chat__instruction.tmpl",
	},
	{
		Name:             "NPC Response (experimental)",
		Type:             ChatInstruction,
		Temperature:      1.1,
		MaxTokens:        300,
		TopP:             0.95,
		PresencePenalty:  1.1,
		FrequencyPenalty: 1.1,
		Stream:           true,
		StopSequences:    nil,

		ReasoningPrefix:   "<think>",
		ReasoningSuffix:   "</think>",
		CharacterIdPrefix: "<characterid>",
		CharacterIdSuffix: "</characterid>",

		SystemPrompt: util.StrAsPointer("templates/npc_chat__system_prompt.tmpl"),
		WorldSetup:   util.StrAsPointer("templates/npc_chat__world_setup.tmpl"),
		Instruction:  "templates/npc_chat__instruction.tmpl",
	},
	{
		Name:             "Default Memories",
		Type:             MemoriesInstruction,
		Temperature:      0.7,
		MaxTokens:        300,
		TopP:             0.95,
		PresencePenalty:  1.1,
		FrequencyPenalty: 1.1,
		Stream:           false,
		StopSequences:    nil,

		ReasoningPrefix:   "<think>",
		ReasoningSuffix:   "</think>",
		CharacterIdPrefix: "<characterid>",
		CharacterIdSuffix: "</characterid>",

		SystemPrompt: util.StrAsPointer("templates/default_memories__system_prompt.tmpl"),
		WorldSetup:   util.StrAsPointer("templates/default_memories__world_setup.tmpl"),
		Instruction:  "templates/default_memories__instruction.tmpl",
	},
}

func reifyInstructionTemplate(instruction *Instruction) (*Instruction, error) {
	props := []*string{
		instruction.SystemPrompt,
		instruction.WorldSetup,
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
