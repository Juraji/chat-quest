package characters

import (
	"go.uber.org/zap"
	"juraji.nl/chat-quest/core/database"
	"juraji.nl/chat-quest/core/log"
	"juraji.nl/chat-quest/core/util"
	"strings"
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
	Tags      []Tag      `json:"tags,omitempty"`
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

func AllCharacterListViews() ([]CharacterListView, bool) {
	query := "SELECT id, created_at, name, favorite, avatar_url FROM characters"

	characters, err := database.QueryForList(query, nil, characterListViewScanner)
	if err != nil {
		log.Get().Error("Error fetching character list views", zap.Error(err))
		return characters, false
	}

	for i := range characters {
		char := characters[i]
		tags, _ := TagsByCharacterId(char.ID)
		char.Tags = tags
	}

	return characters, true
}

func CharacterById(id int) (*Character, bool) {
	query := "SELECT * FROM characters WHERE id = ?"
	args := []any{id}

	list, err := database.QueryForRecord(query, args, CharacterScanner)
	if err != nil {
		log.Get().Error("Error fetching character",
			zap.Int("id", id), zap.Error(err))
		return nil, false
	}

	return list, true
}

func CreateCharacter(newCharacter *Character) bool {
	util.EmptyStrPtrToNil(&newCharacter.AvatarUrl)
	util.EmptyStrPtrToNil(&newCharacter.Appearance)
	util.EmptyStrPtrToNil(&newCharacter.Personality)
	util.EmptyStrPtrToNil(&newCharacter.History)

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

	if err := database.InsertRecord(query, args, &newCharacter.ID, &newCharacter.CreatedAt); err != nil {
		log.Get().Error("Error creating character", zap.Error(err))
		return false
	}

	CharacterCreatedSignal.EmitBG(newCharacter)
	return true
}

func UpdateCharacter(id int, character *Character) bool {
	util.EmptyStrPtrToNil(&character.AvatarUrl)
	util.EmptyStrPtrToNil(&character.Appearance)
	util.EmptyStrPtrToNil(&character.Personality)
	util.EmptyStrPtrToNil(&character.History)

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
	if err != nil {
		log.Get().Error("Error updating character", zap.Int("id", id), zap.Error(err))
		return false
	}

	CharacterUpdatedSignal.EmitBG(character)
	return true
}

func DeleteCharacterById(id int) bool {
	query := "DELETE FROM characters WHERE id = ?"
	args := []any{id}

	err := database.DeleteRecord(query, args)
	if err != nil {
		log.Get().Error("Error deleting character", zap.Int("id", id), zap.Error(err))
		return false
	}

	CharacterDeletedSignal.EmitBG(id)
	return true
}

func DialogueExamplesByCharacterId(characterId int) ([]string, bool) {
	query := "SELECT text FROM character_dialogue_examples WHERE character_id = ?"
	args := []any{characterId}
	scanFunc := func(rows database.RowScanner, dest *string) error {
		return rows.Scan(dest)
	}

	list, err := database.QueryForList(query, args, scanFunc)
	if err != nil {
		log.Get().Error("Error fetching dialogue examples",
			zap.Int("characterId", characterId), zap.Error(err))
		return nil, false
	}

	return list, true
}

func SetDialogueExamplesByCharacterId(characterId int, examples []string) bool {
	err := database.Transactional(func(ctx *database.TxContext) error {
		deleteQuery := "DELETE FROM character_dialogue_examples WHERE character_id = ?"
		if err := ctx.DeleteRecord(deleteQuery, []any{characterId}); err != nil {
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

	return err == nil
}

func CharacterGreetingsByCharacterId(characterId int) ([]string, bool) {
	query := "SELECT text FROM character_greetings WHERE character_id = ?"
	args := []any{characterId}

	list, err := database.QueryForList(query, args, database.StringScanner)
	if err != nil {
		log.Get().Error("Error fetching character greetings",
			zap.Int("characterId", characterId), zap.Error(err))
		return nil, false
	}

	return list, true
}

func SetGreetingsByCharacterId(characterId int, greetings []string) bool {
	err := database.Transactional(func(ctx *database.TxContext) error {
		deleteQuery := "DELETE FROM character_greetings WHERE character_id = ?"
		if err := ctx.DeleteRecord(deleteQuery, []any{characterId}); err != nil {
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

	return err == nil
}

func CharacterGroupGreetingsByCharacterId(characterId int) ([]string, bool) {
	query := "SELECT text FROM character_group_greetings WHERE character_id = ?"
	args := []any{characterId}

	list, err := database.QueryForList(query, args, database.StringScanner)
	if err != nil {
		log.Get().Error("Error fetching character group greetings",
			zap.Int("characterId", characterId), zap.Error(err))
		return nil, false
	}

	return list, true
}

func SetGroupGreetingsByCharacterId(characterId int, greetings []string) bool {
	err := database.Transactional(func(ctx *database.TxContext) error {
		deleteQuery := "DELETE FROM character_group_greetings WHERE character_id = ?"
		if err := ctx.DeleteRecord(deleteQuery, []any{characterId}); err != nil {
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

	return err == nil
}

func RandomGreetingByCharacterId(characterId int, useGroupGreetings bool) (*string, bool) {
	query := `SELECT text FROM character_greetings WHERE character_id = ? ORDER BY RANDOM() LIMIT 1;`
	if useGroupGreetings {
		query = `SELECT text FROM character_group_greetings WHERE character_id = ? ORDER BY RANDOM() LIMIT 1;`
	}
	args := []any{characterId, useGroupGreetings}

	greeting, err := database.QueryForRecord(query, args, database.StringScanner)
	if err != nil {
		log.Get().Error("Error fetching random greeting",
			zap.Int("characterId", characterId), zap.Error(err))
		return nil, false
	}

	return greeting, true
}

func AllTags() ([]Tag, bool) {
	query := "SELECT * FROM tags"
	list, err := database.QueryForList(query, nil, tagScanner)
	if err != nil {
		log.Get().Error("Error fetching tags", zap.Error(err))
		return nil, false
	}

	return list, true
}

func TagById(id int) (*Tag, bool) {
	query := "SELECT * FROM tags WHERE id = ?"
	args := []any{id}
	tag, err := database.QueryForRecord(query, args, tagScanner)
	if err != nil {
		log.Get().Error("Error fetching tag",
			zap.Int("id", id), zap.Error(err))
		return nil, false
	}

	return tag, true
}

func CreateTag(newTag *Tag) bool {
	newTag.Lowercase = strings.ToLower(newTag.Label)

	query := "INSERT INTO tags(label, lowercase) VALUES (?, ?) RETURNING id"
	args := []any{newTag.Label, newTag.Lowercase}

	err := database.InsertRecord(query, args, &newTag.ID)
	if err != nil {
		log.Get().Error("Error inserting tag", zap.Error(err))
		return false
	}

	return true
}

func UpdateTag(id int, tag *Tag) bool {
	tag.Lowercase = strings.ToLower(tag.Label)

	query := "UPDATE tags SET label = ?, lowercase = ? WHERE id = ?"
	args := []any{id, tag.Label, tag.Lowercase}

	err := database.UpdateRecord(query, args)
	if err != nil {
		log.Get().Error("Error updating tag",
			zap.Int("id", id), zap.Error(err))
		return false
	}

	return true
}

func DeleteTagById(id int) bool {
	query := "DELETE FROM tags WHERE id = ?"
	args := []any{id}

	err := database.DeleteRecord(query, args)
	if err != nil {
		log.Get().Error("Error deleting tag",
			zap.Int("id", id), zap.Error(err))
		return false
	}

	return true
}
