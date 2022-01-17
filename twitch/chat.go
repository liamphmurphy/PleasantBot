package twitch

import (
	"fmt"
	"strings"
	"time"
)

func (t *Twitch) Message(msg string) error {
	return t.Bot.WriteToConn(fmt.Sprintf("PRIVMSG #%s :%s", t.Bot.ChannelName, msg))
}

// Handler will contain the root logic for handling any kind of message from twitch; whether from the IRC server itself,
// or messages from the Twitch chat
func (t *Twitch) Handler(msg string, defaultActions map[string]ActionTaker) error {
	var at []ActionTaker

	// see if the message will prompt a default action
	for k, v := range defaultActions {
		if msg == k {
			at = append(at, v)
		}
	}

	if len(at) == 0 {
		at = append(at, &NoOpAction{})
	}

	// perform actions
	for _, act := range at {
		err := act.Action(ActionPayload{}, t.Bot, t)
		if err != nil {
			return err
		}
	}

	return nil
}

// purges a user by sending a timeout of 1 second
func (t *Twitch) purgeUser(username string) {
	t.Message(fmt.Sprintf("/timeout %s 1", username))
}

func (t *Twitch) banUser(username string, reason string) {
	t.Bot.Storage.DB.Insert("ban_history", []string{"user", "reason", "timestamp"}, []string{username, reason, time.Now().Format("2006-01-02 15:04:05")}) // insert into ban_history table
	t.Message(fmt.Sprintf("/ban %s", username))
}

func cleanUpTwitchResponse(response string) string {
	rawSplit := strings.Split(response, ";")
	messageSplit := strings.Split(rawSplit[len(rawSplit)-1], ":")
	return messageSplit[len(messageSplit)-1]
}
