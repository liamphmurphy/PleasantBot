// This file acts as a driver for the bot package. Theoretically this repo could be forked, and an entirely new interface could be written by replacing this file.

package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"net/textproto"
	"os"
	"regexp"

	"github.com/murnux/pleasantbot/bot"
)

func main() {
	pleasant := bot.CreateBot()
	var err error
	pleasant.DB, err = sql.Open("sqlite3", pleasant.DBPath)
	if err != nil {
		log.Fatalf("failed to open database in main function: %s", err)
	}
	defer pleasant.DB.Close()
	connErr := pleasant.Connect()
	if connErr != nil {
		fmt.Printf("error connecting to irc.twitch.tv: %s\n", connErr)
		os.Exit(1)
	}

	// Pass info to HTTP request
	fmt.Fprintf(pleasant.Conn, "PASS %s\r\n", pleasant.OAuth)
	fmt.Fprintf(pleasant.Conn, "NICK %s\r\n", pleasant.Name)
	fmt.Fprintf(pleasant.Conn, "JOIN #%s\r\n", pleasant.ChannelName)

	// Twitch specific information, like badges, mod status etc.
	fmt.Fprintf(pleasant.Conn, "CAP REQ :twitch.tv/membership\r\n")
	fmt.Fprintf(pleasant.Conn, "CAP REQ :twitch.tv/tags\r\n")
	fmt.Fprintf(pleasant.Conn, "CAP REQ :twitch.tv/commands\r\n")

	defer pleasant.Conn.Close()

	reader := bufio.NewReader(pleasant.Conn) // prepare net line readers
	proto := textproto.NewReader(reader)

	msgRegex, _ := regexp.Compile("[;]+") // regexp object used to split messages

	fmt.Printf("Bot: %s\nChannel: %s\n", pleasant.Name, pleasant.ChannelName)

	// keep reading messages until some end condition is reached
	for {
		line, err := proto.ReadLine()
		if err != nil {
			fmt.Printf("error receiving message: %s\n", err)
			os.Exit(1)
		}

		lineSplit := msgRegex.Split(line, -1) // TODO: apparently regex splitting isn't efficient, need to research this
		if lineSplit[0] == "PING" {           // anticipate PING message
			pleasant.WriteToTwitch("PONG :tmi.twitch.tv")
			log.Println("INFO -- replied to PING with a PONG")
		}

		if len(lineSplit) <= 13 { // at this point the message should be from chat, so confirm the length (this is just an approx)
			continue
		}

		message := pleasant.ParseMessage(lineSplit) // create readable message from the user
		fmt.Printf("%s: %s\n", message.Name, message.Content)

		if message.IsCommand { // if first character in a chat message is ! (unicoode value is 64), it's probably a command
			if pleasant.DefaultCommands(message.Content) { // see if message is a default command request
				continue // match is found and the bot took action, move on
			}
			com, err := pleasant.FindCommand(message.Content) // see if it is a custom command request
			if err == nil {
				pleasant.SendMessage(com.Response)
			}
		}
	}
}
