package scenarios

import (
	"juraji.nl/chat-quest/core"
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

func AllScenarios(cq *core.ChatQuestContext) ([]*Scenario, error) {
	query := "SELECT * FROM scenarios"
	return database.QueryForList(cq.DB(), query, nil, scenarioScanner)
}

func ScenarioById(cq *core.ChatQuestContext, id int) (*Scenario, error) {
	query := "SELECT * FROM scenarios WHERE id=?"
	args := []any{id}
	return database.QueryForRecord(cq.DB(), query, args, scenarioScanner)
}

func CreateScenario(cq *core.ChatQuestContext, scenario *Scenario) error {
	util.EmptyStrPtrToNil(&scenario.AvatarUrl)

	query := `INSERT INTO scenarios (name, description, avatar_url, linked_character_id)
            VALUES (?, ?, ?, ?) RETURNING id`
	args := []interface{}{scenario.Name, scenario.Description, scenario.AvatarUrl, scenario.LinkedCharacterId}

	err := database.InsertRecord(cq.DB(), query, args, &scenario.ID)
	if err != nil {
		return err
	}

	ScenarioCreatedSignal.Emit(cq.Context(), scenario)
	return nil
}

func UpdateScenario(cq *core.ChatQuestContext, id int, scenario *Scenario) error {
	util.EmptyStrPtrToNil(&scenario.AvatarUrl)

	query := `UPDATE scenarios
            SET name=?,
                description=?,
                avatar_url=?,
                linked_character_id=?
            WHERE id=?`
	args := []interface{}{scenario.Name, scenario.Description, scenario.AvatarUrl, scenario.LinkedCharacterId, id}

	err := database.UpdateRecord(cq.DB(), query, args)
	if err != nil {
		return err
	}

	ScenarioUpdatedSignal.Emit(cq.Context(), scenario)
	return nil
}

func DeleteScenario(cq *core.ChatQuestContext, id int) error {
	query := "DELETE FROM scenarios WHERE id=?"
	args := []interface{}{id}

	err := database.DeleteRecord(cq.DB(), query, args)
	if err != nil {
		return err
	}

	ScenarioDeletedSignal.Emit(cq.Context(), id)
	return nil
}
