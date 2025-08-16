package characters

import (
	"database/sql"
	"fmt"
	"juraji.nl/chat-quest/core/database"
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
	Tags      []Tag      `json:"tags"`
}

type CharacterTextBlock struct {
	CharacterId int    `json:"characterId"`
	Text        string `json:"text"`
}

type Tag struct {
	ID        int    `json:"id"`
	Label     string `json:"label"`
	Lowercase string `json:"lowercase"`
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
		&dest.Tags,
	)
}

func tagScanner(scanner database.RowScanner, dest *Tag) error {
	return scanner.Scan(
		&dest.ID,
		&dest.Label,
		&dest.Lowercase,
	)
}

func AllCharacterListViews() ([]CharacterListView, error) {
	query := `SELECT
                c.id,
                c.created_at,
                c.name,
                c.favorite,
                c.avatar_url,
                GROUP_CONCAT(json_object('id', t.id, 'label', t.label)) as tags
            FROM characters c
            LEFT JOIN character_tags ct ON c.id = ct.character_id
            LEFT JOIN tags t ON ct.tag_id = t.id
            GROUP BY c.id;`

	return database.QueryForList(database.GetDB(), query, nil, characterListViewScanner)
}

func CharacterById(id int) (*Character, error) {
	query := "SELECT * FROM characters WHERE id = ?"
	args := []any{id}

	return database.QueryForRecord(database.GetDB(), query, args, CharacterScanner)
}

func CreateCharacter(newCharacter *Character) error {
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

	if err := database.InsertRecord(database.GetDB(), query, args, &newCharacter.ID, &newCharacter.CreatedAt); err != nil {
		return err
	}

	util.Emit(CharacterCreatedSignal, newCharacter)
	return nil
}

func UpdateCharacter(id int, character *Character) error {
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

	err := database.UpdateRecord(database.GetDB(), query, args)
	if err != nil {
		return err
	}

	util.Emit(CharacterUpdatedSignal, character)
	return nil
}

func DeleteCharacterById(id int) error {
	query := "DELETE FROM characters WHERE id = ?"
	args := []any{id}

	err := database.DeleteRecord(database.GetDB(), query, args)
	if err != nil {
		return err
	}

	util.Emit(CharacterDeletedSignal, id)
	return nil
}

func TagsByCharacterId(characterId int) ([]Tag, error) {
	query := `
    SELECT t.*
    FROM character_tags ct
        JOIN tags t ON ct.tag_id = t.id
    WHERE ct.character_id = ?
  `
	args := []any{characterId}

	return database.QueryForList(database.GetDB(), query, args, tagScanner)
}

func AddCharacterTag(characterId int, tagId int) error {
	query := "INSERT INTO character_tags (character_id, tag_id) VALUES (?, ?)"
	args := []any{characterId, tagId}

	return database.UpdateRecord(database.GetDB(), query, args)
}

func RemoveCharacterTag(characterId int, tagId int) error {
	query := "DELETE FROM character_tags WHERE character_id = ? AND tag_id = ?"
	args := []any{characterId, tagId}

	return database.DeleteRecord(database.GetDB(), query, args)
}

func SetCharacterTags(characterId int, tagIds []int) error {
	tx, err := database.GetDB().Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer database.RollBackOnErr(tx, err)

	deleteQuery := "DELETE FROM character_tags WHERE character_id = ?"
	if err := database.DeleteRecord(tx, deleteQuery, []any{characterId}); err != nil {
		return fmt.Errorf("failed to delete existing tag ids: %w", err)
	}

	if len(tagIds) == 0 {
		return tx.Commit()
	}

	insertQuery := "INSERT INTO character_tags (character_id, tag_id) VALUES (?, ?)"
	for _, tagId := range tagIds {
		if err := database.InsertRecord(tx, insertQuery, []any{characterId, tagId}); err != nil {
			return fmt.Errorf("failed to insert tag id: %w", err)
		}
	}

	return tx.Commit()
}

func DialogueExamplesByCharacterId(characterId int) ([]string, error) {
	query := "SELECT text FROM character_dialogue_examples WHERE character_id = ?"
	args := []any{characterId}
	scanFunc := func(rows database.RowScanner, dest *string) error {
		return rows.Scan(dest)
	}

	return database.QueryForList(database.GetDB(), query, args, scanFunc)
}

