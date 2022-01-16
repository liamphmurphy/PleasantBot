package bot

import (
	"fmt"
	"net"
	"reflect"
	"testing"

	"github.com/liamphmurphy/pleasantbot/storage"
	"github.com/spf13/viper"
)

type ConnMock struct {
	net.Conn
}

func (conn ConnMock) Write(b []byte) (n int, err error) {
	return -1, nil
}

func mockInitNoErr(path string, sq *storage.Sqlite) error { return nil }

func loaderStubNoError(bot *Bot) error { return nil }

func TestCreateBot(t *testing.T) {
	tests := []struct {
		description     string
		inputViper      *viper.Viper
		inputInitFunc   storage.InitFunc
		inputLoaderFunc BotLoaderFunc
		wantBot         *Bot
		wantErr         error
	}{
		{
			description:     "should succeed creating a standard bot",
			inputViper:      &viper.Viper{},
			inputInitFunc:   mockInitNoErr,
			inputLoaderFunc: loaderStubNoError,
			wantBot:         &Bot{Storage: &Database{}, Config: &viper.Viper{}},
			wantErr:         nil,
		},
		{
			description:     "should fail due to a nil viper struct",
			inputViper:      nil,
			inputInitFunc:   mockInitNoErr,
			inputLoaderFunc: LoadBot,
			wantBot:         &Bot{},
			wantErr:         fmt.Errorf("a fatal error occurred: %s", errNoViper),
		},
	}

	for _, test := range tests {
		bot, err := CreateBot(test.inputViper, Database{}, test.inputInitFunc, test.inputLoaderFunc)
		if err != nil {
			if test.wantErr == nil {
				t.Errorf("got an error when none was expected: %v", err)
			} else if test.wantErr.Error() != err.Error() {
				t.Errorf("did not get the expected error\ngot - %v\nwant - %v", err, test.wantErr)
			}
		}

		if !reflect.DeepEqual(bot, test.wantBot) {
			t.Errorf("did not get the expected bot\ngot - %v\nwant - %v", *bot, test.wantBot)
		}
	}
}

func TestHandlePing(t *testing.T) {
	tests := []struct {
		description   string
		inputMessage  string
		pingIndicator string
		wantResult    bool
		wantErr       error
	}{
		{
			description:   "should handle a standard Twitch PING",
			inputMessage:  "PING :tmi.twitch.tv",
			pingIndicator: "PING :tmi.twitch.tv",
			wantResult:    true,
			wantErr:       nil,
		},
	}

	for _, test := range tests {
		bot := &Bot{Conn: &ConnMock{}}
		result, err := bot.HandlePing(test.inputMessage, test.pingIndicator, "PONG")
		if err != nil {
			if err.Error() != test.wantErr.Error() {
				t.Errorf("got an unexpected error\nwant - %v\ngot - %v", test.wantErr, err)
			}
		}

		if result != test.wantResult {
			t.Errorf("did not get expected result\nwant - %v\ngot - %v", test.wantResult, result)
		}
	}
}
