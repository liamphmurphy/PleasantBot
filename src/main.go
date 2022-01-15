// This file acts as a driver for the bot package. Theoretically this repo could be forked, and an entirely new interface could be written by replacing this file.

package main

import (
	"bufio"
	"fmt"
	"log"
	"net/textproto"
	"os"
	"regexp"
	"time"

	"github.com/murnux/pleasantbot/storage"

	"github.com/murnux/pleasantbot/bot"
)

func main() {
	run()
}

func run() {
	configPath, err := os.UserConfigDir()
	if err != nil {
		log.Fatalf(err.Error())
	}
	configPath = fmt.Sprintf("%s/pleasantbot", configPath)

	var db storage.Sqlite
	pleasant, err := bot.CreateBot(configPath, "config", "pleasantbot.db", true, &db)
	if err != nil {
		log.Fatalf(err.Error())
	}

	if pleasant.EnableServer {
		go pleasant.StartAPI()
	}

	for {
		if !pleasant.Authenticated {
			time.Sleep(1 * time.Second)
		} else {
			break
		}
	}
	err = pleasant.Connect()
	if err != nil {
		log.Fatalf(err.Error())
	}

	// Pass info to HTTP request
	fmt.Fprintf(pleasant.Conn, "PASS %s\r\n", pleasant.GetOAuth())
	fmt.Fprintf(pleasant.Conn, "NICK %s\r\n", pleasant.Name)
	fmt.Fprintf(pleasant.Conn, "JOIN #%s\r\n", pleasant.ChannelName)

	// Twitch specific information, like badges, mod status etc.
	fmt.Fprintf(pleasant.Conn, "CAP REQ :twitch.tv/membership\r\n")
	fmt.Fprintf(pleasant.Conn, "CAP REQ :twitch.tv/tags\r\n")
	fmt.Fprintf(pleasant.Conn, "CAP REQ :twitch.tv/commands\r\n")

	defer pleasant.Conn.Close()

	reader := bufio.NewReader(pleasant.Conn) // prepare network line readers
	proto := textproto.NewReader(reader)

	msgRegex, _ := regexp.Compile("[;]+") // regexp object used to split messages
	pingIndicator := "PING :tmi.twitch.tv"

	fmt.Printf("Bot: %s\nChannel: %s\n", pleasant.Name, pleasant.ChannelName)

	pleasant.RunTimers()

	// keep reading messages until some end condition is reached
	for {
		line, err := proto.ReadLine()
		if err != nil {
			fmt.Printf("error receiving message: %s\n", err)
			os.Exit(1)
		}
		fmt.Println(line)
		ponged, err := pleasant.HandlePing(line, pingIndicator)
		if err != nil {
			log.Fatalf(err.Error())
		}

		// skip further execution if we PONGed
		if ponged {
			fmt.Printf("sent a PONG back to %s", pleasant.ServerName)
			continue
		}

		lineSplit := msgRegex.Split(line, -1) // TODO: apparently regex splitting isn't efficient, need to research this

		if len(lineSplit) <= 13 { // at this point the message should be from chat, so confirm the length (this is just an approx)
			continue
		}

		message := pleasant.ParseMessage(lineSplit) // create readable message from the user
		pleasant.FilterForSpam(message)
		fmt.Printf("%s: %s\n", message.Name, message.Content)

		if message.IsCommand { // if first character in a chat message is !, it's probably a command
			if pleasant.DefaultCommands(message) { // see if message is a default command request
				continue // match is found and the bot took action, move on
			}
			com, err := pleasant.FindCommand(message.Content)
			if err != nil {
				fmt.Printf("error finding command: %s\n", err)
				continue
			}

			comPerm, err := pleasant.ConvertPermToInt(com.Perm) // see if it is a custom command request
			if message.Perm >= comPerm && err == nil {          // send message if user has permission and there were no errors finding the command
				pleasant.SendTwitchMessage(com.Response)
				go pleasant.IncrementCommandCount(message.Content) // increment command count
			}
		}
	}
}
