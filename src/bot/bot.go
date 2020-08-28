// This file handles some of the core functionality of the bot. Such as creation of the bot and connecting to Twitch.

package bot

import (
	"fmt"
	"net"
	"os"

	"github.com/spf13/viper"
)

// Bot struct contains the necessary data to run an instance of a bot
type Bot struct {
	ChannelName string
	ServerName  string
	BotOAuth    string
	BotName     string
	Conn        net.Conn
	Commands    []Command
	BadWords    []BadWord
}

// writeConfig is run whenever the config.toml file doesn't exist, usually after a fresh download of the bot.
func writeConfig(path string) {
	path = path + "/config.toml"
	// prepare default values, will be used when viper writes the new config file
	viper.SetDefault("ChannelName", "<enter channel name to moderate here>")
	viper.SetDefault("ServerName", "irc.twitch.tv:6667")
	viper.SetDefault("BotName", "<enter bot username here>")
	viper.SetDefault("BotOAuth", "<bot oauth>")

	viper.WriteConfigAs(path)
	fmt.Println(fmt.Sprintf("Config file did not exist, so it has been made. Please go to %s and edit the settings.", path))
}

// CreateBot creates an instance of a bot
func CreateBot() *Bot {
	configDir, _ := os.UserConfigDir()        // follows the standard config dir used by the OS
	pleasantDir := configDir + "/pleasantbot" // config.toml will go here
	if _, err := os.Stat(pleasantDir); os.IsNotExist(err) {
		os.Mkdir(pleasantDir, 0755)
	}

	// prepare toml config file
	viper.SetConfigName("config")
	viper.AddConfigPath(pleasantDir)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok { // if config not found, write a default one
			writeConfig(pleasantDir)
			os.Exit(1) // write a better way to exit this gracefully
		} else {
			panic(fmt.Errorf("error reading in config file: %s", err))
		}
	}

	var commands []Command
	var badwords []BadWord
	return &Bot{ // create Bot instance
		ChannelName: viper.GetString("ChannelName"),
		ServerName:  viper.GetString("ServerName"),
		BotOAuth:    viper.GetString("BotOAuth"),
		BotName:     viper.GetString("BotName"),
		Commands:    commands,
		BadWords:    badwords,
	}
}

// Connect establishes a connection to the Twitch IRC server
func (bot *Bot) Connect() error {
	var err error
	if bot.Conn != nil {
		return fmt.Errorf("ERROR: the bot has already established an IRC connection")
	}
	bot.Conn, err = net.Dial("tcp", bot.ServerName)
	return err
}

// ChannelConnect writes the necessary scopes to Twitch
func (bot *Bot) ChannelConnect() {
	// Pass info to HTTP request
	fmt.Fprintf(bot.Conn, "PASS %s\r\n", bot.BotOAuth)
	fmt.Fprintf(bot.Conn, "NICK %s\r\n", bot.BotName)
	fmt.Fprintf(bot.Conn, "JOIN #%s\r\n", bot.ChannelName)

	// Twitch specific information, like badges, mod status etc.
	fmt.Fprintf(bot.Conn, "CAP REQ :twitch.tv/membership\r\n")
	fmt.Fprintf(bot.Conn, "CAP REQ :twitch.tv/tags\r\n")
	fmt.Fprintf(bot.Conn, "CAP REQ :twitch.tv/commands\r\n")

	defer bot.Conn.Close()
}

// WriteToTwitch when given a string sends a properly formatted message to Twitch, easy replacement for using fmt.Fprintf
func (bot *Bot) WriteToTwitch(msg string) {
	fmt.Fprintf(bot.Conn, "%s\r\n", msg)
}

// SendMessage prepares and sends a string to the channel's Twitch chat
func (bot *Bot) SendMessage(msg string) {
	fullMessage := fmt.Sprintf("PRIVMSG #%s :%s", bot.ChannelName, msg)
	bot.WriteToTwitch(fullMessage)
}
