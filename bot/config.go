// this file will contain some helper funcs for preparing the bot's config.

package bot

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// CreateViperConfig creates a viper object. The path is the directory where the config should reside, and name is the
// filename.
func CreateViperConfig(path, name, configType, serverName string) (*viper.Viper, error) {
	var v *viper.Viper
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0755) // create the config directory if needed
	}

	v = viper.New()
	v.SetConfigName(name)
	v.SetConfigType(configType)
	v.AddConfigPath(path)

	fullPath := fmt.Sprintf("%s/%s", path, name)
	// attempt to read in the config
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// write a default config file at path/name.
			writeConfig(fmt.Sprintf("%s/%s", path, name), serverName, v)
			return nil, NonFatalError{Err: fmt.Errorf("had to create a default config file, please go to %s and edit values as needed", fullPath)}
		} else {
			return nil, FatalError{Err: err}
		}
	}

	return v, nil
}

// GetConfigDirectory the "default" location to put in pleasantbot's config files, so Runners can use this function if desired.
func GetConfigDirectory() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return home, FatalError{err}
	}
	return fmt.Sprintf("%s/.config/pleasantbot", home), nil
}

// writeConfig is run whenever the config.toml file doesn't exist, usually after a fresh download of the bot.
func writeConfig(path, serverName string, configObject *viper.Viper) {
	// prepare default values, will be used when viper writes the new config file
	configObject.SetDefault("ChannelName", "<enter channel name to moderate here>")
	configObject.SetDefault("ServerName", serverName)
	configObject.SetDefault("BotName", "<enter bot username here>")
	configObject.SetDefault("BotOAuth", "<bot oauth>")
	configObject.SetDefault("PurgeForLinks", true)
	configObject.SetDefault("PurgeForLongMsg", true)
	configObject.SetDefault("LongMsgAmount", 400)
	configObject.SetDefault("EnableServer", true)
	configObject.SetDefault("PostLinkPerm", uint(1)) // Minimum permission needed for non-purging links, in this case subscriber

	configObject.WriteConfigAs(path)
}
