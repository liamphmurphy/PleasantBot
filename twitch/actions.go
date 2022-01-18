package twitch

import (
	"fmt"

	"github.com/liamphmurphy/pleasantbot/bot"
)

type ActionTaker interface {
	Condition(payload bot.Item, bot *bot.Bot) bool
	Action(payload bot.Item, bot *bot.Bot, messenger bot.Messenger) error
}

// NoopAction will be used as a placeholder when no other action was found
type NoOpAction struct{}

type CommandAction struct {
	Key     string
	Command bot.CommandValue
}

type PingAction struct{}

func (ca *CommandAction) Condition(item bot.Item, bot *bot.Bot) bool {
	return item.Type != ""
}

func (ca *CommandAction) Action(item bot.Item, bot *bot.Bot, messenger bot.Messenger) error {
	var err error
	var response string

	if item.Type == "!com" {
		switch item.Command {
		case "add":
			err = bot.AddCommand(item)
			response = fmt.Sprintf("%s was successfully added", item.Key)
		}
	} else {
		found, com := bot.FindCommand(item.Type)
		if found {
			messenger.Message(com.Response)
		}
	}

	if err == nil {
		messenger.Message(response)
	}

	return err
}

func (pa *PingAction) Condition(item bot.Item, bot *bot.Bot) bool {
	return true
}

func (noop *NoOpAction) Action(item bot.Item, bot *bot.Bot, messenger bot.Messenger) error {
	return nil
}

func (noop *NoOpAction) Condition(item bot.Item, bot *bot.Bot) bool { return true }

func setupDefaultActions() []ActionTaker {
	return []ActionTaker{&CommandAction{}}
}