func SetDialogueExamplesByCharacterId(characterId int, examples []string) error {
	tx, err := database.GetDB().Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer database.RollBackOnErr(tx, err)

	deleteQuery := "DELETE FROM character_dialogue_examples WHERE character_id = ?"
	if err := database.DeleteRecord(tx, deleteQuery, []any{characterId}); err != nil {
		return fmt.Errorf("failed to delete existing dialogue examples: %w", err)
	}

	if len(examples) == 0 {
		return tx.Commit()
	}

	insertQuery := "INSERT INTO character_dialogue_examples (character_id, text) VALUES (?, ?)"
	for _, example := range examples {
		if err := database.InsertRecord(tx, insertQuery, []any{characterId, example}); err != nil {
			return fmt.Errorf("failed to insert dialogue example: %w", err)
		}
	}

	return tx.Commit()
}

func CharacterGreetingsByCharacterId(characterId int) ([]string, error) {
	query := "SELECT text FROM character_greetings WHERE character_id = ?"
	args := []any{characterId}
	scanFunc := func(rows database.RowScanner, dest *string) error {
		return rows.Scan(dest)
	}

	return database.QueryForList(database.GetDB(), query, args, scanFunc)
}

func SetGreetingsByCharacterId(characterId int, greetings []string) error {
	tx, err := database.GetDB().Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer database.RollBackOnErr(tx, err)

	deleteQuery := "DELETE FROM character_greetings WHERE character_id = ?"
	if err := database.DeleteRecord(tx, deleteQuery, []any{characterId}); err != nil {
		return fmt.Errorf("failed to delete existing greetings: %w", err)
	}

	if len(greetings) == 0 {
		return tx.Commit()
	}

	insertQuery := "INSERT INTO character_greetings (character_id, text) VALUES (?, ?)"
	for _, greeting := range greetings {
		if err := database.InsertRecord(tx, insertQuery, []any{characterId, greeting}); err != nil {
			return fmt.Errorf("failed to insert greeting: %w", err)
		}
	}

	return tx.Commit()
}

func CharacterGroupGreetingsByCharacterId(characterId int) ([]string, error) {
	query := "SELECT text FROM character_group_greetings WHERE character_id = ?"
	args := []any{characterId}
	scanFunc := func(rows database.RowScanner, dest *string) error {
		return rows.Scan(dest)
	}

	return database.QueryForList(database.GetDB(), query, args, scanFunc)
}

func SetGroupGreetingsByCharacterId(characterId int, greetings []string) error {
	tx, err := database.GetDB().Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func(tx *sql.Tx, err error) {
		if err != nil {
			_ = tx.Rollback()
		}
	}(tx, err)

	deleteQuery := "DELETE FROM character_greetings WHERE character_id = ?"
	if err := database.DeleteRecord(tx, deleteQuery, []any{characterId}); err != nil {
		return fmt.Errorf("failed to delete existing greetings: %w", err)
	}

	if len(greetings) == 0 {
		return tx.Commit()
	}

	insertQuery := "INSERT INTO character_greetings (character_id, text) VALUES (?, ?)"
	for _, greeting := range greetings {
		if err := database.InsertRecord(tx, insertQuery, []any{characterId, greeting}); err != nil {
			return fmt.Errorf("failed to insert greeting: %w", err)
		}
	}

	return tx.Commit()
}

func RandomGreetingByCharacterId(characterId int, useGroupGreetings bool) (*string, error) {
	query := `SELECT text FROM character_greetings WHERE character_id = ? ORDER BY RANDOM() LIMIT 1;`
	if useGroupGreetings {
		query = `SELECT text FROM character_group_greetings WHERE character_id = ? ORDER BY RANDOM() LIMIT 1;`
	}
	args := []any{characterId, useGroupGreetings}
	return database.QueryForRecord(database.GetDB(), query, args, database.StringScanner)
}

func AllTags() ([]Tag, error) {
	query := "SELECT * FROM tags"
	return database.QueryForList(database.GetDB(), query, nil, tagScanner)
}

func TagById(id int) (*Tag, error) {
	query := "SELECT * FROM tags WHERE id = ?"
	args := []any{id}
	return database.QueryForRecord(database.GetDB(), query, args, tagScanner)
}

func CreateTag(newTag *Tag) error {
	newTag.Lowercase = strings.ToLower(newTag.Label)

	query := "INSERT INTO tags(label, lowercase) VALUES (?, ?) RETURNING id"
	args := []any{newTag.Label, newTag.Lowercase}

	return database.InsertRecord(database.GetDB(), query, args, &newTag.ID)
}

func UpdateTag(id int, tag *Tag) error {
	tag.Lowercase = strings.ToLower(tag.Label)

	query := "UPDATE tags SET label = ?, lowercase = ? WHERE id = ?"
	args := []any{id, tag.Label, tag.Lowercase}

	return database.UpdateRecord(database.GetDB(), query, args)
}

func DeleteTagById(id int) error {
	query := "DELETE FROM tags WHERE id = ?"
	args := []any{id}

	return database.DeleteRecord(database.GetDB(), query, args)
}
