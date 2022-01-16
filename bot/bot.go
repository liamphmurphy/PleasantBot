// This file handles some of the core functionality of the bot. Such as creation of the bot and connecting to Twitch.

package bot

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"regexp"
	"strings"
	"time"

	"github.com/liamphmurphy/pleasantbot/storage"

	"github.com/spf13/viper"
)

var (
	urlRegex, _ = regexp.Compile(`[-a-zA-Z-1-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)`)
)

// Bot struct contains the necessary data to run an instance of a bot
type Bot struct {
	Name            string
	ChannelName     string
	ServerName      string
	oauth           string       `json:"-"`
	Config          *viper.Viper `json:"-"`
	Authenticated   bool
	Conn            net.Conn `json:"-"`
	Storage         *Database
	PurgeForLinks   bool
	PurgeForLongMsg bool
	LongMsgAmount   int
	EnableServer    bool
	PostLinkPerm    uint8
	Perms           []string                 `json:"-"` // holds a list of users that can post a link
	DefaultCommands []DefaultCommand         `json:"-"`
	Commands        map[string]*CommandValue `json:"-"`
	BadWords        []BadWord                `json:"-"`
	Quotes          map[int]*QuoteValues     `json:"-"`
	Timers          map[string]*TimedValue   `json:"-"`
	PermittedUsers  map[string]struct{}      // list of users that can post links

}

// defines some DB options, specifically the columns for the various tables
type Database struct {
	DB             storage.Sqlite `json:"-"`
	Path           string         `json:"-"` // path to the database file
	CommandColumns []string
	QuoteColumns   []string
	TimerColumns   []string
}

type BotLoaderFunc func(bot *Bot) error

var errNoViper = errors.New("the bot has a nil viper config struct, this needs to be made first")

// CreateBot creates an instance of a bot. The function makes no assumptions on what the viper config, Database, and
// loader funcs for the database and bot data loading should be. It just makes the assumption that these loaders MUST exist, so
// the caller may pass in custom ones based on the needs of the service (e.g. twitch, youtube, etc.) or it may use the default ones
// listed in this package.
func CreateBot(viper *viper.Viper, storage Database, dbInit storage.InitFunc, loader BotLoaderFunc) (*Bot, error) {
	var bot Bot
	if viper == nil {
		return &bot, FatalError{Err: errNoViper}
	}
	bot.Config = viper

	err := dbInit(storage.Path, &storage.DB)
	if err != nil {
		return &Bot{}, err
	}
	bot.Storage = &storage

	// load bot data using the passed in loader
	err = loader(&bot)

	return &bot, err
}

// LoadBot populates many of the misc. struct field values
func LoadBot(bot *Bot) error {
	if bot.Config == nil {
		return FatalError{Err: errNoViper}
	}

	// assign bot values provided by the config file
	bot.ChannelName = bot.Config.GetString("ChannelName")
	bot.ServerName = bot.Config.GetString("ServerName")
	bot.oauth = bot.Config.GetString("BotOAuth")
	bot.Name = bot.Config.GetString("BotName")
	bot.PurgeForLinks = bot.Config.GetBool("PurgeForLinks")
	bot.PurgeForLongMsg = bot.Config.GetBool("PurgeForLongMsg")
	bot.LongMsgAmount = bot.Config.GetInt("LongMsgAmount")
	bot.EnableServer = bot.Config.GetBool("EnableServer")

	var err error
	// load data
	bot.Commands = make(map[string]*CommandValue)
	err = bot.LoadCommands()
	if err != nil {
		return err
	}

	bot.Quotes = make(map[int]*QuoteValues)
	err = bot.LoadQuotes()
	if err != nil {
		return err
	}

	err = bot.LoadBadWords()
	if err != nil {
		return err
	}

	bot.Timers = make(map[string]*TimedValue)
	err = bot.LoadTimers()
	if err != nil {
		return err
	}

	return err
}

// Connect establishes a connection to the Twitch IRC server
func (bot *Bot) Connect() error {
	var err error
	if bot.Conn != nil {
		return fmt.Errorf("ERROR: the bot has already established a connection")
	}
	bot.Conn, err = tls.Dial("tcp", bot.ServerName, &tls.Config{})
	return err
}

// WriteToConn when given a string sends a properly formatted message, easy replacement for using fmt.Fprintf
func (bot *Bot) WriteToConn(msg string) error {
	_, err := fmt.Fprintf(bot.Conn, "%s\r\n", msg)
	return err
}

// SendTwitchMessage prepares and sends a string to the channel's Twitch chat
func (bot *Bot) SendTwitchMessage(msg string) {
	bot.WriteToConn(fmt.Sprintf("PRIVMSG #%s :%s", bot.ChannelName, msg))
}

// Itob converts an integer (0 or 1) to a corresponding boolean. Mainly used for command moderator perms
func Itob(i int) bool {
	return i == 1
}

// FilterForSpam parses user message for some config options such as PurgeForLinks to see if message could be spam
func (bot *Bot) FilterForSpam(message User) {
	if bot.PurgeForLinks { // if enabled, check if message contains a link
		// regex obtained from top answer here: https://stackoverflow.com/questions/3809401/what-is-a-good-regular-expression-to-match-a-url
		if urlRegex.MatchString(message.Content) {
			if _, permitted := bot.PermittedUsers[message.Name]; permitted || message.Perm >= bot.PostLinkPerm { // let user post if they are in the permit list or is a moderator / broadcaster
				delete(bot.PermittedUsers, message.Name) // found, do nothing except delete from the map
			} else {
				bot.purgeUser(message.Name) // not permitted, so purge user
				bot.SendTwitchMessage(fmt.Sprintf("%s, you do not have permissions to post links.", message.Name))
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
	bot.SendTwitchMessage(fmt.Sprintf("/timeout %s 1", username))
}

func (bot *Bot) banUser(username string, reason string) {
	bot.Storage.DB.Insert("ban_history", []string{"user", "reason", "timestamp"}, []string{username, reason, time.Now().Format("2006-01-02 15:04:05")}) // insert into ban_history table
	bot.SendTwitchMessage(fmt.Sprintf("/ban %s", username))
}

// GetOAuth returns the bot's oauth token
func (bot *Bot) GetOAuth() string {
	if !strings.Contains(bot.oauth, "oauth") {
		return fmt.Sprintf("oauth:%s", bot.oauth)
	}
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

// HandlePing will send a PONG as a response to a PING. Returns true if a PONG had to be sent
// pingIndicator is the string to check for
func (bot *Bot) HandlePing(message, pingIndicator, response string) (bool, error) {
	contained := strings.Contains(message, pingIndicator)
	if contained {
		err := bot.WriteToConn(response)
		if err != nil {
			return false, err
		}
	}
	return contained, nil
}
