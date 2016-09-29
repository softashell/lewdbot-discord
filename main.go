package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/softashell/lewdbot-discord/brain"
	"github.com/softashell/lewdbot-discord/commands"
	"github.com/softashell/lewdbot-discord/config"
	"github.com/softashell/lewdbot-discord/lewd"
	"github.com/softashell/lewdbot-discord/regex"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	os.Mkdir("./data", 0744)

	config.Init()
	brain.Init()

	c := config.Get()

	go func() {
		start := time.Now()

		log.Println("Starting to fill brain")

		for _, b := range c.Brain {
			log.Println("Parsing", b.File)

			err := brain.LearnFileLines(b.File, b.Simple)

			if err != nil {
				log.Println(err)
				return
			}
		}

		log.Println("Parsing ./data/chatlog.txt")
		err := brain.LearnFileLines("./data/chatlog.txt", false)
		if err != nil {
			log.Println(err)
		}

		log.Println("Brain filled in", time.Since(start))
	}()

	connectToDiscord()

	// Simple way to keep program running until any key press.
	var input string
	fmt.Scanln(&input)
}

func connectToDiscord() {
	log.Println("Connecting to discord")

	var err error

	c := config.Get()

	dg, err := discordgo.New(c.Email, c.Password, "Bot "+c.Token)
	if err != nil {
		log.Println(err.Error())
		return
	}

	// Register messageCreate as a callback for the OnMessageCreate event.
	dg.AddHandler(messageCreate)

	// Retry after broken websocket
	dg.ShouldReconnectOnError = true

	// Verify the Token is valid and grab user information
	dg.State.User, err = dg.User("@me")
	if err != nil {
		log.Printf("error fetching user information, %s\n", err)
	}

	// Open websocket connection
	err = dg.Open()
	if err != nil {
		log.Printf("error opening connection to Discord, %s\n", err)
		return
	}

	log.Println("Connected")
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		// Ignore self
		return
	}

	channel, _ := s.State.Channel(m.ChannelID)

	if channel.IsPrivate {
		channel.Name = "direct message"
	}

	isMentioned := isUserMentioned(s.State.User, m.Mentions) || m.MentionEveryone

	if shouldIgnore(m.Author) {
		return
	}

	text := m.ContentWithMentionsReplaced()
	text = strings.Replace(text, "@everyone", "", -1)

	// Log cleaned up message
	fmt.Printf("%20s %20s %20s > %s\n", channel.Name, time.Now().Format(time.Stamp), m.Author.Username, text)

	commandFound, reply := commands.ParseMessage(s, m, text)

	if commandFound {
		_, err := s.ChannelMessageSend(m.ChannelID, reply)
		if err != nil {
			fmt.Println("s.ChannelMessageSend >> ", err)
		}
		return
	} else if strings.HasPrefix(text, "!") || strings.HasPrefix(text, ".") || strings.HasPrefix(text, "bot.") {
		// Ignore shit meant for other bots
		return
	}

	if config.ChannelIsLewd(channel.GuildID, m.ChannelID) {
		linksFound, reply := lewd.ParseLinks(text)

		if linksFound {
			s.ChannelMessageSend(m.ChannelID, reply)
			return
		}
	}

	// Accept the legacy mention as well and trim it off from text
	if strings.HasPrefix(strings.ToLower(text), "lewdbot, ") {
		text = text[9:]
		isMentioned = true
	}

	if channel.IsPrivate || isMentioned {
		reply := brain.Reply(text)
		reply = regex.Lewdbot.ReplaceAllString(reply, m.Author.Username)

		// Log our reply
		fmt.Printf("%20s %20s %20s > %s\n", channel.Name, time.Now().Format(time.Stamp), s.State.User.Username, reply)

		s.ChannelMessageSend(m.ChannelID, reply)
	} else if !config.GuildIsDumb(channel.GuildID) {
		// Just learn
		brain.Learn(text, true)
	}
}

func shouldIgnore(user *discordgo.User) bool {
	c := config.Get()

	if user.Bot {
		return true
	}

	for _, id := range c.Blacklist {
		if id == user.ID {
			return true
		}
	}

	return false
}

func isUserMentioned(user *discordgo.User, mentions []*discordgo.User) bool {
	for _, mention := range mentions {
		if mention.ID == user.ID {
			return true
		}
	}

	return false
}
