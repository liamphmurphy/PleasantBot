// this file will contain some helper funcs for preparing the bot's config.

package bot

import "github.com/spf13/viper"

// writeConfig is run whenever the config.toml file doesn't exist, usually after a fresh download of the bot.
func writeConfig(path string, configObject *viper.Viper) {
	path = path + "/config.toml"
	// prepare default values, will be used when viper writes the new config file
	configObject.SetDefault("ChannelName", "<enter channel name to moderate here>")
	configObject.SetDefault("ServerName", "irc.chat.twitch.tv:6697")
	configObject.SetDefault("BotName", "<enter bot username here>")
	configObject.SetDefault("BotOAuth", "<bot oauth>")
	configObject.SetDefault("PurgeForLinks", true)
	configObject.SetDefault("PurgeForLongMsg", true)
	configObject.SetDefault("LongMsgAmount", 400)
	configObject.SetDefault("EnableServer", true)
	configObject.SetDefault("PostLinkPerm", uint(1)) // Minimum permission needed for non-purging links, in this case subscriber

	configObject.WriteConfigAs(path)
}
