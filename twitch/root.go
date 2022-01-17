// This file acts as the launch point for a Twitch bot

package twitch

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"net/textproto"

	"github.com/liamphmurphy/pleasantbot/bot"
	"github.com/liamphmurphy/pleasantbot/storage"
)

var (
	twitchServer     = "irc.chat.twitch.tv:6697"
	configFileName   = "twitch"
	configFileType   = "toml"
	databaseFileName = "pleasantbot.db"
	commandCols      = []string{"commandname", "commandresponse", "perm", "count"}
	quoteCols        = []string{"quote", "timestamp", "submitter"}
	timerCols        = []string{"timername", "message", "minutes", "enabled"}
)

type Twitch struct {
	Bot *bot.Bot
}

// this should only run in a sqlite Init call, when the database file is not found in the config directory
func prepareDatabase(db *sql.DB) error {
	stmt := `
	CREATE TABLE IF NOT EXISTS commands (id INTEGER PRIMARY KEY, commandname TEXT UNIQUE, commandresponse TEXT, perm TEXT, count INTEGER);
	CREATE TABLE IF NOT EXISTS badwords (id INTEGER PRIMARY KEY, phrase TEXT, severity INTEGER);
	CREATE TABLE IF NOT EXISTS quotes (id INTEGER PRIMARY KEY, quote TEXT, timestamp TEXT, submitter TEXT);
	CREATE TABLE IF NOT EXISTS ban_history (user TEXT, reason TEXT, timestamp TEXT);
	CREATE TABLE IF NOT EXISTS chatters (username TEXT PRIMARY KEY, count INT);
	CREATE TABLE IF NOT EXISTS timers (timername TEXT UNIQUE, message TEXT, minutes INTEGER, enabled INTEGER);
	`
	_, err := db.Exec(stmt)

	return err
}

// HandleMessage will contain the root logic for handling any kind of message from twitch; whether from the IRC server itself,
// or messages from the Twitch chat
func (t *Twitch) HandleMessage(msg string) error {
	return nil
}

// Run defines the main entry point for a Twitch bot
func (t *Twitch) Run() error {
	configDir, err := bot.GetConfigDirectory()
	if err != nil {
		return err
	}

	v, err := bot.CreateViperConfig(configDir, configFileName, configFileType, twitchServer)
	if err != nil {
		return err
	}

	database := bot.Database{Path: fmt.Sprintf("%s/%s", configDir, databaseFileName), CommandColumns: commandCols,
		QuoteColumns: quoteCols, TimerColumns: timerCols}

	t.Bot, err = bot.CreateBot(v, database, storage.Init, bot.LoadBot, prepareDatabase)
	if err != nil {
		return bot.FatalError{Err: err}
	}

	// Prepate the bot's net.Conn struct
	err = t.Bot.Connect()
	if err != nil {
		return bot.FatalError{Err: err}
	}

	// send initial messages to Twitch as specified in https://dev.twitch.tv/docs/irc/guide#connecting-to-twitch-irc
	err = t.initialConn(t.craftInitialConnMessages())
	if err != nil {
		return bot.FatalError{Err: err}
	}

	defer t.Bot.Conn.Close()

	reader := bufio.NewReader(t.Bot.Conn)
	proto := textproto.NewReader(reader)

	var line string
	// keep running as long as the error is not a fatal error
	for !errors.As(err, &bot.FatalError{}) {
		line, err = proto.ReadLine()
		if err != nil {
			continue
		}

		err = t.HandleMessage(line)
	}

	return nil
}
