// This file handles some of the core functionality of the bot. Such as creation of the bot and connecting to Twitch.

package bot

import (
	"crypto/tls"
	"fmt"
	"github.com/murnux/pleasantbot/db"
	"log"
	"net"
	"os"
	"regexp"
	"time"

	"github.com/spf13/viper"
)

// Bot struct contains the necessary data to run an instance of a bot
type Bot struct {
	Name             string
	ChannelName      string
	ServerName       string
	oauth            string       `json:"-"`
	Config           *viper.Viper `json:"-"`
	Authenticated    bool         // used to tell the front-end GUI whether the bot has been authenticated yet
	Conn             net.Conn     `json:"-"`
	DB               db.Database      `json:"-"`
	DBPath           string       `json:"-"`
	PurgeForLinks    bool
	PurgeForLongMsg  bool
	LongMsgAmount    int
	EnableServer     bool
	PostLinkPerm     uint8
	Perms            []string                 `json:"-"` // holds a list of users that can post a link
	Commands         map[string]*CommandValue `json:"-"`
	BadWords         []BadWord                `json:"-"`
	Quotes           map[int]*QuoteValues     `json:"-"`
	PermittedUsers   map[string]struct{}      // list of users that can post links
	CommandDBColumns []string                 `json:"-"` // used for InsertIntoDB calls
	QuoteDBColumns   []string                 `json:"-"`
}

// writeConfig is run whenever the config.toml file doesn't exist, usually after a fresh download of the bot.
func writeConfig(path string, configObject *viper.Viper) {
	path = path + "/config.toml"
	// prepare default values, will be used when viper writes the new config file
	configObject.SetDefault("ChannelName", "<enter channel name to moderate here>")
	configObject.SetDefault("ServerName", "irc.twitch.tv:6667")
	configObject.SetDefault("BotName", "<enter bot username here>")
	configObject.SetDefault("BotOAuth", "<bot oauth>")
	configObject.SetDefault("PurgeForLinks", true)
	configObject.SetDefault("PurgeForLongMsg", true)
	configObject.SetDefault("LongMsgAmount", 400)
	configObject.SetDefault("EnableServer", true)
	configObject.SetDefault("PostLinkPerm", uint(1)) // Minimum permission needed for non-purging links, in this case subscriber

	configObject.WriteConfigAs(path)
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
	viperConfig := viper.New()
	viperConfig.SetConfigName("config")
	viperConfig.AddConfigPath(pleasantDir)
	if err := viperConfig.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok { // if config not found, write a default one
			writeConfig(pleasantDir, viperConfig)
			os.Exit(1) // write a better way to exit this gracefully
		} else {
			panic(fmt.Errorf("error reading in config file: %s", err))
		}
	}

	var bot Bot
	var db db.Sqlite
	err := db.Init(pleasantDir)
	if err != nil {
		return nil
	}

	bot.Config = viperConfig

	bot.DB = &db

	// load data
	bot.Commands = make(map[string]*CommandValue)
	err = bot.LoadCommands()
	if err != nil {
		log.Fatalf("error loading commands from the database: %s\n", err)
	}

	bot.Quotes = make(map[int]*QuoteValues)
	err = bot.LoadQuotes()
	if err != nil {
		log.Fatalf("error loading quotes from the database: %s\n", err)
	}

	err = bot.LoadBadWords()
	if err != nil {
		log.Fatalf("error loading bannable words from the database: %s\n", err)
	}

	// assign bot values provided by the config file
	bot.ChannelName = bot.Config.GetString("ChannelName")
	bot.ServerName = bot.Config.GetString("ServerName")
	bot.oauth = "oauth:" + bot.Config.GetString("BotOAuth")
	bot.Name = bot.Config.GetString("BotName")
	bot.PurgeForLinks = bot.Config.GetBool("PurgeForLinks")
	bot.PurgeForLongMsg = bot.Config.GetBool("PurgeForLongMsg")
	bot.LongMsgAmount = bot.Config.GetInt("LongMsgAmount")
	bot.EnableServer = bot.Config.GetBool("EnableServer")
	bot.CommandDBColumns = []string{"commandname", "commandresponse", "perm", "count"}
	bot.QuoteDBColumns = []string{"quote", "timestamp", "submitter"}

	// determine using the oauth string whether the user has logged in yet
	if bot.oauth != "oauth:" && bot.oauth != "" {
		bot.Authenticated = true
	} else {
		bot.Authenticated = false
	}

	return &bot
}

// Connect establishes a connection to the Twitch IRC server
func (bot *Bot) Connect() error {
	var err error
	conf := &tls.Config{
		//InsecureSkipVerify: true,
	}
	if bot.Conn != nil {
		return fmt.Errorf("ERROR: the bot has already established an IRC connection")
	}
	bot.Conn, err = tls.Dial("tcp", bot.ServerName, conf)
	///bot.Conn, err = net.Dial("tcp", bot.ServerName)
	return err
}

// ChannelConnect writes the necessary scopes to Twitch
func (bot *Bot) ChannelConnect() {

	// Pass info to HTTP request
	fmt.Fprintf(bot.Conn, "PASS %s\r\n", bot.oauth)
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
	return i == 1
}

// FilterForSpam parses user message for some config options such as PurgeForLinks to see if message could be spam
func (bot *Bot) FilterForSpam(message User) {
	if bot.PurgeForLinks { // if enabled, check if message contains a link
		// regex obtained from top answer here: https://stackoverflow.com/questions/3809401/what-is-a-good-regular-expression-to-match-a-url
		urlRegex, _ := regexp.Compile(`[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)`)
		if urlRegex.MatchString(message.Content) {
			if _, permitted := bot.PermittedUsers[message.Name]; permitted || message.Perm >= bot.PostLinkPerm { // let user post if they are in the permit list or is a moderator / broadcaster
				delete(bot.PermittedUsers, message.Name) // found, do nothing except delete from the map
			} else {
				bot.purgeUser(message.Name) // not permitted, so purge user
				bot.SendMessage(fmt.Sprintf("%s, you do not have permissions to post links.", message.Name))
			}
		}
	}

	if bot.PurgeForLongMsg { // if enabled, check if a message is very long
		if len(message.Content) >= bot.LongMsgAmount {
			bot.purgeUser(message.Name)
		}
	}
}

// AddPermittedUser adds a user to the PermittedUsers slice, allowing them to post a link without being purged
func (bot *Bot) AddPermittedUser(username string) {
	bot.PermittedUsers[username] = struct{}{}
}

// purges a user by sending a timeout of 1 second
func (bot *Bot) purgeUser(username string) {
	bot.SendMessage(fmt.Sprintf("/timeout %s 1", username))
}

func (bot *Bot) banUser(username string, reason string) {
	bot.DB.Insert("ban_history", []string{"user", "reason", "timestamp"}, []string{username, reason, time.Now().Format("2006-01-02 15:04:05")}) // insert into ban_history table
	bot.SendMessage(fmt.Sprintf("/ban %s", username))
}

// GetOAuth returns the bot's oauth token
func (bot *Bot) GetOAuth() string {
	return bot.oauth
}

// SetOAuth creates or replaces the bot Oauth token.
func (bot *Bot) SetOAuth(token string) {
	newToken := fmt.Sprintf("oauth:%s", token)
	bot.Config.Set("botoauth", newToken)
	bot.oauth = newToken // set runtime oauth
	bot.Authenticated = true
	bot.Config.WriteConfig() // update config with new oauth token
}
