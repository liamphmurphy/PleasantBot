// This file handles some of the core functionality of the bot. Such as creation of the bot and connecting to Twitch.

package bot

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"regexp"

	_ "github.com/mattn/go-sqlite3" // docs have a blank import so I'm using that
	"github.com/spf13/viper"
)

// Bot struct contains the necessary data to run an instance of a bot
type Bot struct {
	ChannelName     string
	ServerName      string
	OAuth           string
	Name            string
	Conn            net.Conn
	Commands        []Command
	BadWords        []BadWord
	Quotes          []string
	PermittedUsers  []string // list of users that can post links
	DB              *sql.DB
	DBPath          string
	PurgeForLinks   bool
	PurgeForLongMsg bool
	Perms           []string // holds a list of users that can post a link
}

// writeConfig is run whenever the config.toml file doesn't exist, usually after a fresh download of the bot.
func writeConfig(path string) {
	path = path + "/config.toml"
	// prepare default values, will be used when viper writes the new config file
	viper.SetDefault("ChannelName", "<enter channel name to moderate here>")
	viper.SetDefault("ServerName", "irc.twitch.tv:6667")
	viper.SetDefault("BotName", "<enter bot username here>")
	viper.SetDefault("BotOAuth", "<bot oauth>")
	viper.SetDefault("PurgeForLinks", true)
	viper.SetDefault("PurgeForLongMsg", true)

	viper.WriteConfigAs(path)
	fmt.Println(fmt.Sprintf("Config file did not exist, so it has been made. Please go to %s and edit the settings.", path))
}

// CreateBot creates an instance of a bot
func CreateBot() *Bot {
	configDir, _ := os.UserConfigDir()        // follows the standard config dir used by the OS
	pleasantDir := configDir + "/pleasantbot" // config.toml will go here
	if _, err := os.Stat(pleasantDir); os.IsNotExist(err) {
		os.Mkdir(pleasantDir, 0755)
	}

	// prepare toml config file
	viper.SetConfigName("config")
	viper.AddConfigPath(pleasantDir)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok { // if config not found, write a default one
			writeConfig(pleasantDir)
			os.Exit(1) // write a better way to exit this gracefully
		} else {
			panic(fmt.Errorf("error reading in config file: %s", err))
		}
	}

	var bot Bot
	var db *sql.DB
	// prepare Sqlite 3 database
	dbFile := pleasantDir + "/pleasantbot.db"
	if _, err := os.Stat(dbFile); os.IsNotExist(err) { // make database file if it doesn't exist
		os.Create(dbFile)
	}
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatalf("error trying to open the sqlite3 db file: %s\n", err)
	}

	defer db.Close()

	prepareDatabase(db) // creates and prepares the bot's database
	bot.DB = db
	bot.DBPath = dbFile

	// load data
	err = bot.LoadCommands()
	if err != nil {
		log.Fatalf("error loading commands from the database: %s\n", err)
	}

	err = bot.LoadQuotes()
	if err != nil {
		log.Fatalf("error loading quotes from the database: %s\n", err)
	}

	err = bot.LoadBadWords()
	if err != nil {
		log.Fatalf("error loading bannable words from the database: %s\n", err)
	}

	// assign bot values provided by the config file
	bot.ChannelName = viper.GetString("ChannelName")
	bot.ServerName = viper.GetString("ServerName")
	bot.OAuth = viper.GetString("BotOAuth")
	bot.Name = viper.GetString("BotName")
	bot.PurgeForLinks = viper.GetBool("PurgeForLinks")
	bot.PurgeForLongMsg = viper.GetBool("PurgeForLongMsg")

	return &bot
}

// Connect establishes a connection to the Twitch IRC server
func (bot *Bot) Connect() error {
	var err error
	if bot.Conn != nil {
		return fmt.Errorf("ERROR: the bot has already established an IRC connection")
	}
	bot.Conn, err = net.Dial("tcp", bot.ServerName)
	return err
}

// ChannelConnect writes the necessary scopes to Twitch
func (bot *Bot) ChannelConnect() {
	// Pass info to HTTP request
	fmt.Fprintf(bot.Conn, "PASS %s\r\n", bot.OAuth)
	fmt.Fprintf(bot.Conn, "NICK %s\r\n", bot.Name)
	fmt.Fprintf(bot.Conn, "JOIN #%s\r\n", bot.ChannelName)

	// Twitch specific information, like badges, mod status etc.
	fmt.Fprintf(bot.Conn, "CAP REQ :twitch.tv/membership\r\n")
	fmt.Fprintf(bot.Conn, "CAP REQ :twitch.tv/tags\r\n")
	fmt.Fprintf(bot.Conn, "CAP REQ :twitch.tv/commands\r\n")

	defer bot.Conn.Close()
}

// WriteToTwitch when given a string sends a properly formatted message to Twitch, easy replacement for using fmt.Fprintf
func (bot *Bot) WriteToTwitch(msg string) {
	fmt.Fprintf(bot.Conn, "%s\r\n", msg)
}

// SendMessage prepares and sends a string to the channel's Twitch chat
func (bot *Bot) SendMessage(msg string) {
	bot.WriteToTwitch(fmt.Sprintf("PRIVMSG #%s :%s", bot.ChannelName, msg))
}

// Itob converts an integer (0 or 1) to a corresponding boolean. Mainly used for command moderator perms
func Itob(i int) bool {
	if i == 1 {
		return true
	}

	return false
}

// FilterForSpam parses user message for some config options such as PurgeForLinks to see if message could be spam
func (bot *Bot) FilterForSpam(message User) {
	if bot.PurgeForLinks { // if enabled, check if message contains a link
		// regex gathered from top answer here: https://stackoverflow.com/questions/3809401/what-is-a-good-regular-expression-to-match-a-url
		urlRegex, _ := regexp.Compile(`[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)`)
		if urlRegex.MatchString(message.Content) {
			bot.SendMessage("Please don't type links :)")
		}
	}
}
