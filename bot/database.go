package bot

import (
	"github.com/liamphmurphy/pleasantbot/storage"
)

// defines some DB options, specifically the columns for the various tables
type Database struct {
	DB             storage.Sqlite `json:"-"`
	Path           string         `json:"-"` // path to the database file
	CommandColumns []string
	QuoteColumns   []string
	TimerColumns   []string
}
