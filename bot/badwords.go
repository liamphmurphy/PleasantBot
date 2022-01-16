// this file handles words that can cause purges and bans.
// timeouts may be supported later, but for now I'm thinking that's best left for manual moderation.

package bot

import (
	"strings"
)

// BadWord contains info useful for bannable / purgeable phrases
type BadWord struct {
	Phrase   string
	Severity int // 0 for purge, 1 for perma ban
}

// ParseForBadWord reads in every chat message and sees if a bad word was found in the message.
func (bot *Bot) ParseForBadWord(user User) {
	for i := range bot.BadWords { // search through all bad words
		if strings.Contains(user.Content, bot.BadWords[i].Phrase) {
			if bot.BadWords[i].Severity == 0 { // purge condition
				bot.purgeUser(user.Name)
			} else { // ban condition
				bot.banUser(user.Name, "inputting a banned word")
			}
		}
	}
}

// LoadBadWords loads all badwords from the databases
// TODO: generalize this for all bot data
func (bot *Bot) LoadBadWords() error {
	rows, err := bot.Storage.DB.Query("select phrase, severity from badwords")
	if err != nil {
		return err
	}

	defer rows.Close()
	for rows.Next() { // scan through results from query and assign to the Commands slice
		var phrase string
		var severity int
		err = rows.Scan(&phrase, &severity)
		if err != nil {
			return err
		}
		badWord := BadWord{Phrase: phrase, Severity: severity}
		bot.BadWords = append(bot.BadWords, badWord)
	}
	return nil
}
