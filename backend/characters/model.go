package characters

import (
	"database/sql"
	"fmt"
	"juraji.nl/chat-quest/database"
	"juraji.nl/chat-quest/util"
	"strings"
	"time"
)

type Character struct {
	ID        int64      `json:"id"`
	CreatedAt *time.Time `json:"createdAt"`
	Name      string     `json:"name"`
	Favorite  bool       `json:"favorite"`
	AvatarUrl *string    `json:"avatarUrl"`
}

type CharacterWithTags struct {
	Character
	Tags []*Tag `json:"tags"`
}

type CharacterDetails struct {
	CharacterId        int64   `json:"characterId"`
	Appearance         *string `json:"appearance"`
	Personality        *string `json:"personality"`
	History            *string `json:"history"`
	GroupTalkativeness float64 `json:"groupTalkativeness"`
}

type CharacterTextBlock struct {
	CharacterId int64  `json:"characterId"`
	Text        string `json:"text"`
}

type Tag struct {
	ID        int    `json:"id"`
	Label     string `json:"label"`
	Lowercase string `json:"lowercase"`
}

func characterScanner(scanner database.RowScanner, dest *Character) error {
	return scanner.Scan(
		&dest.ID,
		&dest.CreatedAt,
		&dest.Name,
		&dest.Favorite,
		&dest.AvatarUrl,
	)
}
func characterDetailsScanner(scanner database.RowScanner, dest *CharacterDetails) error {
	return scanner.Scan(
		&dest.CharacterId,
		&dest.Appearance,
		&dest.Personality,
		&dest.History,
		&dest.GroupTalkativeness,
	)
}

func tagScanner(scanner database.RowScanner, dest *Tag) error {
	return scanner.Scan(
		&dest.ID,
		&dest.Label,
		&dest.Lowercase,
	)
}

func AllCharacters(db *sql.DB) ([]*Character, error) {
	query := "SELECT * FROM characters"
	return database.QueryForList(db, query, nil, characterScanner)
}

func AllCharactersWithTags(db *sql.DB) ([]*CharacterWithTags, error) {
	query := "SELECT * FROM characters"
	characters, err := database.QueryForList(db, query, nil, characterScanner)
	if err != nil {
		return nil, err
	}

	var charactersWithTags []*CharacterWithTags
	for _, character := range characters {
		tags, err := TagsByCharacterId(db, character.ID)
		if err != nil {
			return nil, err
		}
		charactersWithTags = append(charactersWithTags, &CharacterWithTags{*character, tags})
	}

	return charactersWithTags, nil
}

func CharacterById(db *sql.DB, id int64) (*Character, error) {
	query := "SELECT * FROM characters WHERE id = ?"
	args := []any{id}

	return database.QueryForRecord(db, query, args, characterScanner)
}

func CreateCharacter(db *sql.DB, newCharacter *Character) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer database.RollBackOnErr(tx, err)
	defer util.EmitOnSuccess(CharacterCreatedSignal, newCharacter, err)

	util.EmptyStrPtrToNil(&newCharacter.AvatarUrl)

	query := "INSERT INTO characters (name, favorite, avatar_url) VALUES (?, ?, ?) RETURNING id, created_at"
	args := []any{newCharacter.Name, newCharacter.Favorite, newCharacter.AvatarUrl}

	if err := database.InsertRecord(tx, query, args, &newCharacter.ID, &newCharacter.CreatedAt); err != nil {
		return err
	}

	// Create empty character details (so it exists when fetched)
	var newCharacterDetails CharacterDetails
	if err = updateCharacterDetails(tx, newCharacter.ID, &newCharacterDetails); err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func UpdateCharacter(db *sql.DB, id int64, character *Character) error {
	util.EmptyStrPtrToNil(&character.AvatarUrl)

	query := "UPDATE characters SET name = ?, favorite = ?, avatar_url = ? WHERE id = ?"
	args := []any{character.Name, character.Favorite, character.AvatarUrl, id}

	err := database.UpdateRecord(db, query, args)
	defer util.EmitOnSuccess(CharacterUpdatedSignal, character, err)

	return err
}

func DeleteCharacterById(db *sql.DB, id int64) error {
	query := "DELETE FROM characters WHERE id = ?"
	args := []any{id}

	err := database.DeleteRecord(db, query, args)
	defer util.EmitOnSuccess(CharacterDeletedSignal, id, err)
	return err
}

func CharacterDetailsByCharacterId(db *sql.DB, characterId int64) (*CharacterDetails, error) {
	query := "SELECT * FROM character_details WHERE character_id = ?"
	args := []any{characterId}

	return database.QueryForRecord(db, query, args, characterDetailsScanner)
}

