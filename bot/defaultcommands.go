// default commands consist of the commands that you want to be available at all times regardless of user additions.
// The driver for the bot should make these default commands.

package bot

import "errors"

type DefaultFunc func(Item, *Bot) error

type DefaultCommand struct {
	Type     string
	Command  string
	ExecFunc DefaultFunc
}

var SliceEmptyErorr = "the bot's slice of default commands is empty"

// RunDefaultCommands will run a default command if caller contains the right invocation
func (bot *Bot) RunDefaultCommands(caller Item) (bool, error) {
	var hit bool
	var err error

	if len(bot.DefaultCommands) == 0 {
		return false, errors.New(SliceEmptyErorr)
	}

	for _, cmd := range bot.DefaultCommands {
		if (caller.Type == cmd.Type) && (caller.Command == cmd.Command) {
			err = cmd.ExecFunc(caller, bot) // if found, execute the tied function
			hit = true
			break
		}
	}

	return hit, err
}
