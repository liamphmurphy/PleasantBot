// this file will define custom errors used throughout the bot

package bot

import "fmt"

// FatalError is an error that would warrant an immediate end to the bot
type FatalError struct {
	Err error
}

// NonFatalError explains that an error occurred due to a non fatal condition and shouldn't be considered
// a terrible issue.
type NonFatalError struct {
	Err error
}

func (fe FatalError) Error() string {
	return fmt.Sprintf("a fatal error occurred: %v", fe.Err)
}

func (nfe NonFatalError) Error() string {
	return fmt.Sprintf("a non-fatal error occurred: %v", nfe.Err)
}
