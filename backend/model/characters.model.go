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
