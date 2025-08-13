package worlds

import (
	"juraji.nl/chat-quest/core"
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

func GetAllWorlds(cq *core.ChatQuestContext) ([]*World, error) {
	query := "SELECT * FROM worlds"
	return database.QueryForList(cq.DB(), query, nil, worldScanner)
}

func WorldById(cq *core.ChatQuestContext, id int) (*World, error) {
	query := "SELECT * FROM worlds WHERE id=?"
	args := []any{id}

	return database.QueryForRecord(cq.DB(), query, args, worldScanner)
}

func CreateWorld(cq *core.ChatQuestContext, newWorld *World) error {
	util.EmptyStrPtrToNil(&newWorld.Description)
	util.EmptyStrPtrToNil(&newWorld.AvatarUrl)

	query := "INSERT INTO worlds (name, description, avatar_url) VALUES (?, ?, ?) RETURNING id"
	args := []any{newWorld.Name, newWorld.Description, newWorld.AvatarUrl}

	err := database.InsertRecord(cq.DB(), query, args, &newWorld.ID)
	if err != nil {
		return err
	}

	WorldCreatedSignal.Emit(cq.Context(), newWorld)
	return nil
}

func UpdateWorld(cq *core.ChatQuestContext, id int, world *World) error {
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

	err := database.UpdateRecord(cq.DB(), query, args)
	if err != nil {
		return err
	}

	WorldUpdatedSignal.Emit(cq.Context(), world)
	return nil
}

func DeleteWorld(cq *core.ChatQuestContext, id int) error {
	query := "DELETE FROM worlds WHERE id=?"
	args := []any{id}

	err := database.DeleteRecord(cq.DB(), query, args)
	if err != nil {
		return err
	}

	WorldDeletedSignal.Emit(cq.Context(), id)
	return nil
}

func GetChatPreferences(cq *core.ChatQuestContext) (*ChatPreferences, error) {
	query := "SELECT chat_model_id, chat_instruction_id FROM chat_preferences WHERE id = 0"
	return database.QueryForRecord(cq.DB(), query, nil, chatPreferencesScanner)
}

func UpdateChatPreferences(cq *core.ChatQuestContext, prefs *ChatPreferences) error {
	query := `UPDATE chat_preferences
            SET chat_model_id = ?,
                chat_instruction_id = ?
            WHERE id = 0`
	args := []any{
		prefs.ChatModelID,
		prefs.ChatInstructionID,
	}

	err := database.UpdateRecord(cq.DB(), query, args)
	if err != nil {
		return err
	}

	ChatPreferencesUpdatedSignal.Emit(cq.Context(), prefs)
	return nil
}
