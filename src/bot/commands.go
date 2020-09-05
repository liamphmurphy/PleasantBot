// This file handles all command related operations for the bot.
// The only time commands will be loaded from the bot's internal database is at initial start, otherwise as commands are added
// they'll be added separately to the bot's internal Commands slice and into the database.

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

// AddCommandString takes in a string of the form !addcom !comtitle <command response>
func (bot *Bot) AddCommandString(msg string) error {
	msgSplit := strings.Split(msg, " ")
	if len(msgSplit) < 2 || msgSplit[0][0] != '!' { // should contain at least the command name + one word or more as the response
		return errors.New("Please make sure your addcom call is like so: !addcom !commandname <full text response for the command>")
	}

	var newCmd Command // create new command
	newCmd.Name = msgSplit[0]
	newCmd.Response = strings.Join(msgSplit[1:], " ")
	newCmd.ModeratorPerms = true
	bot.Commands = append(bot.Commands, newCmd)
	err := bot.InsertIntoDB("commands", []string{"commandname", "commandresponse", "modperms"}, []string{newCmd.Name, newCmd.Response, "0"})
	if err != nil {
		return err
	}
	return nil
}

// FindCommand takes in a key (command name) and returns matching command, if found
func (bot *Bot) FindCommand(key string) (Command, error) {
	var com Command

	// linear search for command. TODO: make this better
	for i := range bot.Commands {
		if bot.Commands[i].Name == key {
			return bot.Commands[i], nil
		}
	}
	return com, errors.New("could not find command")
}

// RemoveCommand takes in a command name as a string, presumably from the chat, and removes it from the slice
func (bot *Bot) RemoveCommand(cmd string) bool {
	cmdFound := false
	index := 0
	for i := range bot.Commands {
		if bot.Commands[i].Name == cmd {
			index = i
			cmdFound = true
		}
	}

	// many thanks to this article for the best way to remove an item from a slice https://www.delftstack.com/howto/go/how-to-delete-an-element-from-a-slice-in-golang/
	if cmdFound {
		bot.Commands[index] = bot.Commands[len(bot.Commands)-1]
		bot.Commands = bot.Commands[:len(bot.Commands)-1]
	}
	return cmdFound
}

// DefaultCommands takes in a potential command request and sees if it is one of the default commands
func (bot *Bot) DefaultCommands(msg string) bool {
	cmdFound := true
	msgSplit := strings.Split(msg, " ")

	switch msgSplit[0] { // start cycling through potential default commands
	case "!help":
		bot.SendMessage("Some helpful help message. :)")

	case "!addcom": // add a new custom command
		err := bot.AddCommandString(strings.Join(msgSplit[1:], " "))
		if err == nil {
			bot.SendMessage("The new command was added successfully.")
		} else {
			bot.SendMessage(fmt.Sprintf("The bot returned the following error: %s", err))
		}

	case "!delcom":
		bot.RemoveCommand(msgSplit[1])

	case "!subon": // turn on subscribers only mode
		bot.SendMessage("/subscribers")
		bot.SendMessage("Subscriber only mode is now on.")
	case "!suboff": // turn off subscribers only mode
		bot.SendMessage("/subscribersoff")
		bot.SendMessage("Subscriber only mode is now off.")
	default:
		cmdFound = false
	}

	return cmdFound
}

// LoadCommands queries the sqlite3 database for existing commands
func (bot *Bot) LoadCommands() error {
	rows, err := bot.DB.Query("select commandname, commandresponse, modperms from commands")
	if err != nil {
		return err
	}

	defer rows.Close()
	for rows.Next() { // scan through results from query and assign to the Commands slice
		var name, response string
		var perm int
		err = rows.Scan(&name, &response, &perm)
		if err != nil {
			return err
		}
		com := Command{Name: name, Response: response, ModeratorPerms: Itob(perm)}
		bot.Commands = append(bot.Commands, com)
	}
	return nil
}
