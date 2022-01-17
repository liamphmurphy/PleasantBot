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

func (t *Twitch) craftInitialConnMessages() []string {
	return []string{
		fmt.Sprintf("PASS %s", t.Bot.GetOAuth()),
		fmt.Sprintf("NICK %s", t.Bot.Name),
		fmt.Sprintf("JOIN #%s", t.Bot.ChannelName),
		"CAP REQ :twitch.tv/tags",
		"CAP REQ :twitch.tv/commands",
	}
}
