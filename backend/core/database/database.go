package database

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./chat-quest.db")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	if err = RunMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return db, nil
}
