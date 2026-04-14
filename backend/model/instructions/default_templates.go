package instructions

import (
	"embed"
	"encoding/json"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

//go:embed templates/*
var templatesFs embed.FS

var defaultTemplates = sync.OnceValue(func() map[string]string {
	dirEntries, err := templatesFs.ReadDir("templates")
	if err != nil {
		panic(errors.Wrap(err, "failed to read default templates"))
	}

	tpls := make(map[string]string)
	for _, e := range dirEntries {
		if strings.HasSuffix(e.Name(), ".json") {
			key := strings.TrimSuffix(e.Name(), ".json")
			var sparseInstruction struct{ Name string }

			jsondata, err := templatesFs.ReadFile("templates/" + e.Name())
			if err != nil {
				panic(errors.Wrapf(err, "failed to read raw data for instruction '%s'", e.Name()))
			}
			err = json.Unmarshal(jsondata, &sparseInstruction)
			if err != nil {
				panic(errors.Wrapf(err, "failed to unmarshal JSON data for instruction '%s'", e.Name()))
			}

			tpls[key] = sparseInstruction.Name
		}
	}

	return tpls
})

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
	instruction.SystemPrompt = new(string(systemPromptData))

	// Reify world setup
	worldSetupPath := *instruction.WorldSetup
	worldSetupData, err := templatesFs.ReadFile(worldSetupPath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load world setup template from '%s'", worldSetupPath)
	}
	instruction.WorldSetup = new(string(worldSetupData))

	// Reify instruction
	instructionPath := instruction.Instruction
	instructionData, err := templatesFs.ReadFile(instructionPath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load instruction template from '%s'", instructionPath)
	}
	instruction.Instruction = string(instructionData)

	return &instruction, nil
}
