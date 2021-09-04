// Experimental file. Going to try and re-use as much code for database calls as much as possible.package bot

package bot

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
)

// this should only run in CreateBot, when the database file is not found in the config directory
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

// InsertIntoDB inserts some n values into a table of tableName
// creates a statement of form: "insert into foo(col1, col2) values(?, ?)" and inserts into DB
func (bot *Bot) InsertIntoDB(tableName string, columns, values []string) error {
	if len(columns) != len(values) { // columns and values must be the same length
		return errors.New("the columns and values arrays must be of the same size")
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
	_, err := bot.DB.Exec(stmt)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") { // duplicate command is the most expected error to occur
			return fmt.Errorf(fmt.Sprintf("the command '%s' already exists", values[0]))
		}
		return err // if the exact error isn't known, return the original error
	}
	return nil
}

// DeleteFromDB takes in the necessary parameters to delete a row from the DB
// currently only supports one column and value due to the structure of the current tables in the DB.
func (bot *Bot) DeleteFromDB(tableName string, column, value string) error {
	value = fmt.Sprintf("'%s'", value) // sqlite require quotes around value
	stmt := fmt.Sprintf("delete from %s where %s = %s", tableName, column, value)
	_, err := bot.DB.Exec(stmt)
	if err != nil {
		return fmt.Errorf("Error deleting value %s from column %s due to error: %s", value, column, err)
	}
	return nil
}

// GetTopFromTable returns the top result using the given parameters for table, target column, and the count column used for the MAX call.
func (bot *Bot) GetTopFromTable(table, column, countColumn string) (string, int) {
	topCom := bot.DB.QueryRow(fmt.Sprintf("SELECT %s, MAX(%s) from %s;", column, countColumn, table))

	var name string
	var count int
	err := topCom.Scan(&name, &count)
	if err != nil {
		fmt.Printf("Error getting top result: %s\n", err)
	}

	return name, count
}
