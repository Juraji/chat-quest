package instructions

import (
	"context"
	"embed"
	"github.com/pkg/errors"
	"juraji.nl/chat-quest/core/database"
	"juraji.nl/chat-quest/core/util"
	"strings"
)

//go:embed templates/*.tmpl
var templatesFs embed.FS

func init() {
	const migrationsListenerKey = "InstallDefaultInstructions"
	const migrationsListenerExecAtVersion = 3

	database.MigrationsCompletedSignal.AddListener(migrationsListenerKey, func(ctx context.Context, event database.MigratedEvent) {
		if event.IsUpIncludingVersion(migrationsListenerExecAtVersion) {
			creators := map[string]func() (*InstructionTemplate, error){
				"Default Chat Instructions":     newChatInstructionTemplateFromDefault,
				"Default Memories Instructions": newMemoriesInstructionTemplateFromDefault,
			}

			for name, creator := range creators {
				template, err := creator()
				if err != nil {
					panic(err)
				}

				template.Name = name
				err = CreateInstruction(template)
				if err != nil {
					panic(errors.Wrap(err, "failed to save instruction template "+name))
				}
			}
		}
	})
}

func newChatInstructionTemplateFromDefault() (*InstructionTemplate, error) {
	const systemPromptTplPath = "templates/default_chat_instruction__system_prompt.tmpl"
	const worldSetupTplPath = "templates/default_chat_instruction__world_setup.tmpl"
	const userInstructionTplPath = "templates/default_chat_instruction__instruction.tmpl"

	return newInstructionTemplateFromDefault(
		ChatInstructionType,
		systemPromptTplPath,
		worldSetupTplPath,
		userInstructionTplPath,
	)
}

func newMemoriesInstructionTemplateFromDefault() (*InstructionTemplate, error) {
	const systemPromptTplPath = "templates/default_memories_instruction__system_prompt.tmpl"
	const worldSetupTplPath = "templates/default_memories_instruction__world_setup.tmpl"
	const userInstructionTplPath = "templates/default_memories_instruction__instruction.tmpl"

	return newInstructionTemplateFromDefault(
		MemoriesInstructionType,
		systemPromptTplPath,
		worldSetupTplPath,
		userInstructionTplPath,
	)
}

func newInstructionTemplateFromDefault(
	instructionType InstructionType,
	systemPromptTplPath string,
	worldSetupTplPath string,
	userInstructionTplPath string,
) (*InstructionTemplate, error) {
	systemPromptTpl, err := util.ReadFileAsString(templatesFs, systemPromptTplPath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load template from '%s'", systemPromptTplPath)
	}
	worldSetupTpl, err := util.ReadFileAsString(templatesFs, worldSetupTplPath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load template from '%s'", worldSetupTplPath)
	}
	instructionTpl, err := util.ReadFileAsString(templatesFs, userInstructionTplPath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load template from '%s'", userInstructionTplPath)
	}

	systemPromptTpl = strings.SplitN(systemPromptTpl, "\n", 2)[1]
	worldSetupTpl = strings.SplitN(worldSetupTpl, "\n", 2)[1]
	instructionTpl = strings.SplitN(instructionTpl, "\n", 2)[1]

	return &InstructionTemplate{
		Type:         instructionType,
		SystemPrompt: systemPromptTpl,
		WorldSetup:   worldSetupTpl,
		Instruction:  instructionTpl,
	}, nil
}
