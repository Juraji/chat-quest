package database

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

var dbInstance *sql.DB

func GetDB() *sql.DB {
	if dbInstance == nil {
		panic("database not initialized")
	}

	return dbInstance
}

func InitDB() (func(), error) {
	db, err := sql.Open("sqlite3", "./chat-quest.db")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	if err = runLatestMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	dbInstance = db
	closeDB := func() {
		err := db.Close()
		if err != nil {
			panic(fmt.Errorf("failed to close database: %w", err))
		}
	}
	return closeDB, nil
}
