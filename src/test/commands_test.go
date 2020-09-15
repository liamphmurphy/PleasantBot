package test

import (
	"testing"

	"github.com/murnux/pleasantbot/bot"
)

func TestAddCommandString(t *testing.T) {
	testBot := bot.CreateBot()

	command := "!addcom !testadd This is a test command."
	added := testBot.AddCommandString(command)

	// test that the command was added
	if added == nil || !(len(testBot.Commands) == 1) && testBot.Commands[0].Name == "!testadd" {
		t.Errorf("the command '%s' was not added correctly.", command)
	}
}

func TestFindCommand(t *testing.T) {
	testBot := bot.CreateBot()

	commandsToAdd := make([]string, 3)
	// create some test commands
	commandsToAdd[0] = "!addcom !test This is a test command"
	commandsToAdd[1] = "!addcom !4head JUST DO IT 4Head"
	commandsToAdd[2] = "!addcom !what chicken butt"

	for i := range commandsToAdd { // add test commands
		testBot.AddCommandString(commandsToAdd[i])
	}

	for i := range testBot.Commands {
		_, err := testBot.FindCommand(testBot.Commands[i].Name) // start testing FindCommand
		if err != nil {
			t.Errorf("Had trouble finding the %s command. Something is wrong.", testBot.Commands[i].Name)
		}
	}
}

func TestRemoveCommand(t *testing.T) {
	testBot := bot.CreateBot()

	commandsToAdd := make([]string, 3)
	// create some test commands
	commandsToAdd[0] = "!addcom !test This is a test command"
	commandsToAdd[1] = "!addcom !4head JUST DO IT 4Head"
	commandsToAdd[2] = "!addcom !what chicken butt"

	for i := range commandsToAdd { // add test commands
		testBot.AddCommandString(commandsToAdd[i])
	}

	testBot.AddCommandString("!addcom !remove Pleae remove this command. :)")

	removed := testBot.RemoveCommand("!what")
	if !removed && len(testBot.Commands) == 3 {
		t.Errorf("Failed to remove command: %s", "!remove")
	}
}
