package test

import (
	"regexp"
	"testing"

	"github.com/murnux/pleasantbot/bot"
)

func TestParseMessage(t *testing.T) {
	testBot := bot.CreateBot()
	lineToTest := "@badge-info=subscriber/75;badges=broadcaster/1,subscriber/3000,premium/1;client-nonce=b0fccf06b9fced4a5cd0cd04546f965b;color=#D3D3D3;display-name=LimePH;emote-only=1;emotes=25:0-4;flags=;id=5db09c3e-92dd-4a34-9eea-a9c3cea76dcc;mod=0;room-id=26692942;subscriber=1;tmi-sent-ts=1600141270651;turbo=0;user-id=26692942;user-type= :limeph!limeph@limeph.tmi.twitch.tv PRIVMSG #limeph :Kappa"
	msgRegex, _ := regexp.Compile("[;]+")
	lineSplit := msgRegex.Split(lineToTest, -1)

	user := testBot.ParseMessage(lineSplit)
	if !(user.Name == "LimePH" && user.Content == "Kappa" && user.IsSubscriber && user.IsModerator && !user.IsCommand) {
		t.Errorf("ParseMessage did not return the right values.")
	}
}
