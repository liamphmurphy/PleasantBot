package bot

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// contains one row of data from the ban_history table
type banHistory struct {
	User      string
	Reason    string
	Timestamp string
}

// used for the "Quick Stats" section fo the dashboard
type stats struct {
	Commands     int
	Quotes       int
	Bans         int
	TopCommand   string
	TopComCount  int
	TopChatter   string
	TopChatCount int
}

// struct to store data from a addcom POST request
type commandPost struct {
	Name string `json:"CommandName"`
	Resp string `json:"Response"`
	Perm string `json:"Perm"`
}

// return JSON representation of the commands
func (bot *Bot) getComHandler(c *gin.Context) {
	c.JSON(http.StatusOK, bot.Commands)
}

// return JSON represntation of bot config data
func (bot *Bot) getBotData(c *gin.Context) {
	c.JSON(http.StatusOK, bot)
}

func (bot *Bot) getQuotes(c *gin.Context) {
	c.JSON(http.StatusOK, bot.Quotes)
}

// return JSON representation of ban_history events
func (bot *Bot) getBanHistory(c *gin.Context) {
	rows, _ := bot.DB.Query("select user, reason, timestamp from ban_history ORDER BY timestamp DESC;")
	var history []banHistory

	defer rows.Close()
	for rows.Next() {
		var user, reason, timestamp string
		rows.Scan(&user, &reason, &timestamp)
		history = append(history, banHistory{User: user, Reason: reason, Timestamp: timestamp})
	}

	if len(history) == 0 {
		c.JSON(http.StatusBadRequest, "There is no ban history data.")
		return
	}

	c.JSON(http.StatusOK, history) // serve the ban_history slice
}

func (bot *Bot) checkAuthHandler(c *gin.Context) {
	c.JSON(200, bot.Authenticated)
}

// return JSON representation of quick stats
/*func (bot *Bot) getStats(c *gin.Context) {
	topCom, count := bot.GetTopFromTable("commands", "commandname", "count")
	topChat, chatCount := bot.GetTopFromTable("chatters", "username", "count")
	c.JSON(http.StatusOK, stats{Commands: len(bot.Commands), Quotes: len(bot.Quotes), Bans: 0, TopCommand: topCom, TopComCount: count, TopChatter: topChat, TopChatCount: chatCount})
}*/

func (bot *Bot) addComHandler(c *gin.Context) {
	var comValues commandPost

	err := c.BindJSON(&comValues)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Status": "Error decoding command JSON"})
	}

	// First add new command to database
	err = bot.DB.Insert("commands", bot.CommandDBColumns, []string{comValues.Name, comValues.Resp, comValues.Perm, "0"})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Status": "Error adding command to DB"})
		log.Fatal(err)
	}

	// Now that it is in the DB, add new command to bot.Commands map
	bot.Commands[comValues.Name] = &CommandValue{Response: comValues.Resp, Perm: comValues.Perm, Count: 0}
	c.JSON(http.StatusOK, gin.H{"Status": "Command successfully added"})
}

func (bot *Bot) delComHandler(c *gin.Context) {
	var commandNames []string

	err := c.BindJSON(&commandNames)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Status": "Error binding JSON in delcom handler"})
	}

	deleteFailed := false
	for _, com := range commandNames {
		if !bot.RemoveCommand(com) {
			deleteFailed = true
		}
	}

	if deleteFailed {
		c.JSON(http.StatusBadRequest, gin.H{"Status": "Couldn't delete all commands"})
	} else {
		c.JSON(http.StatusOK, gin.H{"Status": "Commands deleted"})
	}
}

func (bot *Bot) addOAuthHandler(c *gin.Context) {
	var token string
	err := c.BindJSON(&token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Status:": "Error receiving new oauth token."})
	}

	bot.SetOAuth(token)
	c.JSON(http.StatusOK, gin.H{"Status": "New oauth set."})
}

// StartAPI starts the gin router for the bot's API
func (bot *Bot) StartAPI() {
	router := gin.Default()

	router.Use(cors.Default())

	// API GET endpoints
	router.GET("/getcoms", bot.getComHandler)
	router.GET("/getbot", bot.getBotData)
	router.GET("/getbanhistory", bot.getBanHistory)
	//router.GET("/getstats", bot.getStats)
	router.GET("/getquotes", bot.getQuotes)
	router.GET("/checkauth", bot.checkAuthHandler)

	// API POST endpoints
	router.POST("/addcom", bot.addComHandler)
	router.POST("/delcom", bot.delComHandler)
	router.POST("/addoauth", bot.addOAuthHandler)

	router.Run()
}
