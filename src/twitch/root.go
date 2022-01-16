// This file acts as the launch point for a Twitch bot

package twitch

import (
	"fmt"

	"github.com/liamphmurphy/pleasantbot/bot"
)

type Twitch struct {
	Bot *bot.Bot
}

func (t *Twitch) Run() error {
	fmt.Println("allo gov")
	return nil
}
