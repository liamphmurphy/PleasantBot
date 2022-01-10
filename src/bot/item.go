// item.go will handle the creation / validation of a "item" request. I'm calling item any kind of new command requests
// so for actual commands, timers, etc.

package bot

import (
	"errors"
	"regexp"
	"strings"
)

// Item represents a new key / value item for the bot such a command.
// e.g. a new item request would have a structure such as: !addcom !command <contents>
// This struct will hold the !command and <contents> values respectively.
// Each field will hold the value with any leading ! chars removed.
type Item struct {
	Type     string // ex: addcom
	Key      string // ex: !somecommand
	Contents string // ex: this is the value of some command
}

var (
	keyValueContentRegex = regexp.MustCompile(`^(\![\w]*)\s(\![\w]*)\s(.)*`) // regexp for new item of form !itemcommand !key <some value>, such as "!addcom !com this is a test command"
	keyValueRegex        = regexp.MustCompile(`^(\![\w]*)\s(\![\w]*)$`)      // regexp for new item of form !itemcommand !key, such as !delcom !comtodelete
	keyRegex             = regexp.MustCompile(`^(\![\w]*)$`)                 // regexp for new item of form 1itemcommand, such as "!quote" (note the absence of any values / content)
	typeContentRegex     = regexp.MustCompile(`^(\![\w]*)\s([^!]*)$`)
)

// NewItem takes in a string request, confirms the structure, and creates a new Item struct
func NewItem(req string) (Item, error) {
	var item Item
	split := strings.Split(req, " ")
	splitLen := len(split)

	if splitLen == 0 {
		return item, errors.New("NewItem received an empty string request")
	}

	if keyRegex.MatchString(split[0]) {
		item.Type = strings.ReplaceAll(split[0], "!", "")
	}

	// covers the case of type and content, such as: !addquote <quote content>
	if typeContentRegex.MatchString(req) {
		item.Contents = strings.Join(split[1:], " ")
	} else { // covers remaining cases
		if splitLen > 1 && keyValueRegex.MatchString(strings.Join(split[:], " ")) {
			item.Key = split[1]
		} else if keyValueContentRegex.MatchString(req) {
			item.Key = split[1]
			item.Contents = strings.Join(split[2:], " ")
		}
	}

	return item, nil
}
