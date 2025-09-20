package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
	"juraji.nl/chat-quest/core"
	"juraji.nl/chat-quest/core/log"
)

var dbInstance *sql.DB

func GetDB() *sql.DB {
	if dbInstance == nil {
		panic("database not initialized")
	}

	return dbInstance
}

// InitDB sets up the database, populates dbInstance and runs any migrations.
func InitDB(env core.Environment) func() {
	// DB Setup
	dbLogger := log.Get()

	dbPath := env.MkDataDir("chat-quest.db") + "?_foreign_keys=true"
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		dbLogger.Fatal("Failed to connect to database", zap.Error(err))
	}
	dbInstance = db
	closeDB := func() {
		err := db.Close()
		if err != nil {
			dbLogger.Fatal("Failed to close database", zap.Error(err))
		}
	}

	// Run migrations (will panic if fails)
	runLatestMigrations(db)

	// Return the db closer
	return closeDB
}
