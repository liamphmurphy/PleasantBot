# PleasantBot

PleasantBot is planned to be a full-stack application to help streamers moderate their twitch streams with a helpful bot. It should have the features that many streamers expect, such as:

- Commands
- Quotes
- Ban / purge users for using bad language
- Misc. moderation for links, long messages etc. 

## Running

To run the bot as of now, run the following command in the /src directory:
`go build && ./pleasantbot`

# Goals

- Respect the OS's default settings for config locations using golang's os module.
- Maintain a small database of user added commands, bannable words and more with SQLite.
- Create API middle-man between backend and front-end UI
- Add support for SSL/TLS
- Make it work in Docker

# HappyBot

If you're curious, I worked on another bot on the exact same stack, but I was pretty new to programming and was not far in my formal studies. It's kind of trash, but it did work (except I never finished the GUI D:)  https://gitlab.com/murnux/HappyBot