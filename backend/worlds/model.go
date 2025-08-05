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

func worldScanner(scanner database.RowScanner, dest *World) error {
	return scanner.Scan(
		&dest.ID,
		&dest.Name,
		&dest.Description,
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
	scanFunc := func(scanner database.RowScanner) error {
		return scanner.Scan(&newWorld.ID)
	}

	return database.InsertRecord(db, query, args, scanFunc)
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
