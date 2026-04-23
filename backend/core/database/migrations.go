package database

import (
	"database/sql"
	"errors"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"go.uber.org/zap"
	"juraji.nl/chat-quest/core/log"
)

import "embed"

//go:embed migrations/*.sql
var migrationsFs embed.FS

func runLatestMigrations(db *sql.DB) {
	runUsingMigrations(db, func(m *migrate.Migrate) error {
		return m.Up()
	})
}

func GoToVersion(db *sql.DB, version uint) {
	runUsingMigrations(db, func(m *migrate.Migrate) error {
		return m.Migrate(version)
	})
}

func runUsingMigrations(db *sql.DB, action func(m *migrate.Migrate) error) {
	logger := log.Get()
	var err error

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
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
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

	// Emit migration event and wait for all listeners to complete
	migratedEvent := MigratedEvent{FromVersion: fromVersion, ToVersion: toVersion}
	err = MigrationsVersionUpgradeCompletedSignal.EmitBG(migratedEvent).Wait()
	if err != nil {
		logger.Fatal("Failed running version upgrade handlers")
	}

	// Run "post_migration.sql" (Only when migrating up)
	if fromVersion >= toVersion {
		postMigrationSqlRaw, err := migrationsFs.ReadFile("migrations/post_migration.sql")
		if err != nil {
			logger.Fatal("Failed to read post migrations file", zap.Error(err))
		}
		postMigrationSql := string(postMigrationSqlRaw)
		_, err = db.Exec(postMigrationSql)
		if err != nil {
			logger.Fatal("Failed to execute post migration", zap.Error(err))
		}
	}

	err = MigrationsPostMigrationCompletedSignal.EmitBG(migratedEvent).Wait()
	if err != nil {
		logger.Fatal("Failed running post upgrade handlers")
	}
}
