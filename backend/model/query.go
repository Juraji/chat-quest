package model

import (
	"database/sql"
	"errors"
)

// rowScanner interface that works with both sql.Row and sql.Rows
type rowScanner interface {
	Scan(dest ...any) error
}

// queryExecutor interface that works with both *sql.DB and *sql.Tx
type queryExecutor interface {
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	Exec(query string, args ...any) (sql.Result, error)
}

func queryForList[T any](
	q queryExecutor,
	query string,
	args []any,
	scanFunc func(scanner rowScanner, dest *T) error,
) ([]*T, error) {
	records := make([]*T, 0)
	var err error

	rows, err := q.Query(query, args...)
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
	q queryExecutor,
	query string,
	args []any,
	scanFunc func(scanner rowScanner, dest *T) error,
) (*T, error) {
	var dest T

	row := q.QueryRow(query, args...)
	err := scanFunc(row, &dest)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &dest, nil
}

func insertRecord(
	q queryExecutor,
	query string, args []any,
	scanFunc func(scanner rowScanner) error,
) error {
	row := q.QueryRow(query, args...)
	err := scanFunc(row)

	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}

	return err
}

func updateRecord(
	q queryExecutor,
	query string,
	args []any,
) error {
	res, err := q.Exec(query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func deleteRecord(
	q queryExecutor,
	query string,
	args []any,
) error {
	_, err := q.Exec(query, args...)
	return err
}

//goland:noinspection GoUnusedParameter
func noopScanFunc(scanner rowScanner) error {
	return nil
}
