package model

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func InitDB() (*sql.DB, error) {
	log.Println("Connecting to database...")
	db, err := sql.Open("sqlite3", "./chat-quest.db")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	log.Println("Enabling foreign keys...")
	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %v", err)
	}

	log.Println("Running migrations...")
	if err = RunMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %v", err)
	}

	log.Println("Database initialized successfully!")
	return db, nil
}

func RunMigrations(db *sql.DB) error {
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("failed to create driver instance: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://./migrations",
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrations: %v", err)
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to migrate database: %v", err)
	}

	return nil
}

func QueryForList[T any](
	db *sql.DB,
	query string,
	scanFunc func(rows *sql.Rows, dest *T) error,
) ([]T, error) {
	records := make([]T, 0)

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	for rows.Next() {
		var dest T

		err := scanFunc(rows, &dest)
		if err != nil {
			return nil, err
		}

		records = append(records, dest)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return records, nil
}
