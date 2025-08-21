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

func TagsByCharacterId(characterId int) ([]Tag, bool) {
	query := `
    SELECT t.*
    FROM character_tags ct
        JOIN tags t ON ct.tag_id = t.id
    WHERE ct.character_id = ?
  `
	args := []any{characterId}

	list, err := database.QueryForList(query, args, tagScanner)
	if err != nil {
		log.Get().Error("Error fetching character tags",
			zap.Int("characterId", characterId), zap.Error(err))
		return list, false
	}

	return list, true
}

func AddCharacterTag(characterId int, tagId int) bool {
	query := "INSERT INTO character_tags (character_id, tag_id) VALUES (?, ?)"
	args := []any{characterId, tagId}

	err := database.UpdateRecord(query, args)
	if err != nil {
		log.Get().Error("Error adding tag",
			zap.Int("characterId", characterId),
			zap.Int("tagId", tagId),
			zap.Error(err))
		return false
	}

	return true
}

func RemoveCharacterTag(characterId int, tagId int) bool {
	query := "DELETE FROM character_tags WHERE character_id = ? AND tag_id = ?"
	args := []any{characterId, tagId}

	err := database.DeleteRecord(query, args)
	if err != nil {
		log.Get().Error("Error removing tag",
			zap.Int("characterId", characterId),
			zap.Int("tagId", tagId),
			zap.Error(err))
		return false
	}

	return true
}

func SetCharacterTags(characterId int, tagIds []int) bool {
	err := database.Transactional(func(ctx *database.TxContext) error {
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

	return err == nil
}
