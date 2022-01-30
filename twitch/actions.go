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

type TimerAction struct{}

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
			if err == nil {
				response = fmt.Sprintf("%s was successfully added", item.Key)
			}
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
		case "edit":
			err = bot.EditCommand(item)
			if err == nil {
				response = fmt.Sprintf("'%s' has been updated.", item.Key)
			}
		}
		// handle finding custom commands
	} else {
		found, com := bot.FindCommand(item.Type)
		if found {
			messenger.Message(com.Response)
		}
	}

	if err == nil {
		messenger.Message(response)
	} else {
		messenger.Message(err.Error())
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
	case "add", "new":
		err = bot.AddQuote(item.Contents, item.Sender.Name)
		if err == nil {
			response = fmt.Sprintf("new quote added by @%s", item.Sender.Name)
		}
	case "del", "rm", "delete", "remove":
		err = bot.DeleteQuote(item.Key)
		if err == nil {
			response = fmt.Sprintf("deleted quote with ID: %s", item.Key)
		}
	}

	if err == nil {
		messenger.Message(response)
	} else {
		messenger.Message(err.Error())
	}

	return err
}

func (ta *TimerAction) Condition(item bot.Item, bot *bot.Bot) bool {
	return item.Type == "!timer"
}

func (ta *TimerAction) Action(item bot.Item, bot *bot.Bot, messenger bot.Messenger) error {
	var err error
	var response string

	switch item.Command {
	case "add", "new":
		err = bot.AddTimer(item)
		if err != nil {
			response = "new timer has been added"
		}
	case "del", "rm", "delete", "remove":
		err = bot.DeleteTimer(item)
		if err != nil {
			response = fmt.Sprintf("'%s' has been removed", item.Key)
		}
	}

	if err == nil {
		messenger.Message(response)
	} else {
		messenger.Message(err.Error())
	}

	return err
}

// Action for a NoOpAction returns a nil error, in other words, this is a stub that does nothing
func (noop *NoOpAction) Action(item bot.Item, bot *bot.Bot, messenger bot.Messenger) error {
	return nil
}

func (noop *NoOpAction) Condition(item bot.Item, bot *bot.Bot) bool { return true }

// setupDefaultActions prepares the default ActionTaker pipeline items
func setupDefaultActions() []ActionTaker {
	return []ActionTaker{&CommandAction{}, &QuoteAction{}, &TimerAction{}}
}
