package instructions

import (
	"github.com/pkg/errors"
	"juraji.nl/chat-quest/core/util"
	"strings"
)

func ApplyInstructionTemplates(instruction InstructionTemplate, variables any) (*InstructionTemplate, error) {
	fields := []*string{
		&instruction.SystemPrompt,
		&instruction.WorldSetup,
		&instruction.Instruction,
	}

	errChan := make(chan error, len(fields))
	defer close(errChan)

	for _, fieldPtr := range fields {
		go func() {
			if util.HasTemplateVars(*fieldPtr) {
				tpl, err := util.NewTextTemplate("Template", *fieldPtr)
				if err != nil {
					errChan <- errors.Wrap(err, "Error creating template for instruction template")
					return
				}

				*fieldPtr = strings.TrimSpace(util.WriteToString(tpl, variables))
			}
			errChan <- nil
		}()
	}

	for i := 0; i < len(fields); i++ {
		err := <-errChan
		if err != nil {
			return nil, errors.Wrap(err, "Error processing instruction template")
		}
	}
	return &instruction, nil
}
