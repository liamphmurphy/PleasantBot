package twitch

import (
	"fmt"
	"time"
)

func (t *Twitch) Message(msg string) error {
	return t.sendTwitchMessage(msg)
}

// purges a user by sending a timeout of 1 second
func (t *Twitch) purgeUser(username string) {
	t.Message(fmt.Sprintf("/timeout %s 1", username))
}

func (t *Twitch) banUser(username string, reason string) {
	t.Bot.Storage.DB.Insert("ban_history", []string{"user", "reason", "timestamp"}, []string{username, reason, time.Now().Format("2006-01-02 15:04:05")}) // insert into ban_history table
	t.Message(fmt.Sprintf("/ban %s", username))
}
