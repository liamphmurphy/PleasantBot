// this file will have some of the handling of Twitch API stuff

package twitch

import (
	"fmt"

	"github.com/liamphmurphy/pleasantbot/bot"
)

// initialConn passes in the needed info to twitch's IRC server; channel name, oauth token, etc.
func (t *Twitch) initialConn(messages []string) error {
	if len(messages) == 0 {
		return bot.FatalError{Err: fmt.Errorf("to connect to twitch, some initial conn messages must be provided")}
	}

	// send initial conn messages to Twitch
	for _, message := range messages {
		err := t.Bot.WriteToConn(message)
		if err != nil {
			return err
		}
	}

	return nil
}

// SendTwitchMessage prepares and sends a string to the channel's Twitch chat. It is recommended that the twitch.Message
// func be called instead of this.
func (t *Twitch) sendTwitchMessage(msg string) error {
	return t.Bot.WriteToConn(fmt.Sprintf("PRIVMSG #%s :%s", t.Bot.ChannelName, msg))
}

func (t *Twitch) craftInitialConnMessages() []string {
	return []string{
		fmt.Sprintf("PASS %s", t.Bot.GetOAuth()),
		fmt.Sprintf("NICK %s", t.Bot.Name),
		fmt.Sprintf("JOIN #%s", t.Bot.ChannelName),
		"CAP REQ :twitch.tv/tags",
		"CAP REQ :twitch.tv/commands",
	}
}
