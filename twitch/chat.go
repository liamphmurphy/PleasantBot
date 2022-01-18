package twitch

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/liamphmurphy/pleasantbot/bot"
)

var (
	typeRegex             = regexp.MustCompile(`^(\![\w]*)$`)                            // regexp for new item of form !itemcommand, such as "!quote" (note the absence of any values / content)
	commandNoContentRegex = regexp.MustCompile(`^(\![\w]*)\s(.)*\s(\![.\w]*)$`)          // regexp for request of form '!com del !somecommand'
	fullCommandRegex      = regexp.MustCompile(`^(\![\w]*)\s(.)*\s(\![.\w]*)\s([.\w]*)`) // regexp for request of form '!com add !somecommand this is a test command'
	errComParse           = errors.New("command invocation failed")
)

func (t *Twitch) Message(msg string) error {
	return t.Bot.WriteToConn(fmt.Sprintf("PRIVMSG #%s :%s", t.Bot.ChannelName, msg))
}

// Handler will contain the root logic for handling any kind of message from twitch; whether from the IRC server itself,
// or messages from the Twitch chat
func (t *Twitch) Handler(item bot.Item, defaultActions []ActionTaker) error {
	var at ActionTaker

	// see if the message will prompt a default action
	for _, v := range defaultActions {
		if v.Condition(item, t.Bot) {
			at = v
		}
	}

	if at == nil {
		at = &NoOpAction{}
	}

	// perform actions
	return at.Action(item, t.Bot, t)
}

// purges a user by sending a timeout of 1 second
func (t *Twitch) purgeUser(username string) {
	t.Message(fmt.Sprintf("/timeout %s 1", username))
}

func (t *Twitch) banUser(username string, reason string) {
	t.Bot.Storage.DB.Insert("ban_history", []string{"user", "reason", "timestamp"}, []string{username, reason, time.Now().Format("2006-01-02 15:04:05")}) // insert into ban_history table
	t.Message(fmt.Sprintf("/ban %s", username))
}

// newTwitchItem is a parser for the raw return from the bot's net.Conn.
func newTwitchItem(response string) (bot.Item, error) {
	// Only time a user sent a message is when PRIVMSG exists
	// TODO: more performant way of doing this check?
	if !strings.Contains(response, "PRIVMSG") {
		return bot.Item{IsServerInfo: true, Contents: response}, nil
	}

	rawSplit := strings.Split(response, ";")
	messageSplit := strings.Split(rawSplit[len(rawSplit)-1], ":")

	var item bot.Item

	msg := strings.TrimSpace(messageSplit[len(messageSplit)-1])

	// detect a potential command invocation, if you're confused on what the match means, look at the comments next to the regexp vars
	if msg[0] == '!' {
		if typeRegex.MatchString(msg) {
			item.Type = msg
		} else if commandNoContentRegex.MatchString(msg) {
			split := strings.Split(msg, " ")
			item.Type = split[0]
			item.Command = split[1]
			item.Key = split[2]
		} else if fullCommandRegex.MatchString(msg) {
			split := strings.Split(msg, " ")
			item.Type = split[0]
			item.Command = split[1]
			item.Key = split[2]
			item.Contents = strings.Join(split[3:], " ")
		} else {
			return bot.Item{}, bot.NonFatalError{Err: errComParse}
		}
	} else {
		// in this case, just a standard chat message
		item.Contents = msg
	}

	return item, nil
}
