package characters

import (
	"go.uber.org/zap"
	"juraji.nl/chat-quest/core/database"
	"juraji.nl/chat-quest/core/log"
	"juraji.nl/chat-quest/core/util"
	"time"
)

type Character struct {
	ID                 int        `json:"id"`
	CreatedAt          *time.Time `json:"createdAt"`
	Name               string     `json:"name"`
	Favorite           bool       `json:"favorite"`
	AvatarUrl          *string    `json:"avatarUrl"`
	Appearance         *string    `json:"appearance"`
	Personality        *string    `json:"personality"`
	History            *string    `json:"history"`
	GroupTalkativeness float32    `json:"groupTalkativeness"`
}

type CharacterListView struct {
	ID        int        `json:"id"`
	CreatedAt *time.Time `json:"createdAt"`
	Name      string     `json:"name"`
	Favorite  bool       `json:"favorite"`
	AvatarUrl *string    `json:"avatarUrl"`
	Tags      []Tag      `json:"tags"`
}

func CharacterScanner(scanner database.RowScanner, dest *Character) error {
	return scanner.Scan(
		&dest.ID,
		&dest.CreatedAt,
		&dest.Name,
		&dest.Favorite,
		&dest.AvatarUrl,
		&dest.Appearance,
		&dest.Personality,
		&dest.History,
		&dest.GroupTalkativeness,
	)
}

func characterListViewScanner(scanner database.RowScanner, dest *CharacterListView) error {
	return scanner.Scan(
		&dest.ID,
		&dest.CreatedAt,
		&dest.Name,
		&dest.Favorite,
		&dest.AvatarUrl,
	)
}

func AllCharacterListViews() ([]CharacterListView, error) {
	query := "SELECT id, created_at, name, favorite, avatar_url FROM characters"

	characters, err := database.QueryForList(query, nil, characterListViewScanner)
	if err != nil {
		return nil, err
	}

	for i := range characters {
		char := &characters[i]
		tags, _ := TagsByCharacterId(char.ID)
		char.Tags = tags
	}

	return characters, nil
}

func CharacterById(id int) (*Character, error) {
	query := "SELECT * FROM characters WHERE id = ?"
	args := []any{id}

	return database.QueryForRecord(query, args, CharacterScanner)
}

func CreateCharacter(newCharacter *Character) error {
	newCharacter.AvatarUrl = util.EmptyStrToNil(newCharacter.AvatarUrl)
	newCharacter.Appearance = util.EmptyStrToNil(newCharacter.Appearance)
	newCharacter.Personality = util.EmptyStrToNil(newCharacter.Personality)
	newCharacter.History = util.EmptyStrToNil(newCharacter.History)

	query := `INSERT INTO characters (name, favorite, avatar_url, appearance, personality, history, group_talkativeness)
            VALUES (?, ?, ?, ?, ?, ?, ?) RETURNING id, created_at`
	args := []any{
		newCharacter.Name,
		newCharacter.Favorite,
		newCharacter.AvatarUrl,
		newCharacter.Appearance,
		newCharacter.Personality,
		newCharacter.History,
		newCharacter.GroupTalkativeness,
	}

	err := database.InsertRecord(query, args, &newCharacter.ID, &newCharacter.CreatedAt)

	if err == nil {
		CharacterCreatedSignal.EmitBG(newCharacter)
	}

	return err
}

func UpdateCharacter(id int, character *Character) error {
	character.AvatarUrl = util.EmptyStrToNil(character.AvatarUrl)
	character.Appearance = util.EmptyStrToNil(character.Appearance)
	character.Personality = util.EmptyStrToNil(character.Personality)
	character.History = util.EmptyStrToNil(character.History)

	query := `UPDATE characters
            SET name = ?,
                favorite = ?,
                avatar_url = ?,
                appearance = ?,
                personality = ?,
                history = ?,
                group_talkativeness = ?
            WHERE id = ?`
	args := []any{
		character.Name,
		character.Favorite,
		character.AvatarUrl,
		character.Appearance,
		character.Personality,
		character.History,
		character.GroupTalkativeness,
		id,
	}

	err := database.UpdateRecord(query, args)

	if err == nil {
		CharacterUpdatedSignal.EmitBG(character)
	}

	return err
}

