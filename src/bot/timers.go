package bot

// TimedValues contains the values needed to run a single timed command
type TimedValues struct {
	Message string
	Minutes int
}

func (bot *Bot) AddTimer(msg string) error {
	return nil
}