func UpdateCharacterDetails(db *sql.DB, characterId int64, characterDetail *CharacterDetails) error {
	return updateCharacterDetails(db, characterId, characterDetail)
}

func updateCharacterDetails(db database.QueryExecutor, characterId int64, characterDetail *CharacterDetails) error {
	util.EmptyStrPtrToNil(&characterDetail.Appearance)
	util.EmptyStrPtrToNil(&characterDetail.Personality)
	util.EmptyStrPtrToNil(&characterDetail.History)

	//language=sqlite
	query := `
    INSERT OR REPLACE INTO character_details
      (character_id, appearance, personality, history, group_talkativeness)
    VALUES (?,?,?,?,?)
  `
	args := []any{
		characterId,
		characterDetail.Appearance,
		characterDetail.Personality,
		characterDetail.History,
		characterDetail.GroupTalkativeness,
	}

	return database.UpdateRecord(db, query, args)
}

func TagsByCharacterId(db *sql.DB, characterId int64) ([]*Tag, error) {
	query := `
    SELECT t.*
    FROM character_tags ct
        JOIN tags t ON ct.tag_id = t.id
    WHERE ct.character_id = ?
  `
	args := []any{characterId}

	return database.QueryForList(db, query, args, tagScanner)
}

func AddCharacterTag(db *sql.DB, characterId int64, tagId int64) error {
	query := "INSERT INTO character_tags (character_id, tag_id) VALUES (?, ?)"
	args := []any{characterId, tagId}

	return database.UpdateRecord(db, query, args)
}

func RemoveCharacterTag(db *sql.DB, characterId int64, tagId int64) error {
	query := "DELETE FROM character_tags WHERE character_id = ? AND tag_id = ?"
	args := []any{characterId, tagId}

	return database.DeleteRecord(db, query, args)
}

func SetCharacterTags(db *sql.DB, characterId int64, tagIds []int64) error {
	tx, err := db.Begin()
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

func DialogueExamplesByCharacterId(db *sql.DB, characterId int64) ([]*string, error) {
	query := "SELECT text FROM character_dialogue_examples WHERE character_id = ?"
	args := []any{characterId}
	scanFunc := func(rows database.RowScanner, dest *string) error {
		return rows.Scan(dest)
	}

	return database.QueryForList(db, query, args, scanFunc)
}

func SetDialogueExamplesByCharacterId(db *sql.DB, characterId int64, examples []string) error {
	tx, err := db.Begin()
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

func CharacterGreetingsByCharacterId(db *sql.DB, characterId int64) ([]*string, error) {
	query := "SELECT text FROM character_greetings WHERE character_id = ?"
	args := []any{characterId}
	scanFunc := func(rows database.RowScanner, dest *string) error {
		return rows.Scan(dest)
	}

	return database.QueryForList(db, query, args, scanFunc)
}

func SetGreetingsByCharacterId(db *sql.DB, characterId int64, greetings []string) error {
	tx, err := db.Begin()
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

func CharacterGroupGreetingsByCharacterId(db *sql.DB, characterId int64) ([]*string, error) {
	query := "SELECT text FROM character_group_greetings WHERE character_id = ?"
	args := []any{characterId}
	scanFunc := func(rows database.RowScanner, dest *string) error {
		return rows.Scan(dest)
	}

	return database.QueryForList(db, query, args, scanFunc)
}

func SetGroupGreetingsByCharacterId(db *sql.DB, characterId int64, greetings []string) error {
	tx, err := db.Begin()
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

func AllTags(db *sql.DB) ([]*Tag, error) {
	query := "SELECT * FROM tags"
	return database.QueryForList(db, query, nil, tagScanner)
}

func TagById(db *sql.DB, id int64) (*Tag, error) {
	query := "SELECT * FROM tags WHERE id = ?"
	args := []any{id}
	return database.QueryForRecord(db, query, args, tagScanner)
}

func CreateTag(db *sql.DB, newTag *Tag) error {
	newTag.Lowercase = strings.ToLower(newTag.Label)

	query := "INSERT INTO tags(label, lowercase) VALUES (?, ?) RETURNING id"
	args := []any{newTag.Label, newTag.Lowercase}

	return database.InsertRecord(db, query, args, &newTag.ID)
}

func UpdateTag(db *sql.DB, id int64, tag *Tag) error {
	tag.Lowercase = strings.ToLower(tag.Label)

	query := "UPDATE tags SET label = ?, lowercase = ? WHERE id = ?"
	args := []any{id, tag.Label, tag.Lowercase}

	return database.UpdateRecord(db, query, args)
}

func DeleteTagById(db *sql.DB, id int64) error {
	query := "DELETE FROM tags WHERE id = ?"
	args := []any{id}

	return database.DeleteRecord(db, query, args)
}
