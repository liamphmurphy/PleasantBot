package storage

import "database/sql"

// Database defines the interface required for PleasantBot's database operations
type Database interface {
	Init(loader string) error
	Close() error
	Query(query string) (*sql.Rows, error)
	Insert(tableName string, columns []string, values []string) error
	Delete(tableName string, keyColumn string, keyValue string) error
	ArbitraryExec(statement string) error // Run any arbitrary query not supported in the other interface methods
}