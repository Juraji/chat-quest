package worlds

import (
	"database/sql"
	"juraji.nl/chat-quest/database"
	"juraji.nl/chat-quest/util"
)

type World struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	AvatarUrl   *string `json:"avatarUrl"`
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
		&dest.AvatarUrl,
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
	util.EmptyStrPtrToNil(&newWorld.AvatarUrl)

	query := "INSERT INTO worlds (name, description, avatar_url) VALUES ($1, $2, $3) RETURNING id"
	args := []any{newWorld.Name, newWorld.Description, newWorld.AvatarUrl}

	err := database.InsertRecord(db, query, args, &newWorld.ID)
	defer util.EmitOnSuccess(WorldCreatedSignal, newWorld, err)

	return err
}

func UpdateWorld(db *sql.DB, id int64, world *World) error {
	util.EmptyStrPtrToNil(&world.AvatarUrl)

	query := `UPDATE worlds
            SET name=$2,
                description=$3,
                avatar_url=$4
            WHERE id=$1`
	args := []any{id, world.Name, world.Description, world.AvatarUrl}

	err := database.UpdateRecord(db, query, args)
	defer util.EmitOnSuccess(WorldUpdatedSignal, world, err)

	return err
}

func DeleteWorld(db *sql.DB, id int64) error {
	query := "DELETE FROM worlds WHERE id=$1"
	args := []any{id}

	err := database.DeleteRecord(db, query, args)
	defer util.EmitOnSuccess(WorldDeletedSignal, id, err)

	return err
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

	err := database.UpdateRecord(db, query, args)
	defer util.EmitOnSuccess(ChatPreferencesUpdatedSignal, prefs, err)

	return err
}
