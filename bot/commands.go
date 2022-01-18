// This file handles all command related operations for the bot.
// The only time commands will be loaded from the bot's internal database is at initial start, otherwise as commands are added
// they'll be added separately to the bot's internal Commands slice and into the database.

package bot

import (
	"fmt"
	"strings"
)

// CommandValue makes up a single command, used as the value in the bot's underyling commands map
type CommandValue struct {
	Response string `json:"response"`
	Perm     string `json:"perm"`
	Count    int    `json:"count"`
}

// AddCommandString takes in a string of the form !addcom !comtitle <command response>
func (bot *Bot) AddCommand(item Item) error {
	bot.Commands[item.Key] = &CommandValue{Response: item.Contents, Perm: "all", Count: 0}

	err := bot.Storage.DB.Insert("commands", bot.Storage.CommandColumns, []string{item.Key, bot.Commands[item.Key].Response, "all", "0"})
	if err != nil {
		return err
	}
	return nil
}

// FindCommand takes in a key (command name) and returns matching command, if found
func (bot *Bot) FindCommand(key string) (bool, CommandValue) {
	var comValue CommandValue
	com, found := bot.Commands[key]

	// have to handle some nil pointer logic so we don't provide a nil pointer to the caller
	if com == nil {
		comValue = CommandValue{}
	} else {
		comValue = *com
	}

	return found, comValue
}

// RemoveCommand takes in a command name as a string, presumably from the chat, and removes it
func (bot *Bot) RemoveCommand(key string) (bool, error) {
	var found bool
	if _, found = bot.Commands[key]; found {
		delete(bot.Commands, key)                                    // deletes from the commands map
		err := bot.Storage.DB.Delete("commands", "commandname", key) // deletes permanently from the DB
		if err != nil {
			return found, err
		}
	}
	return found, nil
}

// IncrementCommandCount takes in a command name (key) and increments the associated count value in the DB
func (bot *Bot) IncrementCommandCount(command string) error {
	stmt := fmt.Sprintf("UPDATE commands SET count = count + 1 WHERE commandname = '%s'", command) // prepare statement tring
	err := bot.Storage.DB.ArbitraryExec(stmt)
	if err != nil {
		return fmt.Errorf("Error updating the count for %s. Error: %s", command, err)
	}
	return nil
}

// LoadCommands queries the sqlite3 database for existing commands
func (bot *Bot) LoadCommands() error {
	rows, err := bot.Storage.DB.Query("select commandname, commandresponse, perm from commands")
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
		if name[0] != '!' {
			name = fmt.Sprintf("!%s", name)
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
		return 255, fmt.Errorf("did not receive a valid permission: %s", perm)
	}
}
