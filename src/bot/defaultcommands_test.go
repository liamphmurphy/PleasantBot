package bot

import (
	"errors"
	"testing"
)

func stubNoError(caller Item, bot *Bot) error { return nil }

func stubWithError(caller Item, bot *Bot) error { return errors.New("stub error") }

func TestRunDefaultCommands(t *testing.T) {
	tests := []struct {
		description string
		defaults    []DefaultCommand
		inputItem   Item
		wantResult  bool
		wantErr     error
	}{
		{
			description: "should find the default command without error",
			defaults:    []DefaultCommand{{Type: "test", Command: "command", ExecFunc: stubNoError}},
			inputItem:   Item{Type: "test", Command: "command"},
			wantResult:  true,
			wantErr:     nil,
		},
		{
			description: "should be found but exec func hits an error",
			defaults:    []DefaultCommand{{Type: "test", Command: "command", ExecFunc: stubWithError}},
			inputItem:   Item{Type: "test", Command: "command"},
			wantResult:  true,
			wantErr:     errors.New("stub error"),
		},
		{
			description: "should error on an empty list of default commands",
			defaults:    []DefaultCommand{},
			inputItem:   Item{Type: "test", Command: "command"},
			wantResult:  false,
			wantErr:     errors.New(SliceEmptyErorr),
		},
		{
			description: "should fail to find a default command",
			defaults:    []DefaultCommand{{Type: "some-other-type", Command: "command", ExecFunc: stubNoError}},
			inputItem:   Item{Type: "test", Command: "command"},
			wantResult:  false,
			wantErr:     nil,
		},
	}
	var bot Bot
	for _, test := range tests {
		bot.DefaultCommands = test.defaults
		found, err := bot.RunDefaultCommands(test.inputItem)
		if test.wantErr == nil {
			if err != nil {
				t.Errorf("got a nil error returned but expected: %v", test.wantErr)
			}
		}

		// test if we got the expected error (if any)
		if err != nil && (test.wantErr.Error() != err.Error()) {
			t.Errorf("did not get an expected error\ngot - %v\nwant - %v", err, test.wantErr)
		}

		if found != test.wantResult {
			t.Errorf("did not get the expected result\ngot - %v\nwant - %v", found, test.wantResult)
		}
	}
}
