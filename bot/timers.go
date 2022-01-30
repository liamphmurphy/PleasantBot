package bot

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// TimedValues contains the values needed to run a single timed command.
type TimedValue struct {
	Message string
	Minutes int
	Enabled bool
}

// AddTimer takes in an item and parses it to add an associated timer. Will assume enabled by default.
func (bot *Bot) AddTimer(item Item) error {
	if _, ok := bot.Timers[item.Key]; ok {
		return fmt.Errorf("a timer with the key %s already exists", item.Key)
	}

	if item.Key == "" {
		return NonFatalError{Err: fmt.Errorf("the name for a new timer cannot be empty")}
	}

	values := strings.Split(item.Contents, " ")
	minutes, err := strconv.Atoi(values[0])
	if err != nil {
		return err
	}

	msg := strings.Join(values[1:], " ")
	bot.Timers[item.Key] = &TimedValue{Minutes: minutes, Message: msg, Enabled: true}
	if bot.Storage != nil {
		bot.Storage.DB.Insert("timers", bot.Storage.TimerColumns, []string{item.Key, msg, strconv.Itoa(minutes), "true"})
	}

	return nil
}

// DeleteTimer will delete a timer from the map and DB.
func (bot *Bot) DeleteTimer(item Item) error {
	if _, ok := bot.Timers[item.Key]; ok {
		delete(bot.Timers, item.Key)
		err := bot.Storage.DB.Delete("timers", "timername", item.Key)
		if err != nil {
			return FatalError{Err: err}
		}
	} else {
		return NonFatalError{Err: fmt.Errorf("the timer command '%s' does not exist", item.Key)}
	}

	return nil
}

// RunTimers will loop over and run any timers in the database that are enabled.
// A service who wants to use RunTimers must define a Messenger.Message definition so this function knows
// how to send the timer's contents.
func (bot *Bot) RunTimers(messenger Messenger) {
	for _, tv := range bot.Timers {
		go func(timedVal *TimedValue) {
			if tv.Enabled {
				for range time.NewTicker(time.Minute * time.Duration(timedVal.Minutes)).C {
					messenger.Message(timedVal.Message)
				}
			}

		}(tv)
	}
}

// LoadTimers gets all of the timers that exist in storage
func (bot *Bot) LoadTimers() error {
	rows, err := bot.Storage.DB.Query(fmt.Sprintf("select %s from timers", strings.Join(bot.Storage.TimerColumns, ",")))
	if err != nil {
		return err
	}

	defer rows.Close()
	for rows.Next() {
		var name, response string
		var minutes int
		var enabled bool
		err = rows.Scan(&name, &response, &minutes, &enabled)
		if err != nil {
			return err
		}
		bot.Timers[name] = &TimedValue{Message: response, Minutes: minutes, Enabled: enabled}
	}

	return nil
}
