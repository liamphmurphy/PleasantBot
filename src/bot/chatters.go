package bot

import (
	"fmt"
	"strings"
)

// User contains data on a single message and the user that sent it.
type User struct {
	Name      string // twitch username
	Content   string // the actual message
	Perm      uint8  // Represents permission level: 3 for broadcaster, 2 for moderator, 1 for subscriber, 0 for viewer. Used to check if the user can use a particular command.
	IsCommand bool
}

// used for debugging / testing. Shows the split of the message by index and content, useful to determine where to access which information.
func debugDisplayMessage(message []string) {
	for i := range message {
		fmt.Printf("i: %d -- %s\n", i, message[i])
	}
}

// a raw line from Twitch IRC will contain some values such as whether they are subscriber, their name etc.
// this takes in a line such as 'subscriber=0' and returns RHS of the equals sign.
func parseMessageValue(line string) string {
	return strings.Split(line, "=")[1]
}

// similar to parseMessageValue. Some twitch values are 0 or 1, so return false or true accordingly.
func parseMessageTrueOrFalse(line string) bool {
	value := parseMessageValue(line)
	if value == "0" {
		return false
	}

	return true
}

// parses a user's message and associates the appropriate uint8 value for their permission level.
func parsePermissions(message []string) uint8 {
	perm := 0

	if parseMessageTrueOrFalse(message[10]) { // checks if user is a subscriber
		if 1 > perm {
			perm = 1
		}
	}

	if parseMessageTrueOrFalse(message[8]) { // checks if user is a moderator
		if 2 > perm {
			perm = 2
		}
	}

	return uint8(perm)
}

// ParseMessage parses a message and its associated data.
func (bot *Bot) ParseMessage(message []string) User {
	var user User
	user.Name = parseMessageValue(message[4])
	user.Perm = parsePermissions(message)

	// This next bit feels horrendous, but it do be like that sometimes
	actualMessage := strings.Split(message[len(message)-1], fmt.Sprintf("PRIVMSG #%s :", bot.ChannelName))
	user.Content = actualMessage[1]

	if user.Content[0] == '!' && len(user.Content) > 1 {
		user.IsCommand = true
	}

	return user
}

// UpdateChatterCount increments the number of times a chatter has chatted
func (bot *Bot) UpdateChatterCount(user string) error {
	var err error

	// check if user is already in DB
	stmt := fmt.Sprintf("INSERT OR IGNORE INTO chatters (username, count) VALUES ('%s', '0')", user) // prepare statement string
	err = bot.DB.ArbitraryExec(stmt)
	if err != nil {
		return fmt.Errorf("Error inserting user %s into chatters. Error: %s", user, err)
	}

	// TODO: figure out way to generalize this and previous statement or reduce lines of code
	stmt = fmt.Sprintf("UPDATE chatters SET count = count + 1 WHERE username='%s'", user)
	err = bot.DB.ArbitraryExec(stmt)
	if err != nil {
		return fmt.Errorf("Error updating the count for %s. Error: %s", user, err)
	}
	return nil
}
