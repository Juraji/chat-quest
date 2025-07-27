package model

import "database/sql"

type Scenario struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Scene string `json:"scene"`
}

func scenarioScanner(scanner RowScanner, dest *Scenario) error {
	return scanner.Scan(
		&dest.ID,
		&dest.Name,
		&dest.Scene,
	)
}

func AllScenarios(db *sql.DB) ([]*Scenario, error) {
	query := "SELECT * FROM scenarios"
	return queryForList(db, query, nil, scenarioScanner)
}

func ScenarioById(db *sql.DB, id int64) (*Scenario, error) {
	query := "SELECT * FROM scenarios WHERE id = $1"
	args := []any{id}
	return queryForRecord(db, query, args, scenarioScanner)
}

func CreateScenario(db *sql.DB, newScenario *Scenario) error {
	query := "INSERT INTO scenarios (name, scene) VALUES($1, $2)"
	args := []any{newScenario.Name, newScenario.Scene}
	scanFunc := func(scanner RowScanner) error {
		return scanner.Scan(&newScenario.ID)
	}

	return insertRecord(db, query, args, scanFunc)
}

func UpdateScenario(db *sql.DB, id int64, scenario *Scenario) error {
	query := "UPDATE scenarios SET name = $1, scene = $2 WHERE id = $3"
	args := []any{scenario.Name, scenario.Scene, id}

	return updateRecord(db, query, args)
}

func DeleteScenario(db *sql.DB, id int64) error {
	query := "DELETE FROM scenarios WHERE id = $1"
	args := []any{id}

	return deleteRecord(db, query, args)
}
