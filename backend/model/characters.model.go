package model

import (
	"database/sql"
	"fmt"
	"juraji.nl/chat-quest/util"
	"time"
)

type Character struct {
	ID        int64      `json:"id"`
	CreatedAt *time.Time `json:"createdAt"`
	Name      string     `json:"name"`
	Favorite  bool       `json:"favorite"`
	AvatarUrl *string    `json:"avatarUrl"`
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

func characterScanner(scanner RowScanner, dest *Character) error {
	return scanner.Scan(
		&dest.ID,
		&dest.CreatedAt,
		&dest.Name,
		&dest.Favorite,
		&dest.AvatarUrl,
	)
}

func characterDetailsScanner(scanner RowScanner, dest *CharacterDetails) error {
	return scanner.Scan(
		&dest.CharacterId,
		&dest.Appearance,
		&dest.Personality,
		&dest.History,
		&dest.GroupTalkativeness,
	)
}

func AllCharacters(db *sql.DB) ([]*Character, error) {
	query := "SELECT * FROM characters"
	return queryForList(db, query, nil, characterScanner)
}

func CharacterById(db *sql.DB, id int64) (*Character, error) {
	query := "SELECT * FROM characters WHERE id = $1"
	args := []any{id}

	return queryForRecord(db, query, args, characterScanner)
}

func CreateCharacter(db *sql.DB, newCharacter *Character) error {
	util.EmptyStrPtrToNil(&newCharacter.AvatarUrl)

	query := "INSERT INTO characters (name, favorite, avatar_url) VALUES ($1, $2, $3) RETURNING id, created_at"
	args := []any{newCharacter.Name, newCharacter.Favorite, newCharacter.AvatarUrl}
	scanFunc := func(scanner RowScanner) error {
		return scanner.Scan(
			&newCharacter.ID,
			&newCharacter.CreatedAt,
		)
	}

	err := insertRecord(db, query, args, scanFunc)
	if err != nil {
		return err
	}

	// Create empty character details (so it exists when fetched)
	var newCharacterDetails CharacterDetails
	return UpdateCharacterDetails(db, newCharacter.ID, &newCharacterDetails)
}

func UpdateCharacter(db *sql.DB, id int64, character *Character) error {
	util.EmptyStrPtrToNil(&character.AvatarUrl)

	query := "UPDATE characters SET name = $1, favorite = $2, avatar_url = $3 WHERE id = $4"
	args := []any{character.Name, character.Favorite, character.AvatarUrl, id}

	return updateRecord(db, query, args)
}

func DeleteCharacterById(db *sql.DB, id int64) error {
	query := "DELETE FROM characters WHERE id = $1"
	args := []any{id}

	return deleteRecord(db, query, args)
}

func CharacterDetailsByCharacterId(db *sql.DB, characterId int64) (*CharacterDetails, error) {
	query := "SELECT * FROM character_details WHERE character_id = $1"
	args := []any{characterId}

	return queryForRecord(db, query, args, characterDetailsScanner)
}

func UpdateCharacterDetails(db *sql.DB, characterId int64, characterDetail *CharacterDetails) error {
	util.EmptyStrPtrToNil(&characterDetail.Appearance)
	util.EmptyStrPtrToNil(&characterDetail.Personality)
	util.EmptyStrPtrToNil(&characterDetail.History)

	//language=sqlite
	query := `
    INSERT OR REPLACE INTO character_details
      (character_id, appearance, personality, history, group_talkativeness)
    VALUES ($1, $2, $3, $4, $5)
  `
	args := []any{
		characterId,
		characterDetail.Appearance,
		characterDetail.Personality,
		characterDetail.History,
		characterDetail.GroupTalkativeness,
	}

	return updateRecord(db, query, args)
}

func TagsByCharacterId(db *sql.DB, characterId int64) ([]*Tag, error) {
	query := `
    SELECT t.*
    FROM character_tags ct
        JOIN tags t ON ct.tag_id = t.id
    WHERE ct.character_id = $1
  `
	args := []any{characterId}

	return queryForList(db, query, args, tagScanner)
}

func AddCharacterTag(db *sql.DB, characterId int64, tagId int64) error {
	query := "INSERT INTO character_tags (character_id, tag_id) VALUES ($1, $2)"
	args := []any{characterId, tagId}

	return updateRecord(db, query, args)
}

func RemoveCharacterTag(db *sql.DB, characterId int64, tagId int64) error {
	query := "DELETE FROM character_tags WHERE character_id = $1 AND tag_id = $2"
	args := []any{characterId, tagId}

	return deleteRecord(db, query, args)
}

func SetCharacterTags(db *sql.DB, characterId int64, tagIds []int64) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func(tx *sql.Tx, err error) {
		if err != nil {
			_ = tx.Rollback()
		}
	}(tx, err)

	deleteQuery := "DELETE FROM character_tags WHERE character_id = $1"
	if _, err := tx.Exec(deleteQuery, characterId); err != nil {
		return fmt.Errorf("failed to delete existing tag ids: %w", err)
	}

	if len(tagIds) == 0 {
		return tx.Commit()
	}

	insertQuery := "INSERT INTO character_tags (character_id, tag_id) VALUES ($1, $2)"
	for _, tagId := range tagIds {
		_, err := tx.Exec(insertQuery, characterId, tagId)
		if err != nil {
			return fmt.Errorf("failed to insert tag id: %w", err)
		}
	}

	return tx.Commit()
}

