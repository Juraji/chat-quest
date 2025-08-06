package worlds

import (
	"database/sql"
	"juraji.nl/chat-quest/database"
)

type World struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ChatPreferences struct {
	ChatModelID       *int64 `json:"chatModelId"`
	ChatInstructionID *int64 `json:"chatInstructionId"`
}

func worldScanner(scanner database.RowScanner, dest *World) error {
	return scanner.Scan(
		&dest.ID,
		&dest.Name,
		&dest.Description,
	)
}
func chatPreferencesScanner(scanner database.RowScanner, dest *ChatPreferences) error {
	return scanner.Scan(
		&dest.ChatModelID,
		&dest.ChatInstructionID,
	)
}

func GetAllWorlds(db *sql.DB) ([]*World, error) {
	query := "SELECT * FROM worlds"
	return database.QueryForList(db, query, nil, worldScanner)
}

func WorldById(db *sql.DB, id int64) (*World, error) {
	query := "SELECT * FROM worlds WHERE id=$1"
	args := []any{id}

	return database.QueryForRecord(db, query, args, worldScanner)
}

func CreateWorld(db *sql.DB, newWorld *World) error {
	query := "INSERT INTO worlds (name, description) VALUES ($1, $2) RETURNING id"
	args := []any{newWorld.Name, newWorld.Description}

	return database.InsertRecord(db, query, args, &newWorld.ID)
}

func UpdateWorld(db *sql.DB, id int64, world *World) error {
	query := `UPDATE worlds
            SET name=$2,
                description=$3
            WHERE id=$1`
	args := []any{id, world.Name, world.Description}

	return database.UpdateRecord(db, query, args)
}

func DeleteWorld(db *sql.DB, id int64) error {
	query := "DELETE FROM worlds WHERE id=$1"
	args := []any{id}

	return database.DeleteRecord(db, query, args)
}

func GetChatPreferences(db *sql.DB) (*ChatPreferences, error) {
	query := "SELECT chat_model_id, chat_instruction_id FROM chat_preferences WHERE id = 0"
	return database.QueryForRecord(db, query, nil, chatPreferencesScanner)
}

func UpdateChatPreferences(db *sql.DB, prefs *ChatPreferences) error {
	query := `UPDATE chat_preferences
            SET chat_model_id = $1,
                chat_instruction_id = $2
            WHERE id = 0`
	args := []any{
		prefs.ChatModelID,
		prefs.ChatInstructionID,
	}
	return database.UpdateRecord(db, query, args)
}
