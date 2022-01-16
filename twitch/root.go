// This file acts as the launch point for a Twitch bot

package twitch

import (
	"fmt"

	"github.com/liamphmurphy/pleasantbot/bot"
	"github.com/liamphmurphy/pleasantbot/storage"
)

var twitchServer = "irc.chat.twitch.tv:6697"

type Twitch struct {
	Bot *bot.Bot
}

// Run defines the main entry point for a Twitch bot
func (t *Twitch) Run() error {
	homeDir, err := bot.GetHomeDirectory()
	if err != nil {
		return err
	}

	v, err := bot.CreateViperConfig(homeDir, "twitch.toml", twitchServer)
	if err != nil {
		return err
	}

	database := bot.Database{Path: fmt.Sprintf("%s/pleasantbot.db", homeDir)}

	t.Bot, err = bot.CreateBot(v, database, storage.Init, bot.LoadBot)
	if err != nil {
		return bot.FatalError{Err: err}
	}

	return nil
}
