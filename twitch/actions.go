package twitch

import "github.com/liamphmurphy/pleasantbot/bot"

type ActionTaker interface {
	Action(payload ActionPayload, bot *bot.Bot, messenger bot.Messenger) error
}

// NoopAction will be used as a placeholder when no other action was found
type NoOpAction struct{}

type ActionPayload struct {
	Message string
}

type CommandAction struct {
	Key     string
	Command bot.CommandValue
}

func (ca *CommandAction) Action(payload ActionPayload, bot *bot.Bot, messenger bot.Messenger) error {
	return messenger.Message("oh hi there!")
}

func (noop *NoOpAction) Action(payload ActionPayload, bot *bot.Bot, messenger bot.Messenger) error {
	return nil
}

func setupDefaultActions() map[string]ActionTaker {
	return map[string]ActionTaker{"!help": &CommandAction{}}
}
