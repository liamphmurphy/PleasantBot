// This file handles all command related operations for the bot.
// The only time commands will be loaded from the bot's internal database is at initial start, otherwise as commands are added
// they'll be added separately to the bot's internal Commands slice and into the database.

package bot

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// CommandValue makes up a single command, used as the value in the bot's underyling commands map
type CommandValue struct {
	Response string
	Perm     string
	Count    int
}

// AddCommandString takes in a string of the form !addcom !comtitle <command response>
func (bot *Bot) AddCommandString(msg string) error {
	msgSplit := strings.Split(msg, " ")
	if len(msgSplit) < 2 || msgSplit[0][0] != '!' { // should contain at least the command name + one word or more as the response
		return errors.New("please make sure your addcom call is like so: !addcom !commandname <full text response for the command>")
	}

	key := msgSplit[0]
	bot.Commands[key] = &CommandValue{Response: (strings.Join(msgSplit[1:], " ")), Perm: "all", Count: 0}

	err := bot.InsertIntoDB("commands", bot.CommandDBColumns, []string{key, bot.Commands[key].Response, "1", "0"})
	if err != nil {
		return err
	}
	return nil
}

// FindCommand takes in a key (command name) and returns matching command, if found
func (bot *Bot) FindCommand(key string) (CommandValue, error) {
	var com *CommandValue
	com, found := bot.Commands[key]
	var err error
	if !found {
		err = errors.New("could not find command")
		return CommandValue{Response: "nil", Perm: "nil"}, err
	}
	return *com, err
}

// RemoveCommand takes in a command name as a string, presumably from the chat, and removes it
func (bot *Bot) RemoveCommand(key string) bool {
	var found bool
	if _, found := bot.Commands[key]; found {
		delete(bot.Commands, key)                               // deletes from the commands map
		err := bot.DeleteFromDB("commands", "commandname", key) // deletes permanently from the DB
		if err == nil {
			bot.SendMessage(fmt.Sprintf("%s has been deleted.", key))
		}
	}
	return found
}

// IncrementCommandCount takes in a command name (key) and increments the associated count value in the DB
func (bot *Bot) IncrementCommandCount(command string) error {
	var err error
	stmt := fmt.Sprintf("UPDATE commands SET count = count + 1 WHERE commandname = '%s'", command) // prepare statement tring
	_, err = bot.DB.Exec(stmt)
	if err != nil {
		return fmt.Errorf("Error updating the count for %s. Error: %s", command, err)
	}
	return nil
}

// DefaultCommands takes in a potential command request and sees if it is one of the default commands
func (bot *Bot) DefaultCommands(user User) bool {
	cmdFound := true
	msgSplit := strings.Split(user.Content, " ")

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

	case "!quote":
		if user.Content == "!quote" { // if entire message is just !quote, user is asking for a random quote
			quote, err := bot.RandomQuote()
			if err != nil {
				bot.SendMessage(fmt.Sprintf("Error finding a quote: %s", err))
				return false
			}
			bot.SendMessage(quote)
		} else { // otherwise, count on them at least attempting in getting a specific quote
			id, err := strconv.Atoi(msgSplit[1])
			if err != nil {
				bot.SendMessage("id for the quote must be a valid positive integer")
			}
			quote, err := bot.GetQuote(id)
			if err != nil {
				bot.SendMessage(fmt.Sprintf("Error finding a quote: %s", err))
				return false
			}
			bot.SendMessage(quote)
		}
	case "!addquote":
		bot.AddQuote(strings.Join(msgSplit[1:], " "), user.Name)

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
	rows, err := bot.DB.Query("select commandname, commandresponse, perm from commands")
	if err != nil {
		return err
	}

	defer rows.Close()
	for rows.Next() { // scan through results from query and assign to the Commands slice
		var name, response string
		var perm string
		err = rows.Scan(&name, &response, &perm)
		if err != nil {
			return err
		}
		bot.Commands[name] = &CommandValue{Response: response, Perm: perm}
	}
	return nil
}

// ConvertPermToInt takes in a string "all", "moderator" etc and converts it to the associated int.
func (bot *Bot) ConvertPermToInt(perm string) (uint8, error) {
	perm = strings.ToLower(perm)
	var err error

	switch perm { // determine permission, or error out if needed
	case "all":
		return 0, err

	case "subscriber":
		return 1, err

	case "moderator":
		return 2, err

	case "broadcaster":
		return 3, err

	default:
		return 255, fmt.Errorf("id not receive a valid permission")
	}
}
