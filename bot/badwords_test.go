package bot

import (
	"reflect"
	"testing"
)

func TestParseForBadWord(t *testing.T) {
	tests := []struct {
		description string
		phrase      string
		badWords    []BadWord
		wantFound   bool
		wantBadWord BadWord
	}{
		{
			description: "should find a bad word",
			badWords:    []BadWord{{Phrase: "cookies", Severity: 0}},
			phrase:      "cookies",
			wantFound:   true,
			wantBadWord: BadWord{Phrase: "cookies", Severity: 0},
		},
		{
			description: "should not find a bad word with a non-empty slice",
			badWords:    []BadWord{{Phrase: "not-a-cookie", Severity: 0}},
			phrase:      "cookies",
			wantFound:   false,
			wantBadWord: BadWord{},
		},
		{
			description: "should not find a bad word with an empty slice",
			badWords:    []BadWord{},
			phrase:      "cookies",
			wantFound:   false,
			wantBadWord: BadWord{},
		},
	}

	for _, test := range tests {
		bot := &Bot{BadWords: test.badWords}
		found, badWord := bot.ParseForBadWord(test.phrase)
		if !reflect.DeepEqual(found, test.wantFound) {
			t.Errorf("did not get expected found value\ngot - %v\nwant - %v", found, test.wantFound)
		}

		if !reflect.DeepEqual(badWord, test.wantBadWord) {
			t.Errorf("did not get expected BadWord struct\ngot - %v\nwant - %v", badWord, test.wantBadWord)
		}
	}
}
