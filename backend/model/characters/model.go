package characters

import (
	"database/sql"
	"fmt"
	"juraji.nl/chat-quest/core"
	"juraji.nl/chat-quest/core/database"
	"juraji.nl/chat-quest/core/util"
	"strings"
	"time"
)

type Character struct {
	ID        int        `json:"id"`
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
	CharacterId        int     `json:"characterId"`
	Appearance         *string `json:"appearance"`
	Personality        *string `json:"personality"`
	History            *string `json:"history"`
	GroupTalkativeness float32 `json:"groupTalkativeness"`
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

func AllCharacters(cq *core.ChatQuestContext) ([]*Character, error) {
	query := "SELECT * FROM characters"
	return database.QueryForList(cq.DB(), query, nil, CharacterScanner)
}

func AllCharactersWithTags(cq *core.ChatQuestContext) ([]*CharacterWithTags, error) {
	query := "SELECT * FROM characters"
	characters, err := database.QueryForList(cq.DB(), query, nil, CharacterScanner)
	if err != nil {
		return nil, err
	}

	var charactersWithTags []*CharacterWithTags
	for _, character := range characters {
		tags, err := TagsByCharacterId(cq, character.ID)
		if err != nil {
			return nil, err
		}
		charactersWithTags = append(charactersWithTags, &CharacterWithTags{*character, tags})
	}

	return charactersWithTags, nil
}

func CharacterById(cq *core.ChatQuestContext, id int) (*Character, error) {
	query := "SELECT * FROM characters WHERE id = ?"
	args := []any{id}

	return database.QueryForRecord(cq.DB(), query, args, CharacterScanner)
}

func CreateCharacter(cq *core.ChatQuestContext, newCharacter *Character) error {
	tx, err := cq.DB().Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer database.RollBackOnErr(cq, tx, err)

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

	CharacterCreatedSignal.Emit(cq.Context(), newCharacter)
	return nil
}

func UpdateCharacter(cq *core.ChatQuestContext, id int, character *Character) error {
	util.EmptyStrPtrToNil(&character.AvatarUrl)

	query := "UPDATE characters SET name = ?, favorite = ?, avatar_url = ? WHERE id = ?"
	args := []any{character.Name, character.Favorite, character.AvatarUrl, id}

	err := database.UpdateRecord(cq.DB(), query, args)
	if err != nil {
		return err
	}

	CharacterUpdatedSignal.Emit(cq.Context(), character)
	return nil
}

func DeleteCharacterById(cq *core.ChatQuestContext, id int) error {
	query := "DELETE FROM characters WHERE id = ?"
	args := []any{id}

	err := database.DeleteRecord(cq.DB(), query, args)
	if err != nil {
		return err
	}

	CharacterDeletedSignal.Emit(cq.Context(), id)
	return nil
}

func CharacterDetailsByCharacterId(cq *core.ChatQuestContext, characterId int) (*CharacterDetails, error) {
	query := "SELECT * FROM character_details WHERE character_id = ?"
	args := []any{characterId}

	return database.QueryForRecord(cq.DB(), query, args, characterDetailsScanner)
}

func UpdateCharacterDetails(cq *core.ChatQuestContext, characterId int, characterDetail *CharacterDetails) error {
	return updateCharacterDetails(cq.DB(), characterId, characterDetail)
}

func updateCharacterDetails(db database.QueryExecutor, characterId int, characterDetail *CharacterDetails) error {
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

func TagsByCharacterId(cq *core.ChatQuestContext, characterId int) ([]*Tag, error) {
	query := `
    SELECT t.*
    FROM character_tags ct
        JOIN tags t ON ct.tag_id = t.id
    WHERE ct.character_id = ?
  `
	args := []any{characterId}

	return database.QueryForList(cq.DB(), query, args, tagScanner)
}

func AddCharacterTag(cq *core.ChatQuestContext, characterId int, tagId int) error {
	query := "INSERT INTO character_tags (character_id, tag_id) VALUES (?, ?)"
	args := []any{characterId, tagId}

	return database.UpdateRecord(cq.DB(), query, args)
}

func RemoveCharacterTag(cq *core.ChatQuestContext, characterId int, tagId int) error {
	query := "DELETE FROM character_tags WHERE character_id = ? AND tag_id = ?"
	args := []any{characterId, tagId}

	return database.DeleteRecord(cq.DB(), query, args)
}

func SetCharacterTags(cq *core.ChatQuestContext, characterId int, tagIds []int) error {
	tx, err := cq.DB().Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer database.RollBackOnErr(cq, tx, err)

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

func DialogueExamplesByCharacterId(cq *core.ChatQuestContext, characterId int) ([]*string, error) {
	query := "SELECT text FROM character_dialogue_examples WHERE character_id = ?"
	args := []any{characterId}
	scanFunc := func(rows database.RowScanner, dest *string) error {
		return rows.Scan(dest)
	}

	return database.QueryForList(cq.DB(), query, args, scanFunc)
}

func SetDialogueExamplesByCharacterId(cq *core.ChatQuestContext, characterId int, examples []string) error {
	tx, err := cq.DB().Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer database.RollBackOnErr(cq, tx, err)

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

func CharacterGreetingsByCharacterId(cq *core.ChatQuestContext, characterId int) ([]*string, error) {
	query := "SELECT text FROM character_greetings WHERE character_id = ?"
	args := []any{characterId}
	scanFunc := func(rows database.RowScanner, dest *string) error {
		return rows.Scan(dest)
	}

	return database.QueryForList(cq.DB(), query, args, scanFunc)
}

func SetGreetingsByCharacterId(cq *core.ChatQuestContext, characterId int, greetings []string) error {
	tx, err := cq.DB().Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer database.RollBackOnErr(cq, tx, err)

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

func CharacterGroupGreetingsByCharacterId(cq *core.ChatQuestContext, characterId int) ([]*string, error) {
	query := "SELECT text FROM character_group_greetings WHERE character_id = ?"
	args := []any{characterId}
	scanFunc := func(rows database.RowScanner, dest *string) error {
		return rows.Scan(dest)
	}

	return database.QueryForList(cq.DB(), query, args, scanFunc)
}

func SetGroupGreetingsByCharacterId(cq *core.ChatQuestContext, characterId int, greetings []string) error {
	tx, err := cq.DB().Begin()
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

func AllTags(cq *core.ChatQuestContext) ([]*Tag, error) {
	query := "SELECT * FROM tags"
	return database.QueryForList(cq.DB(), query, nil, tagScanner)
}

func TagById(cq *core.ChatQuestContext, id int) (*Tag, error) {
	query := "SELECT * FROM tags WHERE id = ?"
	args := []any{id}
	return database.QueryForRecord(cq.DB(), query, args, tagScanner)
}

func CreateTag(cq *core.ChatQuestContext, newTag *Tag) error {
	newTag.Lowercase = strings.ToLower(newTag.Label)

	query := "INSERT INTO tags(label, lowercase) VALUES (?, ?) RETURNING id"
	args := []any{newTag.Label, newTag.Lowercase}

	return database.InsertRecord(cq.DB(), query, args, &newTag.ID)
}

func UpdateTag(cq *core.ChatQuestContext, id int, tag *Tag) error {
	tag.Lowercase = strings.ToLower(tag.Label)

	query := "UPDATE tags SET label = ?, lowercase = ? WHERE id = ?"
	args := []any{id, tag.Label, tag.Lowercase}

	return database.UpdateRecord(cq.DB(), query, args)
}

func DeleteTagById(cq *core.ChatQuestContext, id int) error {
	query := "DELETE FROM tags WHERE id = ?"
	args := []any{id}

	return database.DeleteRecord(cq.DB(), query, args)
}
