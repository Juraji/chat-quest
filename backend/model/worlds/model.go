package worlds

import (
	"juraji.nl/chat-quest/core/database"
	"juraji.nl/chat-quest/core/util"
)

type World struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	AvatarUrl   *string `json:"avatarUrl"`
}

type ChatPreferences struct {
	ChatModelID       *int `json:"chatModelId"`
	ChatInstructionID *int `json:"chatInstructionId"`
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

func GetAllWorlds() ([]World, error) {
	query := "SELECT * FROM worlds"
	return database.QueryForList(database.GetDB(), query, nil, worldScanner)
}

func WorldById(id int) (*World, error) {
	query := "SELECT * FROM worlds WHERE id=?"
	args := []any{id}

	return database.QueryForRecord(database.GetDB(), query, args, worldScanner)
}

func CreateWorld(newWorld *World) error {
	util.EmptyStrPtrToNil(&newWorld.Description)
	util.EmptyStrPtrToNil(&newWorld.AvatarUrl)

	query := "INSERT INTO worlds (name, description, avatar_url) VALUES (?, ?, ?) RETURNING id"
	args := []any{newWorld.Name, newWorld.Description, newWorld.AvatarUrl}

	err := database.InsertRecord(database.GetDB(), query, args, &newWorld.ID)
	if err != nil {
		return err
	}

	util.Emit(WorldCreatedSignal, newWorld)
	return nil
}

func UpdateWorld(id int, world *World) error {
	util.EmptyStrPtrToNil(&world.Description)
	util.EmptyStrPtrToNil(&world.AvatarUrl)

	query := `UPDATE worlds
            SET name=?,
                description=?,
                avatar_url=?
            WHERE id=?`
	args := []any{
		world.Name,
		world.Description,
		world.AvatarUrl,
		id,
	}

	err := database.UpdateRecord(database.GetDB(), query, args)
	if err != nil {
		return err
	}

	util.Emit(WorldUpdatedSignal, world)
	return nil
}

func DeleteWorld(id int) error {
	query := "DELETE FROM worlds WHERE id=?"
	args := []any{id}

	err := database.DeleteRecord(database.GetDB(), query, args)
	if err != nil {
		return err
	}

	util.Emit(WorldDeletedSignal, id)
	return nil
}

func GetChatPreferences() (*ChatPreferences, error) {
	query := "SELECT chat_model_id, chat_instruction_id FROM chat_preferences WHERE id = 0"
	return database.QueryForRecord(database.GetDB(), query, nil, chatPreferencesScanner)
}

func UpdateChatPreferences(prefs *ChatPreferences) error {
	query := `UPDATE chat_preferences
            SET chat_model_id = ?,
                chat_instruction_id = ?
            WHERE id = 0`
	args := []any{
		prefs.ChatModelID,
		prefs.ChatInstructionID,
	}

	err := database.UpdateRecord(database.GetDB(), query, args)
	if err != nil {
		return err
	}

	util.Emit(ChatPreferencesUpdatedSignal, prefs)
	return nil
}
