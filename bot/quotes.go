// Handles all quote related actions. Because only a slice of strings is needed, there is no custom quote struct

package bot

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

// QuoteValues represents the values associated with a quote. The ID in the DB will be the map key
type QuoteValues struct {
	Quote     string
	Timestamp string
	Submitter string
}

// AddQuote adds a quote to the bot's internal slice and the database
func (bot *Bot) AddQuote(quote string, submitter string) error {
	// prepare the quote with an added date and time
	date := time.Now().Format("2006-01-02")
	err := bot.Storage.DB.Insert("quotes", bot.Storage.QuoteColumns, []string{quote, date, submitter})
	bot.LoadQuotes()
	return err
}

// given the id, generates a string containing the quote, timestamp and submitter
func (bot *Bot) generateQuoteString(id int) string {
	values := bot.Quotes[id]
	return fmt.Sprintf("%s -- %s [submitted by %s]", values.Quote, values.Timestamp, values.Submitter)
}

// RandomQuote returns a random string, it does not print. It's up to the caller on what to do with it.
func (bot *Bot) RandomQuote() (string, error) {
	if len(bot.Quotes) == 0 {
		return "", errors.New("no quotes were found")
	}
	randomIndex := rand.Intn(len(bot.Quotes))

	return bot.generateQuoteString(randomIndex + 1), nil // return quote string
}

// LoadQuotes loads all quotes from the DB.
// TODO: generalize this for all loading functions
func (bot *Bot) LoadQuotes() error {
	rows, err := bot.Storage.DB.Query("select id, quote, timestamp, submitter from quotes")
	if err != nil {
		return err
	}

	defer rows.Close()
	for rows.Next() { // scan through results from query and assign to the Commands slice
		var id int
		var quote, timestamp, submitter string
		err = rows.Scan(&id, &quote, &timestamp, &submitter)
		if err != nil {
			return err
		}
		bot.Quotes[id] = &QuoteValues{Quote: quote, Timestamp: timestamp, Submitter: submitter}
	}
	return nil
}

func (bot *Bot) DeleteQuote(quoteID string) error {
	id, err := strconv.Atoi(quoteID)
	if err != nil {
		return err
	}

	if _, found := bot.Quotes[id]; found {
		delete(bot.Quotes, id)
	}

	return bot.Storage.DB.Delete("quotes", "id", quoteID)
}

// GetQuote returns a quote of a specified index / id. Correlates to the automatically generated ID in sqlite.
func (bot *Bot) GetQuote(index int) (string, error) {
	if index <= 0 {
		return "", errors.New("the ID must be a valid integer greater than 0")
	} else if index > len(bot.Quotes) {
		return "", fmt.Errorf("the requested ID %d is greater than the total number of quotes, which is: %d", index, len(bot.Quotes))
	}
	return bot.generateQuoteString(index), nil // return quote string
}
