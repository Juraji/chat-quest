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

func tagScanner(scanner RowScanner, dest *Tag) error {
	return scanner.Scan(
		&dest.ID,
		&dest.Label,
		&dest.Lowercase,
	)
}

func AllTags(db *sql.DB) ([]*Tag, error) {
	query := "SELECT * FROM tags"
	return queryForList(db, query, nil, tagScanner)
}

func TagById(db *sql.DB, id int64) (*Tag, error) {
	query := "SELECT * FROM tags WHERE id = $1"
	args := []any{id}
	return queryForRecord(db, query, args, tagScanner)
}

func CreateTag(db *sql.DB, newTag *Tag) error {
	query := "INSERT INTO tags(label, lowercase) VALUES ($1, $2) RETURNING id"
	args := []any{newTag.Label, newTag.Lowercase}
	scanFunc := func(scanner RowScanner) error {
		return scanner.Scan(&newTag.ID)
	}

	return insertRecord(db, query, args, scanFunc)
}

func UpdateTag(db *sql.DB, id int64, tag *Tag) error {
	tag.Lowercase = strings.ToLower(tag.Label)

	query := "UPDATE tags SET label = $1, lowercase = $2 WHERE id = $3"
	args := []any{id, tag.Label, tag.Lowercase}

	return updateRecord(db, query, args)
}

func DeleteTagById(db *sql.DB, id int64) error {
	query := "DELETE FROM tags WHERE id = $1"
	args := []any{id}

	return deleteRecord(db, query, args)
}
