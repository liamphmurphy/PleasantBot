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
			description:  "should pass a simple request",
			inputRequest: "!newitem !someitem this is a test item",
			wantedItem:   Item{Type: "newitem", Key: "someitem", Contents: "this is a test item"},
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
