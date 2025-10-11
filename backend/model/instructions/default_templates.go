package instructions

import (
	"embed"
	"encoding/json"

	"github.com/pkg/errors"
)

//go:embed templates/*
var templatesFs embed.FS

var defaultTemplates = map[string]string{
	"default_chat":             "Default Chat",
	"default_memories":         "Default Memories",
	"multi_char_response_chat": "Multi-Character Response (experimental)",
	"npc_chat":                 "NPC Response (experimental)",
}

func reifyInstructionTemplate(templateName string) (*Instruction, error) {
	// Load instruction data
	instructionJsonData, err := templatesFs.ReadFile("templates/" + templateName + ".json")
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read instruction data '%s'", templateName)
	}
	var instruction Instruction
	err = json.Unmarshal(instructionJsonData, &instruction)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal instruction data for '%s'", templateName)
	}

	// Reify system prompt
	systemPromptPath := *instruction.SystemPrompt
	systemPromptData, err := templatesFs.ReadFile(systemPromptPath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load system prompt template from '%s'", systemPromptPath)
	}
	systemPrompt := string(systemPromptData)
	instruction.SystemPrompt = &systemPrompt

	// Reify world setup
	worldSetupPath := *instruction.WorldSetup
	worldSetupData, err := templatesFs.ReadFile(worldSetupPath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load world setup template from '%s'", worldSetupPath)
	}
	worldSetup := string(worldSetupData)
	instruction.WorldSetup = &worldSetup

	// Reify instruction
	instructionPath := instruction.Instruction
	instructionData, err := templatesFs.ReadFile(instructionPath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load instruction template from '%s'", instructionPath)
	}
	instruction.Instruction = string(instructionData)

	return &instruction, nil
}
