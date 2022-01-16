package bot

import (
	"errors"
	"reflect"
	"testing"
)

var timers = map[string]*TimedValue{"dup-key": &TimedValue{}}

func TestAddTimer(t *testing.T) {
	tests := []struct {
		description string
		inputItem   Item
		wantedTimer TimedValue
		wantErr     error
	}{
		{
			description: "simple add timer should process correctly",
			inputItem:   Item{Key: "timer-key", Contents: "5 this is what the timer will print"},
			wantedTimer: TimedValue{Message: "this is what the timer will print", Minutes: 5, Enabled: true},
			wantErr:     nil,
		},
		{
			description: "should fail to add since key already exists",
			inputItem:   Item{Key: "dup-key"},
			wantedTimer: TimedValue{},
			wantErr:     errors.New("a timer with the key dup-key already exists"),
		},
		{
			description: "should fail when minutes isn't provided",
			inputItem:   Item{Key: "test-key", Contents: "I forgot the minutes, oh no!"},
			wantedTimer: TimedValue{},
			wantErr:     errors.New(`strconv.Atoi: parsing "I": invalid syntax`),
		},
	}

	for _, test := range tests {
		bot := Bot{Timers: timers}
		err := bot.AddTimer(test.inputItem)
		if err == nil && test.wantErr != nil {
			t.Errorf("got an unexpected error\ngot - %v\nwant - %v", err, test.wantErr)
		}

		if err != nil {
			if err.Error() != test.wantErr.Error() {
				t.Errorf("got an unexpected error\ngot - %v\nwant - %v", err, test.wantErr)
			}
		} else {
			timer := timers[test.inputItem.Key]
			if !reflect.DeepEqual(*timer, test.wantedTimer) {
				t.Errorf("did not get back the expected timer\ngot - %v\nwant - %v", timer, test.wantedTimer)
			}
		}
	}
}
