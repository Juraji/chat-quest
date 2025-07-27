package model

import (
	"database/sql"
	"errors"
)

type RowScanner interface {
	Scan(dest ...any) error
}

type DbEntity struct {
	ID int64 `json:"id"`
}

func queryForList[T any](
	db *sql.DB,
	query string,
	args []any,
	scanFunc func(scanner RowScanner, dest *T) error,
) ([]*T, error) {
	records := make([]*T, 0)
	var err error

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	for rows.Next() {
		var dest T

		err = scanFunc(rows, &dest)
		if err != nil {
			return nil, err
		}

		records = append(records, &dest)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return records, nil
}

func queryForRecord[T any](
	db *sql.DB,
	query string,
	args []any,
	scanFunc func(scanner RowScanner, dest *T) error,
) (*T, error) {
	var dest T

	row := db.QueryRow(query, args...)
	err := scanFunc(row, &dest)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &dest, nil
}

func insertRecord(db *sql.DB, query string, args []any, scanFunc func(scanner RowScanner) error) error {
	row := db.QueryRow(query, args...)
	err := scanFunc(row)

	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}

	return err
}

func updateRecord(db *sql.DB, query string, args []any) error {
	_, err := db.Exec(query, args...)
	return err
}

func deleteRecord(db *sql.DB, query string, args []any) error {
	_, err := db.Exec(query, args...)
	return err
}
