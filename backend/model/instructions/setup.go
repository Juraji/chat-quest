package instructions

import (
	"context"

	"github.com/pkg/errors"
	"juraji.nl/chat-quest/core/database"
	"juraji.nl/chat-quest/core/log"
)

func init() {
	const key = "InstallDefaultInstructions"

	database.MigrationsVersionUpgradeCompletedSignal.AddListener(key, func(ctx context.Context, event database.MigratedEvent) error {
		// Always execute, but only if our latest version is above 3 (instructions).
		if event.ToVersion <= 3 {
			return nil
		}

		logger := log.Get()

		logger.Info("Checking default instructions...")
		existing, err := AllInstructions()
		if err != nil {
			return errors.Wrap(err, "failed to get existing instructions")
		}

		existingNames := make(map[string]struct{}, len(existing))
		for _, inst := range existing {
			existingNames[inst.Name] = struct{}{}
		}

		for tplKey, tplName := range DefaultTemplates() {
			if _, exists := existingNames[tplName]; !exists {
				template, err := ReifyInstructionTemplate(tplKey)
				if err != nil {
					return errors.Wrapf(err, "failed to create default template for key %s", tplKey)
				}
				err = CreateInstruction(template)
				if err != nil {
					return errors.Wrapf(err, "failed to save instruction template for key %s", tplKey)
				}
			}
		}

		return nil
	})
}
