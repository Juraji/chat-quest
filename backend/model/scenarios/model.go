package scenarios

import (
	"go.uber.org/zap"
	"juraji.nl/chat-quest/core/database"
	"juraji.nl/chat-quest/core/log"
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

func AllScenarios() ([]Scenario, bool) {
	query := "SELECT * FROM scenarios"
	list, err := database.QueryForList(query, nil, scenarioScanner)
	if err != nil {
		log.Get().Error("Error fetching scenarios", zap.Error(err))
		return nil, false
	}

	return list, true
}

func ScenarioById(id int) (*Scenario, bool) {
	query := "SELECT * FROM scenarios WHERE id=?"
	args := []any{id}
	scene, err := database.QueryForRecord(query, args, scenarioScanner)
	if err != nil {
		log.Get().Error("Error fetching scenario",
			zap.Int("id", id), zap.Error(err))
		return nil, false
	}

	return scene, true
}

func CreateScenario(scenario *Scenario) bool {
	util.EmptyStrPtrToNil(&scenario.AvatarUrl)

	query := `INSERT INTO scenarios (name, description, avatar_url, linked_character_id)
            VALUES (?, ?, ?, ?) RETURNING id`
	args := []interface{}{scenario.Name, scenario.Description, scenario.AvatarUrl, scenario.LinkedCharacterId}

	err := database.InsertRecord(query, args, &scenario.ID)
	if err != nil {
		log.Get().Error("Error inserting scenario", zap.Error(err))
		return false
	}

	ScenarioCreatedSignal.EmitBG(scenario)
	return true
}

func UpdateScenario(id int, scenario *Scenario) bool {
	util.EmptyStrPtrToNil(&scenario.AvatarUrl)

	query := `UPDATE scenarios
            SET name=?,
                description=?,
                avatar_url=?,
                linked_character_id=?
            WHERE id=?`
	args := []interface{}{scenario.Name, scenario.Description, scenario.AvatarUrl, scenario.LinkedCharacterId, id}

	err := database.UpdateRecord(query, args)
	if err != nil {
		log.Get().Error("Error updating scenario",
			zap.Int("id", id), zap.Error(err))
		return false
	}

	ScenarioUpdatedSignal.EmitBG(scenario)
	return true
}

func DeleteScenario(id int) bool {
	query := "DELETE FROM scenarios WHERE id=?"
	args := []interface{}{id}

	err := database.DeleteRecord(query, args)
	if err != nil {
		log.Get().Error("Error deleting scenario",
			zap.Int("id", id), zap.Error(err))
		return false
	}

	ScenarioDeletedSignal.EmitBG(id)
	return true
}
