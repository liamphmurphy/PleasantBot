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
