package storage

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type Sqlite struct {
	db *sql.DB
}

type InitFunc func(path string, sq *Sqlite, prepareFunc DatabasePrepareFunc) error

type DatabasePrepareFunc func(db *sql.DB) error

func Init(path string, sq *Sqlite, prepareFunc DatabasePrepareFunc) error {
	// prepare Sqlite 3 database
	if _, err := os.Stat(path); os.IsNotExist(err) { // make database file if it doesn't exist
		os.Create(path)
	}
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return err
	}

	sq.db = db
	err = prepareFunc(db) // in general this will prepare the schema of the db for a specific service

	return err
}

// Query takes in a query and returns the resulting rows
func (sq *Sqlite) Query(query string) (*sql.Rows, error) {
	rows, err := sq.db.Query(query)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (sq *Sqlite) Insert(tableName string, columns []string, values []string) error {
	if len(columns) != len(values) { // columns and values must be the same length
		return fmt.Errorf("the columns and values arrays must be of the same size")
	}

	// sqlite query needs quotes around string values, there's probably a better way to do this
	for i := range values {
		values[i] = fmt.Sprintf("'%s'", values[i])
	}

	// format the columns and values to work with the SQLite insert statement
	columnsFormatted := strings.Join(columns, ", ")
	valuesFormatted := strings.Join(values, ", ")

	// insert formatted data into DB
	stmt := fmt.Sprintf("insert into %s(%s) values(%s)", tableName, columnsFormatted, valuesFormatted)
	_, err := sq.db.Exec(stmt)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") { // duplicate entry is the most expected error to occur
			return fmt.Errorf(fmt.Sprintf("the item '%s' already exists", values[0]))
		}
		return err // if the exact error isn't known, return the original error
	}
	return nil
}

func (sq *Sqlite) Delete(tableName string, keyColumn string, keyValue string) error {
	stmt := fmt.Sprintf("delete from %s where %s = '%s'", tableName, keyColumn, keyValue)
	err := sq.ArbitraryExec(stmt)
	if err != nil {
		return fmt.Errorf("error deleting value %s from column %s due to error: %s", keyValue, keyColumn, err)
	}
	return nil
}

func (sq *Sqlite) ArbitraryExec(statement string) error {
	_, err := sq.db.Exec(statement)
	if err != nil {
		return err
	}

	return nil
}

func (sq *Sqlite) Close() error {
	return sq.db.Close()
}
