package worlds

import (
	"errors"
	"go.uber.org/zap"
	"juraji.nl/chat-quest/core/database"
	"juraji.nl/chat-quest/core/log"
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

func (p *ChatPreferences) Validate() error {
	if p == nil {
		return errors.New("chat preferences is nil")
	}
	if p.ChatModelID == nil {
		return errors.New("ChatModelId is nil")
	}
	if p.ChatInstructionID == nil {
		return errors.New("ChatInstructionId is nil")
	}

	return nil
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

func GetAllWorlds() ([]World, bool) {
	query := "SELECT * FROM worlds"
	list, err := database.QueryForList(query, nil, worldScanner)
	if err != nil {
		log.Get().Error("Error fetching worlds", zap.Error(err))
		return nil, false
	}

	return list, true
}

func WorldById(id int) (*World, bool) {
	query := "SELECT * FROM worlds WHERE id=?"
	args := []any{id}

	world, err := database.QueryForRecord(query, args, worldScanner)
	if err != nil {
		log.Get().Error("Error fetching world",
			zap.Int("id", id), zap.Error(err))
		return nil, false
	}

	return world, true
}

func CreateWorld(newWorld *World) bool {
	util.EmptyStrPtrToNil(&newWorld.Description)
	util.EmptyStrPtrToNil(&newWorld.AvatarUrl)

	query := "INSERT INTO worlds (name, description, avatar_url) VALUES (?, ?, ?) RETURNING id"
	args := []any{newWorld.Name, newWorld.Description, newWorld.AvatarUrl}

	err := database.InsertRecord(query, args, &newWorld.ID)
	if err != nil {
		log.Get().Error("Error inserting world",
			zap.Int("id", newWorld.ID), zap.Error(err))
		return false
	}

	WorldCreatedSignal.EmitBG(newWorld)
	return true
}

func UpdateWorld(id int, world *World) bool {
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

	err := database.UpdateRecord(query, args)
	if err != nil {
		log.Get().Error("Error updating world",
			zap.Int("id", id), zap.Error(err))
		return false
	}

	WorldUpdatedSignal.EmitBG(world)
	return true
}

func DeleteWorld(id int) bool {
	query := "DELETE FROM worlds WHERE id=?"
	args := []any{id}

	err := database.DeleteRecord(query, args)
	if err != nil {
		log.Get().Error("Error deleting world",
			zap.Int("id", id), zap.Error(err))
		return false
	}

	WorldDeletedSignal.EmitBG(id)
	return true
}

func GetChatPreferences() (*ChatPreferences, bool) {
	query := "SELECT chat_model_id, chat_instruction_id FROM chat_preferences WHERE id = 0"
	prefs, err := database.QueryForRecord(query, nil, chatPreferencesScanner)
	if err != nil {
		log.Get().Error("Error fetching chat_preferences", zap.Error(err))
		return nil, false
	}

	return prefs, true
}

func UpdateChatPreferences(prefs *ChatPreferences) bool {
	query := `UPDATE chat_preferences
            SET chat_model_id = ?,
                chat_instruction_id = ?
            WHERE id = 0`
	args := []any{
		prefs.ChatModelID,
		prefs.ChatInstructionID,
	}

	err := database.UpdateRecord(query, args)
	if err != nil {
		log.Get().Error("Error updating chat_preferences", zap.Error(err))
		return false
	}

	ChatPreferencesUpdatedSignal.EmitBG(prefs)
	return true
}
