package bot

import (
	"reflect"
	"testing"
)

func TestNewItem(t *testing.T) {
	tests := []struct {
		description  string
		inputRequest string
		wantedItem   Item
	}{
		{
			description:  "should get a full type / key / content item",
			inputRequest: "!newitem !someitem this is a test item",
			wantedItem:   Item{Type: "newitem", Key: "!someitem", Contents: "this is a test item"},
		},
		{
			description:  "should get a type / key item",
			inputRequest: "!delcom !acommand",
			wantedItem:   Item{Type: "delcom", Key: "!acommand"},
		},
		{
			description:  "should get a type / content item",
			inputRequest: "!addquote this is a test quote",
			wantedItem:   Item{Type: "addquote", Contents: "this is a test quote"},
		},
		{
			description:  "should get a type item",
			inputRequest: "!quote",
			wantedItem:   Item{Type: "quote"},
		},
	}

	for _, test := range tests {
		item, err := NewItem(test.inputRequest)
		if err != nil {
			t.Errorf("did not get an expected error: %s", err)
		}

		if !reflect.DeepEqual(item, test.wantedItem) {
			t.Errorf("did not get an expected item:\nwant - %v\ngot - %v", test.wantedItem, item)
		}
	}
}
