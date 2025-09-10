package characters

import (
	"strings"

	"github.com/pkg/errors"
	"juraji.nl/chat-quest/core/database"
)

type Tag struct {
	ID        int    `json:"id"`
	Label     string `json:"label"`
	Lowercase string `json:"lowercase"`
}

func tagScanner(scanner database.RowScanner, dest *Tag) error {
	return scanner.Scan(
		&dest.ID,
		&dest.Label,
		&dest.Lowercase,
	)
}

func AllTags() ([]Tag, error) {
	query := "SELECT * FROM tags"
	return database.QueryForList(query, nil, tagScanner)
}

func TagById(id int) (*Tag, error) {
	query := "SELECT * FROM tags WHERE id = ?"
	args := []any{id}
	return database.QueryForRecord(query, args, tagScanner)
}

func CreateTag(newTag *Tag) error {
	newTag.Lowercase = strings.ToLower(newTag.Label)

	query := "INSERT INTO tags(label, lowercase) VALUES (?, ?) RETURNING id"
	args := []any{newTag.Label, newTag.Lowercase}

	err := database.InsertRecord(query, args, &newTag.ID)

	if err == nil {
		TagCreatedSignal.EmitBG(newTag)
	}

	return err
}

func UpdateTag(id int, tag *Tag) error {
	tag.Lowercase = strings.ToLower(tag.Label)

	query := "UPDATE tags SET label = ?, lowercase = ? WHERE id = ?"
	args := []any{id, tag.Label, tag.Lowercase}

	err := database.UpdateRecord(query, args)

	if err == nil {
		TagUpdatedSignal.EmitBG(tag)
	}

	return err
}

func DeleteTagById(id int) error {
	query := "DELETE FROM tags WHERE id = ?"
	args := []any{id}

	_, err := database.DeleteRecord(query, args)

	if err == nil {
		TagDeletedSignal.EmitBG(id)
	}

	return err
}

func TagsByCharacterId(characterId int) ([]Tag, error) {
	query := `
    SELECT t.*
    FROM character_tags ct
        JOIN tags t ON ct.tag_id = t.id
    WHERE ct.character_id = ?
  `
	args := []any{characterId}

	return database.QueryForList(query, args, tagScanner)
}

func AddCharacterTag(characterId int, tagId int) error {
	query := "INSERT INTO character_tags (character_id, tag_id) VALUES (?, ?)"
	args := []any{characterId, tagId}

	err := database.UpdateRecord(query, args)
	if err == nil {
		CharacterTagAddedSignal.EmitBG([]int{characterId, tagId})
	}
	return err
}

func RemoveCharacterTag(characterId int, tagId int) error {
	query := "DELETE FROM character_tags WHERE character_id = ? AND tag_id = ?"
	args := []any{characterId, tagId}

	_, err := database.DeleteRecord(query, args)

	if err == nil {
		CharacterTagRemovedSignal.EmitBG([]int{characterId, tagId})
	}

	return err
}
func SetCharacterTags(characterId int, tagIds []int) error {
	var deletedTagIds []int

	err := database.Transactional(func(ctx *database.TxContext) error {
		var err error
		deleteQuery := "DELETE FROM character_tags WHERE character_id = ? RETURNING tag_id"
		deletedTagIds, err = ctx.DeleteRecord(deleteQuery, []any{characterId})
		if err != nil {
			return errors.Wrapf(err, "error removing tags for character %d", characterId)
		}

		if len(tagIds) == 0 {
			// Shortcut, no more work to do
			return nil
		}

		insertQuery := "INSERT INTO character_tags (character_id, tag_id) VALUES (?, ?)"
		for _, tagId := range tagIds {
			if err = ctx.InsertRecord(insertQuery, []any{characterId, tagId}); err != nil {
				return errors.Wrapf(err, "error inserting tags for character %d", characterId)
			}
		}

		return nil
	})

	if err == nil {
		for i := range deletedTagIds {
			rec := []int{characterId, deletedTagIds[i]}
			CharacterTagRemovedSignal.EmitBG(rec)
		}

		for i := range tagIds {
			rec := []int{characterId, tagIds[i]}
			CharacterTagAddedSignal.EmitBG(rec)
		}
	}

	return err
}
