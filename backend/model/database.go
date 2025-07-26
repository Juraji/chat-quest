package model

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"juraji.nl/chat-quest/migrations"
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
	if err = migrations.RunMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %v", err)
	}

	log.Println("Database initialized successfully!")
	return db, nil
}
