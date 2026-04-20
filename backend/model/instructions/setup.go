package instructions

import (
	"context"

	"go.uber.org/zap"
	"juraji.nl/chat-quest/core/database"
	"juraji.nl/chat-quest/core/log"
)

func init() {
	const key = "InstallDefaultInstructions"

	database.MigrationsVersionUpgradeCompletedSignal.AddListener(key, func(ctx context.Context, event database.MigratedEvent) {
		// Always execute, but only if our latest version is above 3 (instructions).
		if event.ToVersion <= 3 {
			return
		}

		logger := log.Get()

		logger.Info("Checking default instructions...")
		existing, err := AllInstructions()
		if err != nil {
			logger.Panic("failed to get existing instructions",
				zap.Error(err))
		}

		existingNames := make(map[string]struct{}, len(existing))
		for _, inst := range existing {
			existingNames[inst.Name] = struct{}{}
		}

		for tplKey, tplName := range DefaultTemplates() {
			if _, exists := existingNames[tplName]; !exists {
				template, err := ReifyInstructionTemplate(tplKey)
				if err != nil {
					logger.Panic("failed to create default template",
						zap.String("key", tplKey),
						zap.Error(err))
				}
				err = CreateInstruction(template)
				if err != nil {
					logger.Panic("failed to save instruction template",
						zap.String("key", tplKey),
						zap.Error(err))
				}
			}
		}
	})
}
