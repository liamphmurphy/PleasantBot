package bot

import (
	"net"
	"testing"
)

type ConnMock struct {
	net.Conn
}

func (conn ConnMock) Write(b []byte) (n int, err error) {
	return -1, nil
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
		result, err := bot.HandlePing(test.inputMessage, test.pingIndicator)
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
