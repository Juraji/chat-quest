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
	return database.QueryForList(database.GetDB(), query, nil, scenarioScanner)
}

func ScenarioById(id int) (*Scenario, error) {
	query := "SELECT * FROM scenarios WHERE id=?"
	args := []any{id}
	return database.QueryForRecord(database.GetDB(), query, args, scenarioScanner)
}

func CreateScenario(scenario *Scenario) error {
	util.EmptyStrPtrToNil(&scenario.AvatarUrl)

	query := `INSERT INTO scenarios (name, description, avatar_url, linked_character_id)
            VALUES (?, ?, ?, ?) RETURNING id`
	args := []interface{}{scenario.Name, scenario.Description, scenario.AvatarUrl, scenario.LinkedCharacterId}

	err := database.InsertRecord(database.GetDB(), query, args, &scenario.ID)
	if err != nil {
		return err
	}

	util.Emit(ScenarioCreatedSignal, scenario)
	return nil
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

	err := database.UpdateRecord(database.GetDB(), query, args)
	if err != nil {
		return err
	}

	util.Emit(ScenarioUpdatedSignal, scenario)
	return nil
}

func DeleteScenario(id int) error {
	query := "DELETE FROM scenarios WHERE id=?"
	args := []interface{}{id}

	err := database.DeleteRecord(database.GetDB(), query, args)
	if err != nil {
		return err
	}

	util.Emit(ScenarioDeletedSignal, id)
	return nil
}
