package instructions

import (
	"context"

	"github.com/pkg/errors"
	"juraji.nl/chat-quest/core/database"
)

func init() {
	const key = "InstallDefaultInstructions"
	const execAtVersion = 3

	database.MigrationsCompletedSignal.AddListener(key, func(ctx context.Context, event database.MigratedEvent) {
		if event.IsUpIncludingVersion(execAtVersion) {
			for tpl := range defaultTemplates {
				template, err := reifyInstructionTemplate(tpl)
				if err != nil {
					panic(err)
				}

				err = CreateInstruction(template)
				if err != nil {
					panic(errors.Wrap(err, "failed to save instruction template "+template.Name))
				}
			}
		}
	})
}
