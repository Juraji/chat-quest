package model

import (
	"database/sql"
)

type Character struct {
	ID        int32   `json:"id"`
	CreatedAt string  `json:"createdAt"`
	Name      string  `json:"name"`
	Favorite  bool    `json:"favorite"`
	AvatarUrl *string `json:"avatarUrl"`
}

func AllCharacters(db *sql.DB) ([]*Character, error) {
	query := "SELECT id, name, favorite, created_at, avatar_url FROM characters"
	scanFunc := func(rows *sql.Rows, dest *Character) error {
		return rows.Scan(
			&dest.ID,
			&dest.Name,
			&dest.Favorite,
			&dest.CreatedAt,
			&dest.AvatarUrl,
		)
	}

	return queryForList[Character](db, query, scanFunc)
}

func CharacterById(db *sql.DB, id int32) (*Character, error) {
	query := "SELECT id, name, favorite, created_at, avatar_url FROM characters WHERE id = $1"
	args := []any{id}
	scanFunc := func(row *sql.Row, dest *Character) error {
		return row.Scan(
			&dest.ID,
			&dest.Name,
			&dest.Favorite,
			&dest.CreatedAt,
			&dest.AvatarUrl,
		)
	}

	return queryForRecord[Character](db, query, args, scanFunc)
}

func CreateCharacter(db *sql.DB, newCharacter *Character) error {
	query := "INSERT INTO characters (name, favorite, avatar_url) VALUES ($1, $2, $3) RETURNING id, created_at"
	args := []any{newCharacter.Name, newCharacter.Favorite, newCharacter.AvatarUrl}
	scanFunc := func(row *sql.Row) error {
		return row.Scan(
			&newCharacter.ID,
			&newCharacter.CreatedAt,
		)
	}

	return insertRecord(db, query, args, scanFunc)
}

func UpdateCharacter(db *sql.DB, id int32, character *Character) error {
	query := "UPDATE characters SET name = $1, favorite = $2, avatar_url = $3 WHERE id = $4"
	args := []any{character.Name, character.Favorite, character.AvatarUrl, id}

	return updateRecord(db, query, args)
}

func DeleteCharacterById(db *sql.DB, id int32) error {
	query := "DELETE FROM characters WHERE id = $1"
	args := []any{id}

	return deleteRecord(db, query, args)
}