func DialogueExamplesByCharacterId(db *sql.DB, characterId int64) ([]*string, error) {
	query := "SELECT text FROM character_dialogue_examples WHERE character_id = $1"
	args := []any{characterId}
	scanFunc := func(rows RowScanner, dest *string) error {
		return rows.Scan(dest)
	}

	return queryForList(db, query, args, scanFunc)
}

func SetDialogueExamplesByCharacterId(db *sql.DB, characterId int64, examples []string) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func(tx *sql.Tx, err error) {
		if err != nil {
			_ = tx.Rollback()
		}
	}(tx, err)

	deleteQuery := "DELETE FROM character_dialogue_examples WHERE character_id = $1"
	if _, err := tx.Exec(deleteQuery, characterId); err != nil {
		return fmt.Errorf("failed to delete existing dialogue examples: %w", err)
	}

	if len(examples) == 0 {
		return tx.Commit()
	}

	insertQuery := "INSERT INTO character_dialogue_examples (character_id, text) VALUES ($1, $2)"
	for _, example := range examples {
		_, err := tx.Exec(insertQuery, characterId, example)
		if err != nil {
			return fmt.Errorf("failed to insert dialogue example: %w", err)
		}
	}

	return tx.Commit()
}

func CharacterGreetingsByCharacterId(db *sql.DB, characterId int64) ([]*string, error) {
	query := "SELECT * FROM character_greetings WHERE character_id = $1"
	args := []any{characterId}
	scanFunc := func(rows RowScanner, dest *string) error {
		return rows.Scan(dest)
	}

	return queryForList(db, query, args, scanFunc)
}

func SetGreetingsByCharacterId(db *sql.DB, characterId int64, greetings []string) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func(tx *sql.Tx, err error) {
		if err != nil {
			_ = tx.Rollback()
		}
	}(tx, err)

	deleteQuery := "DELETE FROM character_greetings WHERE character_id = $1"
	if _, err := tx.Exec(deleteQuery, characterId); err != nil {
		return fmt.Errorf("failed to delete existing greetings: %w", err)
	}

	if len(greetings) == 0 {
		return tx.Commit()
	}

	insertQuery := "INSERT INTO character_greetings (character_id, text) VALUES ($1, $2)"
	for _, greeting := range greetings {
		_, err := tx.Exec(insertQuery, characterId, greeting)
		if err != nil {
			return fmt.Errorf("failed to insert greeting: %w", err)
		}
	}

	return tx.Commit()
}

func CharacterGroupGreetingsByCharacterId(db *sql.DB, characterId int64) ([]*string, error) {
	query := "SELECT * FROM character_group_greetings WHERE character_id = $1"
	args := []any{characterId}
	scanFunc := func(rows RowScanner, dest *string) error {
		return rows.Scan(dest)
	}

	return queryForList(db, query, args, scanFunc)
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

	deleteQuery := "DELETE FROM character_greetings WHERE character_id = $1"
	if _, err := tx.Exec(deleteQuery, characterId); err != nil {
		return fmt.Errorf("failed to delete existing greetings: %w", err)
	}

	if len(greetings) == 0 {
		return tx.Commit()
	}

	insertQuery := "INSERT INTO character_greetings (character_id, text) VALUES ($1, $2)"
	for _, greeting := range greetings {
		_, err := tx.Exec(insertQuery, characterId, greeting)
		if err != nil {
			return fmt.Errorf("failed to insert greeting: %w", err)
		}
	}

	return tx.Commit()
}
