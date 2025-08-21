package worlds

import (
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

func worldScanner(scanner database.RowScanner, dest *World) error {
	return scanner.Scan(
		&dest.ID,
		&dest.Name,
		&dest.Description,
		&dest.AvatarUrl,
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
