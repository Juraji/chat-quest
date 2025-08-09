package scenarios

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
	query := "SELECT * FROM scenarios WHERE id=?"
	args := []any{id}
	return database.QueryForRecord(db, query, args, scenarioScanner)
}

func CreateScenario(db *sql.DB, scenario *Scenario) error {
	util.EmptyStrPtrToNil(&scenario.AvatarUrl)

	query := `INSERT INTO scenarios (name, description, avatar_url, linked_character_id)
            VALUES (?, ?, ?, ?) RETURNING id`
	args := []interface{}{scenario.Name, scenario.Description, scenario.AvatarUrl, scenario.LinkedCharacterId}

	err := database.InsertRecord(db, query, args, &scenario.ID)
	defer util.EmitOnSuccess(ScenarioCreatedSignal, scenario, err)

	return err
}

func UpdateScenario(db *sql.DB, id int64, scenario *Scenario) error {
	util.EmptyStrPtrToNil(&scenario.AvatarUrl)

	query := `UPDATE scenarios
            SET name=?,
                description=?,
                avatar_url=?,
                linked_character_id=?
            WHERE id=?`
	args := []interface{}{scenario.Name, scenario.Description, scenario.AvatarUrl, scenario.LinkedCharacterId, id}

	err := database.UpdateRecord(db, query, args)
	defer util.EmitOnSuccess(ScenarioUpdatedSignal, scenario, err)

	return err
}

func DeleteScenario(db *sql.DB, id int64) error {
	query := "DELETE FROM scenarios WHERE id=?"
	args := []interface{}{id}

	err := database.DeleteRecord(db, query, args)
	defer util.EmitOnSuccess(ScenarioDeletedSignal, id, err)

	return err
}
