package model

import (
	"database/sql"
	"strings"
)

type Tag struct {
	ID        int    `json:"id"`
	Label     string `json:"label"`
	Lowercase string `json:"lowercase"`
}

func AllTags(db *sql.DB) ([]*Tag, error) {
	query := "SELECT * FROM tags"
	scanFunc := func(rows *sql.Rows, dest *Tag) error {
		return rows.Scan(
			&dest.ID,
			&dest.Label,
			&dest.Lowercase,
		)
	}

	return queryForList(db, query, scanFunc)
}

func TagById(db *sql.DB, id int32) (*Tag, error) {
	query := "SELECT * FROM tags WHERE id = $1"
	args := []any{id}
	scanFunc := func(row *sql.Row, dest *Tag) error {
		return row.Scan(
			&dest.ID,
			&dest.Label,
			&dest.Lowercase,
		)
	}

	return queryForRecord(db, query, args, scanFunc)
}

func CreateTag(db *sql.DB, newTag *Tag) error {
	query := "INSERT INTO tags(label, lowercase) VALUES ($1, $2) RETURNING id"
	args := []any{newTag.Label, newTag.Lowercase}
	scanFunc := func(row *sql.Row) error {
		return row.Scan(&newTag.ID)
	}

	return insertRecord(db, query, args, scanFunc)
}

func UpdateTag(db *sql.DB, id int32, tag *Tag) error {
	tag.Lowercase = strings.ToLower(tag.Label)

	query := "UPDATE tags SET label = $1, lowercase = $2 WHERE id = $3"
	args := []any{id, tag.Label, tag.Lowercase}

	return updateRecord(db, query, args)
}

func DeleteTagById(db *sql.DB, id int32) error {
	query := "DELETE FROM tags WHERE id = $1"
	args := []any{id}

	return deleteRecord(db, query, args)
}
