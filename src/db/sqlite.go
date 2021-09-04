package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
)

type Sqlite struct {
	db *sql.DB
}

// this should only run in a sqlite Init call, when the database file is not found in the config directory
func prepareDatabase(db *sql.DB) {
	stmt := `
	CREATE TABLE IF NOT EXISTS commands (id INTEGER PRIMARY KEY, commandname TEXT UNIQUE, commandresponse TEXT, perm TEXT, count INTEGER);
	CREATE TABLE IF NOT EXISTS badwords (id INTEGER PRIMARY KEY, phrase TEXT, severity INTEGER);
	CREATE TABLE IF NOT EXISTS quotes (id INTEGER PRIMARY KEY, quote TEXT, timestamp TEXT, submitter TEXT);
	CREATE TABLE IF NOT EXISTS ban_history (user TEXT, reason TEXT, timestamp TEXT);
	CREATE TABLE IF NOT EXISTS chatters (username TEXT PRIMARY KEY, count INT);
	CREATE TABLE IF NOT EXISTS timers (timername TEXT UNIQUE, message TEXT, minutes INTEGER);
	`
	_, err := db.Exec(stmt)
	if err != nil {
		log.Fatalf("error running create table statements: %s", err)
	}
}

func (sq *Sqlite) Init(dir string) error {
	// prepare Sqlite 3 database
	dbFile := dir + "/pleasantbot.db"
	if _, err := os.Stat(dbFile); os.IsNotExist(err) { // make database file if it doesn't exist
		os.Create(dbFile)
	}
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return err
	}

	sq.db  = db
	prepareDatabase(db) // creates and prepares the bot's database

	defer db.Close()
	return nil
}

// Query takes in a query and returns the resulting rows
func (sq *Sqlite) Query(query string) (*sql.Rows, error)  {
	rows, err := sq.db.Query(query)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (sq *Sqlite) Insert(tableName string, columns []string, values[]string) error {
	if len(columns) != len(values) { // columns and values must be the same length
		return fmt.Errorf("the columns and values arrays must be of the same size")
	}

	// sqlite query needs quotes around string values, there's probably a better way to do this
	for i := range values {
		values[i] = "'" + values[i] + "'"
	}

	// format the columns and values to work with the SQLite insert statement
	columnsFormatted := strings.Join(columns, ", ")
	valuesFormatted := strings.Join(values, ", ")

	// insert formatted data into DB
	stmt := fmt.Sprintf("insert into %s(%s) values(%s)", tableName, columnsFormatted, valuesFormatted)
	_, err := sq.db.Exec(stmt)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") { // duplicate entry is the most expected error to occur
			return fmt.Errorf(fmt.Sprintf("the command '%s' already exists", values[0]))
		}
		return err // if the exact error isn't known, return the original error
	}
	return nil
}

func (sq *Sqlite) Delete(tableName string, keyColumn string, keyValue string) error {
	value := fmt.Sprintf("'%s'", keyValue) // sqlite requires quotes around value
	stmt := fmt.Sprintf("delete from %s where %s = %s", tableName, keyColumn, value)
	err := sq.ArbitraryExec(stmt)
	if err != nil {
		return fmt.Errorf("error deleting value %s from column %s due to error: %s", value, keyColumn, err)
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