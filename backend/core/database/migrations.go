package database

import (
	"context"
	"database/sql"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"go.uber.org/zap"
	"juraji.nl/chat-quest/core/log"
	"juraji.nl/chat-quest/core/util"
)

import "embed"

//go:embed migrations/*.sql
var migrationsFs embed.FS

func init() {
	MigrationsCompletedSignal.AddListener(func(ctx context.Context, event MigratedEvent) {
		logger := log.Get()
		switch {
		case event.FromVersion < event.ToVersion:
			logger.Sugar().Infof("Database migrated up from v%d to v%d",
				event.FromVersion, event.ToVersion)
		case event.ToVersion < event.FromVersion:
			logger.Sugar().Infof("Database migrated down from v%d to v%d",
				event.FromVersion, event.ToVersion)
		default:
			logger.Info("No database migrations necessary")
		}
	})
}

func runLatestMigrations(db *sql.DB) {
	runUsingMigrations(db, func(m *migrate.Migrate) error {
		return m.Up()
	})
}

func GoToVersion(version uint) {
	runUsingMigrations(GetDB(), func(m *migrate.Migrate) error {
		return m.Migrate(version)
	})
}

func runUsingMigrations(db *sql.DB, action func(m *migrate.Migrate) error) {
	logger := log.Get()

	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		logger.Fatal("Failed to open database for migrations", zap.Error(err))
	}

	fsDriver, err := iofs.New(migrationsFs, "migrations")
	if err != nil {
		logger.Fatal("Failed to create fs driver", zap.Error(err))
	}

	m, err := migrate.NewWithInstance("iofs", fsDriver, "main", driver)
	if err != nil {
		logger.Fatal("Failed to create migrations object", zap.Error(err))
	}

	m.PrefetchMigrations = 0

	fromVersion, _, err := m.Version()
	if err != nil {
		logger.Fatal("Failed to get old version from migrations", zap.Error(err))
	}

	err = action(m)
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logger.Fatal("Failed to apply migrations", zap.Error(err))
	}

	toVersion, _, err := m.Version()
	if err != nil {
		logger.Fatal("Failed to get new version from migrations", zap.Error(err))
	}

	util.Emit(MigrationsCompletedSignal, MigratedEvent{FromVersion: fromVersion, ToVersion: toVersion})
}
