package database

import (
	"database/sql"
	"github.com/pkg/errors"
)

// RowScanner interface that works with both sql.Row and sql.Rows
type RowScanner interface {
	Scan(dest ...any) error
}

// queryExecutor interface that works with both *sql.DB and *sql.Tx
type queryExecutor interface {
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	Exec(query string, args ...any) (sql.Result, error)
}

type TxContext struct {
	tx *sql.Tx
}

func QueryForList[T any](query string, args []any, scanFunc func(scanner RowScanner, dest *T) error) ([]T, error) {
	return queryForList(GetDB(), query, args, scanFunc)
}

func QueryForRecord[T any](query string, args []any, scanFunc func(scanner RowScanner, dest *T) error) (*T, error) {
	return queryForRecord(GetDB(), query, args, scanFunc)
}

func InsertRecord(query string, args []any, scanTo ...any) error {
	return insertRecord(GetDB(), query, args, scanTo)
}

func UpdateRecord(query string, args []any) error {
	return updateRecord(GetDB(), query, args)
}

func DeleteRecord(query string, args []any) error {
	return deleteRecord(GetDB(), query, args)
}

func Transactional(action func(ctx *TxContext) error) error {
	tx, err := GetDB().Begin()
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}
	defer func(err error) {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				panic(rbErr)
			}
		}
	}(err)

	err = action(&TxContext{tx: tx})
	if err != nil {
		return errors.Wrap(err, "failed to execute transaction")
	}

	err = tx.Commit()
	if err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	return nil
}

func (tx *TxContext) InsertRecord(query string, args []any, scanTo ...any) error {
	return insertRecord(tx.tx, query, args, scanTo)
}
func (tx *TxContext) UpdateRecord(query string, args []any) error {
	return updateRecord(tx.tx, query, args)
}
func (tx *TxContext) DeleteRecord(query string, args []any) error {
	return deleteRecord(tx.tx, query, args)
}

func queryForList[T any](
	q queryExecutor,
	query string,
	args []any,
	scanFunc func(scanner RowScanner, dest *T) error,
) ([]T, error) {
	records := make([]T, 0)

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

		records = append(records, dest)
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

func insertRecord(
	q queryExecutor,
	query string,
	args []any,
	scanTo []any,
) error {
	if len(scanTo) != 0 {
		row := q.QueryRow(query, args...)
		err := row.Scan(scanTo...)
		return err
	} else {
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
		return errors.New("no rows affected")
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

func StringScanner(scanner RowScanner, dest *string) error {
	return scanner.Scan(dest)
}

func BoolScanner(scanner RowScanner, dest *bool) error {
	return scanner.Scan(dest)
}
