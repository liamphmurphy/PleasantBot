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

type CommandAction struct{}

type QuoteAction struct{}

type PingAction struct{}

func (ca *CommandAction) Condition(item bot.Item, bot *bot.Bot) bool {
	return item.Type != ""
}

// Action for a CommandAction defines the actions that can be made with a command
func (ca *CommandAction) Action(item bot.Item, bot *bot.Bot, messenger bot.Messenger) error {
	var err error
	var response string

	if item.Type == "!com" {
		switch item.Command {
		case "add", "new":
			err = bot.AddCommand(item)
			response = fmt.Sprintf("%s was successfully added", item.Key)
		case "del", "rm", "delete", "remove":
			var found bool
			found, err = bot.RemoveCommand(item.Key)
			if err == nil {
				if !found {
					response = fmt.Sprintf("%s does not exist", item.Key)
				} else {
					response = fmt.Sprintf("%s was successfully deleted", item.Key)
				}
			}
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

func (qa *QuoteAction) Condition(item bot.Item, bot *bot.Bot) bool {
	return item.Type == "!quote"
}

func (qa *QuoteAction) Action(item bot.Item, bot *bot.Bot, messenger bot.Messenger) error {
	var err error
	var response string

	switch item.Command {
	case "":
		response, err = bot.RandomQuote()
	case "add":
		err = bot.AddQuote(item.Contents, "nil")
	}

	if err == nil {
		messenger.Message(response)
	} else {
		messenger.Message(err.Error())
	}

	return err
}

func (pa *PingAction) Condition(item bot.Item, bot *bot.Bot) bool {
	return true
}

// Action for a NoOpAction returns a nil error, in other words, this is a stub that does nothing
func (noop *NoOpAction) Action(item bot.Item, bot *bot.Bot, messenger bot.Messenger) error {
	return nil
}

func (noop *NoOpAction) Condition(item bot.Item, bot *bot.Bot) bool { return true }

// setupDefaultActions prepares the default ActionTaker pipeline items
func setupDefaultActions() []ActionTaker {
	return []ActionTaker{&CommandAction{}, &QuoteAction{}}
}
