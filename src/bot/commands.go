package bot

import (
	"errors"
	"fmt"
	"strings"
)

// Command makes up a single command
type Command struct {
	Name           string
	Response       string
	ModeratorPerms bool
}

// takes in an already made command and appends it to the commands slice
func (bot *Bot) addCommand(command Command) bool {
	bot.Commands = append(bot.Commands, command)
	return true
}

// takes in a string of the form !newcom !comtitle <command response>
func (bot *Bot) addCommandString(msg string) bool {
	msgSplit := strings.Split(msg, " ")
	if len(msgSplit) < 3 { // should be at least 3 elements
		bot.SendMessage("New command should be of the form: !newcom !commandtitle <command content>")
		return false
	}
	fmt.Println(msgSplit)
	var newCmd Command // create new command
	newCmd.Name = msgSplit[1]
	newCmd.Response = strings.Join(msgSplit[2:], " ")
	newCmd.ModeratorPerms = true
	bot.Commands = append(bot.Commands, newCmd)
	return true
}

// FindCommand takes in a key (command name) and returns matching command, if found
func (bot *Bot) FindCommand(key string) (Command, error) {
	var com Command

	// linear search for command. TODO: make this better
	for i := range bot.Commands {
		fmt.Println(bot.Commands[i].Name)
		if bot.Commands[i].Name == key {
			return bot.Commands[i], nil
		}
	}
	return com, errors.New("could not find command")
}

// DefaultCommands takes in a potential command request and sees if it is one of the default commands
func (bot *Bot) DefaultCommands(msg string) bool {
	cmdFound := true
	msgSplit := strings.Split(msg, " ")
	switch msgSplit[0] { // start cycling through potential default commands
	case "!help":
		bot.SendMessage("Some helpful help message. :)")
	case "!addcom": // add a new custom command
		bot.addCommandString(msg)
	case "!subon":
		bot.SendMessage("/subscribers")
		bot.SendMessage("Subscriber only mode is now on.")
	case "!suboff":
		bot.SendMessage("/subscribersoff")
		bot.SendMessage("Subscriber only mode is now off.")
	default:
		cmdFound = false
	}

	return cmdFound
}
