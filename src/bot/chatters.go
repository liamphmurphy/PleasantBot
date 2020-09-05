package bot

import (
	"fmt"
	"strings"
)

// User contains data on a single message and the user that sent it.
type User struct {
	Name         string // twitch username
	Content      string // the actual message
	IsSubscriber bool
	IsModerator  bool
	IsCommand    bool
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

// ParseMessage parses a message and its associated data.
func (bot *Bot) ParseMessage(message []string) User {
	var user User
	user.Name = parseMessageValue(message[4])
	user.IsSubscriber = parseMessageTrueOrFalse(message[10]) // determine if subscriber
	user.IsModerator = parseMessageTrueOrFalse(message[8])   // determine if moderator

	// This next bit feels horrendous, but it do be like that sometimes
	actualMessage := strings.Split(message[len(message)-1], fmt.Sprintf("PRIVMSG #%s :", bot.ChannelName))
	user.Content = actualMessage[1]

	if user.Content[0] == '!' && len(user.Content) > 1 {
		user.IsCommand = true
	}

	return user
}
