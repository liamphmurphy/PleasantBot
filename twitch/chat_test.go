package twitch

import (
	"reflect"
	"testing"

	"github.com/liamphmurphy/pleasantbot/bot"
)

func TestNewTwitchItem(t *testing.T) {
	tests := []struct {
		description string
		inputMsg    string
		wantItem    bot.Item
		wantErr     error
	}{
		{
			description: "should process a standard chat message",
			inputMsg:    "@badge-info=subscriber/91;badges=broadcaster/1,subscriber/3000,premium/1;client-nonce=2d59456ff9792c4aa9521d53f109091b;color=#D3D3D3;display-name=LimePH;emotes=;first-msg=0;flags=;id=a6416f66-c477-47e2-ad6c-44c38a20f919;mod=0;room-id=26692942;subscriber=1;tmi-sent-ts=1642452235079;turbo=0;user-id=26692942;user-type= :limeph!limeph@limeph.tmi.twitch.tv PRIVMSG #limeph :test message",
			wantItem:    bot.Item{Contents: "test message"},
			wantErr:     nil,
		},
		{
			description: "should identify when it is a non-user server message",
			inputMsg:    ":tmi.twitch.tv 372 whitegirlcoffeebot :You are in a maze of twisty passages, all alike.",
			wantItem:    bot.Item{IsServerInfo: true, Contents: ":tmi.twitch.tv 372 whitegirlcoffeebot :You are in a maze of twisty passages, all alike."},
			wantErr:     nil,
		},
		{
			description: "detect a case of a command invocation without any key, e.g. !quote or !help.",
			inputMsg:    "@badge-info=subscriber/91;badges=broadcaster/1,subscriber/3000,premium/1;client-nonce=2d59456ff9792c4aa9521d53f109091b;color=#D3D3D3;display-name=LimePH;emotes=;first-msg=0;flags=;id=a6416f66-c477-47e2-ad6c-44c38a20f919;mod=0;room-id=26692942;subscriber=1;tmi-sent-ts=1642452235079;turbo=0;user-id=26692942;user-type= :limeph!limeph@limeph.tmi.twitch.tv PRIVMSG #limeph :!quote",
			wantItem:    bot.Item{Type: "!quote"},
			wantErr:     nil,
		},
		{
			description: "detect a case of a full command invocation, in this example, !",
			inputMsg:    "@badge-info=subscriber/91;badges=broadcaster/1,subscriber/3000,premium/1;client-nonce=2d59456ff9792c4aa9521d53f109091b;color=#D3D3D3;display-name=LimePH;emotes=;first-msg=0;flags=;id=a6416f66-c477-47e2-ad6c-44c38a20f919;mod=0;room-id=26692942;subscriber=1;tmi-sent-ts=1642452235079;turbo=0;user-id=26692942;user-type= :limeph!limeph@limeph.tmi.twitch.tv PRIVMSG #limeph :!com add !somecommand this is a test command",
			wantItem:    bot.Item{Type: "!com", Command: "add", Key: "!somecommand", Contents: "this is a test command"},
			wantErr:     nil,
		},
		{
			description: "detect a case of an improper command invocation",
			inputMsg:    "@badge-info=subscriber/91;badges=broadcaster/1,subscriber/3000,premium/1;client-nonce=2d59456ff9792c4aa9521d53f109091b;color=#D3D3D3;display-name=LimePH;emotes=;first-msg=0;flags=;id=a6416f66-c477-47e2-ad6c-44c38a20f919;mod=0;room-id=26692942;subscriber=1;tmi-sent-ts=1642452235079;turbo=0;user-id=26692942;user-type= :limeph!limeph@limeph.tmi.twitch.tv PRIVMSG #limeph :!com add keywithoutaexclamationmark woopsie daisy",
			wantItem:    bot.Item{},
			wantErr:     bot.FatalError{Err: errComParse},
		},
	}

	for _, test := range tests {
		item, err := newTwitchItem(test.inputMsg)
		if err != nil {
			if test.wantErr == nil {
				t.Errorf("got an error when we wanted nil: %v", err)
			} else if err.Error() != test.wantErr.Error() {
				t.Errorf("did not get an expected error\ngot - %v\nwant - %v", err, test.wantErr)
			}
		}

		if !reflect.DeepEqual(item, test.wantItem) {
			t.Errorf("did not get an expected Item\ngot - %v\nwant - %v", item, test.wantItem)
		}
	}
}
