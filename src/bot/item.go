// item.go will handle the creation / validation of a "item" request. I'm calling item any kind of new command requests
// so for actual commands, timers, etc.

package bot

import (
	"fmt"
	"regexp"
	"strings"
)

// Item represents a new key / value item for the bot such a command.
// e.g. a new item request would have a structure such as: !addcom !command <contents>
// This struct will hold the !command and <contents> values respectively.
// Each field will hold the value with any leading ! chars removed.
type Item struct {
	Type     string // ex: !addcom
	Key      string // ex: !somecommand
	Contents string // ex: this is the value of some command
}

var (
	itemRegex = regexp.MustCompile(`^(\![\w]*)\s(\![\w]*)\s(.)*`) // regexp for new item of form !newitem !key <some value>
)

// NewItem takes in a string request, confirms the structure, and creates a new Item struct
func NewItem(req string) (Item, error) {
	// confirm that the request matches the new item regex
	if !itemRegex.MatchString(req) {
		return Item{}, fmt.Errorf("the foillowing request does not match the correct form: %s", req)
	}
	var item Item
	split := strings.Split(req, " ")

	item.Type = strings.ReplaceAll(split[0], "!", "")
	item.Key = strings.ReplaceAll(split[1], "!", "")
	item.Contents = strings.Join(split[2:], " ")

	return item, nil
}
