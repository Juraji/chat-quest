package species

import (
	"juraji.nl/chat-quest/core/database"
	"juraji.nl/chat-quest/core/util"
)

type Species struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	AvatarUrl   *string `json:"avatarUrl"`
}

func speciesScanner(scanner database.RowScanner, dest *Species) error {
	return scanner.Scan(
		&dest.ID,
		&dest.Name,
		&dest.Description,
		&dest.AvatarUrl,
	)
}

func AllSpecies() ([]Species, error) {
	query := "SELECT * FROM species"
	return database.QueryForList(query, nil, speciesScanner)
}

func SpeciesByID(id int) (*Species, error) {
	query := "SELECT * FROM species WHERE id=?"
	args := []any{id}
	return database.QueryForRecord(query, args, speciesScanner)
}

func CreateSpecies(species *Species) error {
	species.AvatarUrl = util.EmptyStrToNil(species.AvatarUrl)

	query := `INSERT INTO species(name, description, avatar_url) VALUES (?, ?, ?) RETURNING id`
	args := []any{species.Name, species.Description, species.AvatarUrl}

	err := database.InsertRecord(query, args, &species.ID)

	if err == nil {
		SpeciesCreatedSignal.EmitBG(species)
	}
	return err
}

func UpdateSpecies(id int, species *Species) error {
	species.AvatarUrl = util.EmptyStrToNil(species.AvatarUrl)

	query := `UPDATE species
			SET name=?,
			    description=?,
			    avatar_url=?
			WHERE id=?`
	args := []any{species.Name, species.Description, species.AvatarUrl, id}

	err := database.UpdateRecord(query, args)

	if err == nil {
		SpeciesUpdatedSignal.EmitBG(species)
	}

	return err
}

func DeleteSpecies(id int) error {
	query := "DELETE FROM species WHERE id=?"
	args := []any{id}

	_, err := database.DeleteRecord(query, args)

	if err == nil {
		SpeciesDeletedSignal.EmitBG(id)
	}

	return err
}

func GetSpeciesPresentInSession(sessionId int) ([]Species, error) {
	query := `SELECT DISTINCT s.*
			FROM species s
				LEFT JOIN characters c on c.species_id = s.id
				LEFT JOIN chat_participants cp on c.id = cp.character_id
			WHERE cp.chat_session_id=?
				AND c.species_id IS NOT NULL

			UNION

			SELECT DISTINCT s.*
			FROM species s
				LEFT JOIN characters c on c.species_id = s.id
				LEFT JOIN chat_sessions cs on cs.persona_id = c.id
			WHERE cs.id=?
				AND cs.persona_id IS NOT NULL`
	args := []any{sessionId, sessionId}
	return database.QueryForList(query, args, speciesScanner)
}
