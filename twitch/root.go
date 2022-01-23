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

	// Prepare the bot's net.Conn struct
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

	fmt.Printf("Connected to Twitch!\nBot: %s\nChannel: %s\n", t.Bot.Name, t.Bot.ChannelName)

	var item bot.Item
	var line string

	// Keep running as long as the error is not a fatal error.
	// It's likely that all errors encountered will just invoke a 'continue',
	// so that the for loop condition can determine whether we need to stop or not
	// TODO: I should make sure that using errors.As this much doesn't have any negative performance implications
	for !errors.As(err, &bot.FatalError{}) {
		line, err = proto.ReadLine()
		if err != nil {
			continue
		}
		item, err = newTwitchItem(line)
		if err != nil {
			t.Message(fmt.Sprintf("@%s - %s", item.Sender.Name, err.Error()))
			continue
		}
		err = t.Handler(item, setupDefaultActions())
		if err != nil {
			continue
		}
	}

	return err
}
