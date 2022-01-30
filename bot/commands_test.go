package bot

import (
	"reflect"
	"testing"
)

func TestFindCommand(t *testing.T) {
	tests := []struct {
		description string
		key         string
		commands    map[string]*CommandValue
		wantFound   bool
		wantCommand CommandValue
	}{
		{
			description: "should find and delete a command",
			key:         "cookies",
			commands:    map[string]*CommandValue{"cookies": &CommandValue{Response: "mmm cookies"}},
			wantFound:   true,
			wantCommand: CommandValue{Response: "mmm cookies"},
		},
		{
			description: "should not find a command",
			key:         "cookies",
			commands:    map[string]*CommandValue{"cupcake": &CommandValue{Response: "mmm cupcake"}},
			wantFound:   false,
			wantCommand: CommandValue{},
		},
	}

	for _, test := range tests {
		bot := &Bot{Commands: test.commands}
		found, com := bot.FindCommand(test.key)

		if !reflect.DeepEqual(found, test.wantFound) {
			t.Errorf("did not get expected found value\ngot - %v\nwant - %v", found, test.wantFound)
		}

		if !reflect.DeepEqual(com, test.wantCommand) {
			t.Errorf("did not get expected BadWord struct\ngot - %v\nwant - %v", com, test.wantCommand)
		}
	}
}

func TestEditCommand(t *testing.T) {
	tests := []struct {
		description string
		inputItem   Item
		commands    map[string]*CommandValue
		wantCommand *CommandValue
		wantErr     error
	}{
		{
			description: "should succeed to edit a command",
			inputItem:   Item{Type: "!com", Command: "edit", Key: "!test", Contents: "this is a new response"},
			commands: map[string]*CommandValue{
				"!test": {
					Response: "this is the old response",
					Perm:     "",
					Count:    0,
				},
			},
			wantCommand: &CommandValue{Response: "this is a new response", Perm: "", Count: 0},
			wantErr:     nil,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			bot := &Bot{Commands: test.commands}
			err := bot.EditCommand(test.inputItem)
			if err != nil {
				if test.wantErr == nil {
					t.Errorf("did not get the expected err\ngot - %v\nwant - nil", err)
				} else if err.Error() != test.wantErr.Error() {
					t.Errorf("did not get the expected error\ngot - %v\nwant - %v", err, test.wantErr)
				}
			}

			command, _ := bot.Commands[test.inputItem.Key]
			if !reflect.DeepEqual(command, test.wantCommand) {
				t.Errorf("did not get the expected command\ngot - %v\nwant - %v", command, test.wantCommand)
			}
		})
	}
}
