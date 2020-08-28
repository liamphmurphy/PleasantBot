// this file handles words that can cause purges and bans.
// timeouts may be supported later, but for now I'm thinking that's best left for manual moderation.

package bot

import (
	"fmt"
	"strings"
)

// BadWord contains info useful for bannable / purgeable phrases
type BadWord struct {
	Phrase   string
	Severity int // 0 for purge, 1 for perma ban
}

// purges a user by sending a timeout of 1 second
func (bot *Bot) purgeUser(username string) {
	bot.SendMessage(fmt.Sprintf("/timeout %s 1", username))
}

func (bot *Bot) banUser(username string) {
	bot.SendMessage(fmt.Sprintf("/ban %s", username))
}

// ParseForBadWord reads in every chat message and sees if a bad word was found in the message.
func (bot *Bot) ParseForBadWord(user User) {
	for i := range bot.BadWords { // search through all bad words
		if strings.Contains(user.Content, bot.BadWords[i].Phrase) {
			if bot.BadWords[i].Severity == 0 { // purge condition
				bot.purgeUser(user.Name)
			} else { // ban condition
				bot.banUser(user.Name)
			}
		}
	}
}
