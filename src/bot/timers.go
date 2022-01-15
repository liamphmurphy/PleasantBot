package bot

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// TimedValues contains the values needed to run a single timed command
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

	values := strings.Split(item.Contents, " ")
	minutes, err := strconv.Atoi(values[0])
	if err != nil {
		return err
	}

	bot.Timers[item.Key] = &TimedValue{Minutes: minutes, Message: strings.Join(values[1:], " "), Enabled: true}

	return nil
}

func (bot *Bot) RunTimers() {
	for _, tv := range bot.Timers {
		go func(timedVal *TimedValue) {
			for range time.NewTicker(time.Minute * time.Duration(timedVal.Minutes)).C {
				bot.SendTwitchMessage(timedVal.Message)
			}
		}(tv)
	}
}

func (bot *Bot) LoadTimers() error {
	rows, err := bot.DB.Query(fmt.Sprintf("select %s from timers", strings.Join(bot.TimerDBColumns, ",")))
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
