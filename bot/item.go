// item.go will handle the creation / validation of a "item" request. I'm calling item any kind of new command requests
// so for actual commands, timers, etc.

package bot

// Item represents a new key / value item for the bot such as a command.
// e.g. a new item request would have a structure such as: !addcom !command <contents>
// This struct will hold the !command and <contents> values respectively.
// Each field will hold the value with any leading ! chars removed.
type Item struct {
	IsServerInfo bool   // if true, consider item not from user and can be ignored
	Type         string // ex: com
	Command      string // ex: add
	Key          string // ex: !somecommand
	Contents     string // ex: this is the value of some command
}
