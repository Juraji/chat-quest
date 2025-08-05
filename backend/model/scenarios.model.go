package model

import (
	"database/sql"
	"juraji.nl/chat-quest/database"
	"juraji.nl/chat-quest/util"
)

type Scenario struct {
	ID                int64   `json:"id"`
	Name              string  `json:"name"`
	Description       string  `json:"description"`
	AvatarUrl         *string `json:"avatarUrl"`
	LinkedCharacterId *int64  `json:"linkedCharacterId"`
}

func scenarioScanner(scanner database.RowScanner, dest *Scenario) error {
	return scanner.Scan(
		&dest.ID,
		&dest.Name,
		&dest.Description,
		&dest.AvatarUrl,
		&dest.LinkedCharacterId,
	)
}

func AllScenarios(db *sql.DB) ([]*Scenario, error) {
	query := "SELECT * FROM scenarios"
	return database.QueryForList(db, query, nil, scenarioScanner)
}

func ScenarioById(db *sql.DB, id int64) (*Scenario, error) {
	query := "SELECT * FROM scenarios WHERE id=$1"
	args := []any{id}
	return database.QueryForRecord(db, query, args, scenarioScanner)
}

func CreateScenario(db *sql.DB, scenario *Scenario) error {
	util.EmptyStrPtrToNil(&scenario.AvatarUrl)

	query := `INSERT INTO scenarios (name, description, avatar_url, linked_character_id)
            VALUES ($1, $2, $3, $4) RETURNING id`
	args := []interface{}{scenario.Name, scenario.Description, scenario.AvatarUrl, scenario.LinkedCharacterId}
	scanFunc := func(scanner database.RowScanner) error {
		return scanner.Scan(&scenario.ID)
	}

	return database.InsertRecord(db, query, args, scanFunc)
}

func UpdateScenario(db *sql.DB, id int64, scenario *Scenario) error {
	util.EmptyStrPtrToNil(&scenario.AvatarUrl)

	query := `UPDATE scenarios
            SET name=$1,
                description=$2,
                avatar_url=$3,
                linked_character_id=$4
            WHERE id=$5`
	args := []interface{}{scenario.Name, scenario.Description, scenario.AvatarUrl, scenario.LinkedCharacterId, id}

	return database.UpdateRecord(db, query, args)
}

func DeleteScenario(db *sql.DB, id int64) error {
	query := "DELETE FROM scenarios WHERE id=$1"
	args := []interface{}{id}
	return database.DeleteRecord(db, query, args)
}
