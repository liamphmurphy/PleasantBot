package bot

import "github.com/gin-gonic/gin"

// return JSON representation of the commands
func (bot *Bot) getComHandler(c *gin.Context) {
	c.JSON(200, bot.Commands)
}

// return JSON represntation of bot config data
func (bot *Bot) getBotData(c *gin.Context) {
	c.JSON(200, bot)
}

// StartAPI starts the gin router for the bot's API
func (bot *Bot) StartAPI() {
	router := gin.Default()
	router.GET("/getcoms", bot.getComHandler)
	router.GET("/getbot", bot.getBotData)
	router.Run()
}
