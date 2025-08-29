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

func worldScanner(scanner database.RowScanner, dest *World) error {
	return scanner.Scan(
		&dest.ID,
		&dest.Name,
		&dest.Description,
		&dest.AvatarUrl,
	)
}

func GetAllWorlds() ([]World, error) {
	query := "SELECT * FROM worlds"
	return database.QueryForList(query, nil, worldScanner)
}

func WorldById(id int) (*World, error) {
	query := "SELECT * FROM worlds WHERE id=?"
	args := []any{id}
	return database.QueryForRecord(query, args, worldScanner)
}

func CreateWorld(newWorld *World) error {
	util.EmptyStrPtrToNil(&newWorld.Description)
	util.EmptyStrPtrToNil(&newWorld.AvatarUrl)

	query := "INSERT INTO worlds (name, description, avatar_url) VALUES (?, ?, ?) RETURNING id"
	args := []any{newWorld.Name, newWorld.Description, newWorld.AvatarUrl}

	err := database.InsertRecord(query, args, &newWorld.ID)

	if err == nil {
		WorldCreatedSignal.EmitBG(newWorld)
	}

	return err
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

	err := database.UpdateRecord(query, args)

	if err == nil {
		WorldUpdatedSignal.EmitBG(world)
	}

	return err
}

func DeleteWorld(id int) error {
	query := "DELETE FROM worlds WHERE id=?"
	args := []any{id}

	_, err := database.DeleteRecord(query, args)

	if err == nil {
		WorldDeletedSignal.EmitBG(id)
	}

	return err
}
