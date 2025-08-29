package scenarios

import (
	"juraji.nl/chat-quest/core/database"
	"juraji.nl/chat-quest/core/util"
)

type Scenario struct {
	ID                int     `json:"id"`
	Name              string  `json:"name"`
	Description       string  `json:"description"`
	AvatarUrl         *string `json:"avatarUrl"`
	LinkedCharacterId *int    `json:"linkedCharacterId"`
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

func AllScenarios() ([]Scenario, error) {
	query := "SELECT * FROM scenarios"
	return database.QueryForList(query, nil, scenarioScanner)
}

func ScenarioById(id int) (*Scenario, error) {
	query := "SELECT * FROM scenarios WHERE id=?"
	args := []any{id}
	return database.QueryForRecord(query, args, scenarioScanner)
}

func CreateScenario(scenario *Scenario) error {
	util.EmptyStrPtrToNil(&scenario.AvatarUrl)

	query := `INSERT INTO scenarios (name, description, avatar_url, linked_character_id)
            VALUES (?, ?, ?, ?) RETURNING id`
	args := []interface{}{scenario.Name, scenario.Description, scenario.AvatarUrl, scenario.LinkedCharacterId}

	err := database.InsertRecord(query, args, &scenario.ID)

	if err == nil {
		ScenarioCreatedSignal.EmitBG(scenario)
	}
	return err
}

func UpdateScenario(id int, scenario *Scenario) error {
	util.EmptyStrPtrToNil(&scenario.AvatarUrl)

	query := `UPDATE scenarios
            SET name=?,
                description=?,
                avatar_url=?,
                linked_character_id=?
            WHERE id=?`
	args := []interface{}{scenario.Name, scenario.Description, scenario.AvatarUrl, scenario.LinkedCharacterId, id}

	err := database.UpdateRecord(query, args)

	if err == nil {
		ScenarioUpdatedSignal.EmitBG(scenario)
	}

	return err
}

func DeleteScenario(id int) error {
	query := "DELETE FROM scenarios WHERE id=?"
	args := []interface{}{id}

	_, err := database.DeleteRecord(query, args)

	if err == nil {
		ScenarioDeletedSignal.EmitBG(id)
	}

	return err
}
