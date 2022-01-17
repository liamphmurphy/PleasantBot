// This file handles some of the core functionality of the bot. Such as creation of the bot and connecting to Twitch.

package bot

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"regexp"
	"strings"

	"github.com/liamphmurphy/pleasantbot/storage"

	"github.com/spf13/viper"
)

var (
	// regex obtained from top answer here: https://stackoverflow.com/questions/3809401/what-is-a-good-regular-expression-to-match-a-url
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

type BotLoaderFunc func(bot *Bot) error

// Messenger is used in those times where a function in the bot package HAS to send a message within its definition.
// This is avoided whenever possible, since the caller (service) should be determining what to do. However there are rare
// cases, such as with RunTimers, when a interface for sending messages is needed.
type Messenger interface {
	Message(msg string) error
}

var errNoViper = errors.New("the bot has a nil viper config struct, this needs to be made first")

// CreateBot creates an instance of a bot. The function makes no assumptions on what the viper config, Database, and
// loader funcs for the database and bot data loading should be. It just makes the assumption that these loaders MUST exist, so
// the caller may pass in custom ones based on the needs of the service (e.g. twitch, youtube, etc.) or it may use the default ones
// listed in this package.
func CreateBot(viper *viper.Viper, storage Database, dbInit storage.InitFunc, loader BotLoaderFunc, prepareFunc storage.DatabasePrepareFunc) (*Bot, error) {
	var bot Bot
	if viper == nil {
		return &bot, FatalError{Err: errNoViper}
	}
	bot.Config = viper

	err := dbInit(storage.Path, &storage.DB, prepareFunc)
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

// Itob converts an integer (0 or 1) to a corresponding boolean. Mainly used for command moderator perms
func Itob(i int) bool {
	return i == 1
}

// DetectURl uses urlRegex to determine if the passed in message is a URL. The caller than perform whatever
// action is desired with this information
func (bot *Bot) DetectURL(message string) bool {
	return urlRegex.MatchString(message)
}

// AddPermittedUser adds a user to the PermittedUsers slice, allowing them to post a link without being purged
func (bot *Bot) AddPermittedUser(username string) {
	bot.PermittedUsers[username] = struct{}{}
}

// DeletePermittedUser deletes a user from the permittedusers map if they exist in it
func (bot *Bot) DeletePermittedUser(username string) bool {
	var exists bool
	if _, exists = bot.PermittedUsers[username]; exists {
		delete(bot.PermittedUsers, username)
	}
	return exists
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
func (bot *Bot) HandlePing(message, pingIndicator, response string, messenger Messenger) (bool, error) {
	contained := strings.Contains(message, pingIndicator)
	if contained {
		err := messenger.Message(response)
		if err != nil {
			return false, err
		}
	}
	return contained, nil
}
