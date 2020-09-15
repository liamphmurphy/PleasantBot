// Handles all quote related actions. Because only a slice of strings is needed, there is no custom quote struct

package bot

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// AddQuote adds a quote to the bot's internal slice and the database
func (bot *Bot) AddQuote(quote string) error {
	// prepare the quote with an added date and time
	quoteWithDate := fmt.Sprintf("%s -- %s", quote, time.Now().Format("2006-01-02"))
	bot.Quotes = append(bot.Quotes, quoteWithDate)
	err := bot.InsertIntoDB("quotes", []string{"quote"}, []string{quoteWithDate})
	return err
}

// RandomQuote returns a random string, it does not print. It's up to the caller on what to do with it.
func (bot *Bot) RandomQuote() (string, error) {
	if len(bot.Quotes) == 0 {
		return "", errors.New("no quotes were found")
	}
	randomIndex := rand.Intn(len(bot.Quotes))
	return bot.Quotes[randomIndex], nil
}

// LoadQuotes loads all quotes from the DB.
// TODO: generalize this for all loading functions
func (bot *Bot) LoadQuotes() error {
	rows, err := bot.DB.Query("select quote from quotes")
	if err != nil {
		return err
	}

	defer rows.Close()
	for rows.Next() { // scan through results from query and assign to the Commands slice
		var quote string
		err = rows.Scan(&quote)
		if err != nil {
			return err
		}
		bot.Quotes = append(bot.Quotes, quote)
	}
	return nil
}

// GetQuote returns a quote of a specified index / id. Correlates to the automatically generated ID in sqlite.
func (bot *Bot) GetQuote(index int) (string, error) {
	if index <= 0 {
		return "", errors.New("the ID must be a valid integer greater than 0")
	} else if index > len(bot.Quotes) {
		return "", fmt.Errorf("the requested ID %d is greater than the total number of quotes, which is: %d", index, len(bot.Quotes))
	}
	return bot.Quotes[index-1], nil
}
