package characters

import (
	"go.uber.org/zap"
	"juraji.nl/chat-quest/core/database"
	"juraji.nl/chat-quest/core/log"
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

	return database.UpdateRecord(query, args)
}

func RemoveCharacterTag(characterId int, tagId int) error {
	query := "DELETE FROM character_tags WHERE character_id = ? AND tag_id = ?"
	args := []any{characterId, tagId}

	return database.DeleteRecord(query, args)
}

func SetCharacterTags(characterId int, tagIds []int) error {
	return database.Transactional(func(ctx *database.TxContext) error {
		deleteQuery := "DELETE FROM character_tags WHERE character_id = ?"
		if err := ctx.DeleteRecord(deleteQuery, []any{characterId}); err != nil {
			log.Get().Error("Error removing tags", zap.Int("characterId", characterId), zap.Error(err))
			return err
		}

		if len(tagIds) == 0 {
			// Shortcut, no more work to do
			return nil
		}

		insertQuery := "INSERT INTO character_tags (character_id, tag_id) VALUES (?, ?)"
		for _, tagId := range tagIds {
			if err := ctx.InsertRecord(insertQuery, []any{characterId, tagId}); err != nil {
				log.Get().Error("Error adding tags",
					zap.Int("characterId", characterId),
					zap.Int("tagId", tagId),
					zap.Error(err))
				return err
			}
		}

		return nil
	})
}
