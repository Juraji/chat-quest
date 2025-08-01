package migrations

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed sql/*.sql
var migrationsFs embed.FS

func RunMigrations(db *sql.DB) error {
	return runUsingMigrations(db, func(m *migrate.Migrate) error {
		return m.Up()
	})
}

func GoToVersion(db *sql.DB, version uint) error {
	return runUsingMigrations(db, func(m *migrate.Migrate) error {
		return m.Migrate(version)
	})
}

func runUsingMigrations(db *sql.DB, action func(m *migrate.Migrate) error) error {
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("failed to create databse driver instance: %v", err)
	}

	migrationsFsDriver, err := iofs.New(migrationsFs, "sql")
	if err != nil {
		return fmt.Errorf("failed to create migrations fs driver: %v", err)
	}

	m, err := migrate.NewWithInstance("iofs", migrationsFsDriver, "main", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrations: %v", err)
	}

	if err = action(m); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to migrate database: %v", err)
	}

	return nil
}
