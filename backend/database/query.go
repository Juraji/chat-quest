package database

import (
	"database/sql"
	"errors"
	"log"
)

// RowScanner interface that works with both sql.Row and sql.Rows
type RowScanner interface {
	Scan(dest ...any) error
}

// QueryExecutor interface that works with both *sql.DB and *sql.Tx
type QueryExecutor interface {
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	Exec(query string, args ...any) (sql.Result, error)
}

func QueryForList[T any](
	q QueryExecutor,
	query string,
	args []any,
	scanFunc func(scanner RowScanner, dest *T) error,
) ([]*T, error) {
	records := make([]*T, 0)

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

func QueryForRecord[T any](
	q QueryExecutor,
	query string,
	args []any,
	scanFunc func(scanner RowScanner, dest *T) error,
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

func InsertRecord(
	q QueryExecutor,
	query string, args []any,
	scanTo ...any,
) error {
	row := q.QueryRow(query, args...)

	if len(scanTo) != 0 {
		err := row.Scan(scanTo...)

		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}

		return err
	}

	return nil
}

func UpdateRecord(
	q QueryExecutor,
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
		return errors.New("no rows affected")
	}

	return nil
}

func DeleteRecord(
	q QueryExecutor,
	query string,
	args []any,
) error {
	_, err := q.Exec(query, args...)
	return err
}

func RollBackOnErr(tx *sql.Tx, err error) {
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			log.Printf("Rollback failed: %v", rbErr)
		}
	}
}