func DeleteCharacterById(id int) error {
	query := "DELETE FROM characters WHERE id = ?"
	args := []any{id}

	_, err := database.DeleteRecord(query, args)

	if err != nil {
		CharacterDeletedSignal.EmitBG(id)
	}

	return err
}

func DialogueExamplesByCharacterId(characterId int) ([]string, error) {
	query := "SELECT text FROM character_dialogue_examples WHERE character_id = ?"
	args := []any{characterId}
	scanFunc := func(rows database.RowScanner, dest *string) error {
		return rows.Scan(dest)
	}

	return database.QueryForList(query, args, scanFunc)
}

func SetDialogueExamplesByCharacterId(characterId int, examples []string) error {
	return database.Transactional(func(ctx *database.TxContext) error {
		deleteQuery := "DELETE FROM character_dialogue_examples WHERE character_id = ?"
		if _, err := ctx.DeleteRecord(deleteQuery, []any{characterId}); err != nil {
			log.Get().Error("Error removing dialogue examples",
				zap.Int("characterId", characterId), zap.Error(err))
			return err
		}

		if len(examples) == 0 {
			return nil
		}

		insertQuery := "INSERT INTO character_dialogue_examples (character_id, text) VALUES (?, ?)"
		for _, example := range examples {
			if err := ctx.InsertRecord(insertQuery, []any{characterId, example}); err != nil {
				log.Get().Error("Error adding dialogue examples",
					zap.Int("characterId", characterId), zap.Error(err))
				return err
			}
		}

		return nil
	})
}

func CharacterGreetingsByCharacterId(characterId int) ([]string, error) {
	query := "SELECT text FROM character_greetings WHERE character_id = ?"
	args := []any{characterId}

	return database.QueryForList(query, args, database.StringScanner)
}

func SetGreetingsByCharacterId(characterId int, greetings []string) error {
	return database.Transactional(func(ctx *database.TxContext) error {
		deleteQuery := "DELETE FROM character_greetings WHERE character_id = ?"
		if _, err := ctx.DeleteRecord(deleteQuery, []any{characterId}); err != nil {
			log.Get().Error("Error removing greetings",
				zap.Int("characterId", characterId), zap.Error(err))
			return err
		}

		if len(greetings) == 0 {
			return nil
		}

		insertQuery := "INSERT INTO character_greetings (character_id, text) VALUES (?, ?)"
		for _, greeting := range greetings {
			if err := ctx.InsertRecord(insertQuery, []any{characterId, greeting}); err != nil {
				log.Get().Error("Error adding greetings",
					zap.Int("characterId", characterId), zap.Error(err))
				return err
			}
		}

		return nil
	})
}

func CharacterGroupGreetingsByCharacterId(characterId int) ([]string, error) {
	query := "SELECT text FROM character_group_greetings WHERE character_id = ?"
	args := []any{characterId}

	return database.QueryForList(query, args, database.StringScanner)
}

func SetGroupGreetingsByCharacterId(characterId int, greetings []string) error {
	return database.Transactional(func(ctx *database.TxContext) error {
		deleteQuery := "DELETE FROM character_group_greetings WHERE character_id = ?"
		if _, err := ctx.DeleteRecord(deleteQuery, []any{characterId}); err != nil {
			log.Get().Error("Error removing group greetings",
				zap.Int("characterId", characterId), zap.Error(err))
			return err
		}

		if len(greetings) == 0 {
			return nil
		}

		insertQuery := "INSERT INTO character_group_greetings (character_id, text) VALUES (?, ?)"
		for _, greeting := range greetings {
			if err := ctx.InsertRecord(insertQuery, []any{characterId, greeting}); err != nil {
				log.Get().Error("Error adding group greetings",
					zap.Int("characterId", characterId), zap.Error(err))
				return err
			}
		}

		return nil
	})
}

func RandomGreetingByCharacterId(characterId int, useGroupGreetings bool) (*string, error) {
	query := `SELECT text FROM character_greetings WHERE character_id = ? ORDER BY RANDOM() LIMIT 1;`
	if useGroupGreetings {
		query = `SELECT text FROM character_group_greetings WHERE character_id = ? ORDER BY RANDOM() LIMIT 1;`
	}
	args := []any{characterId, useGroupGreetings}

	return database.QueryForRecord(query, args, database.StringScanner)
}
