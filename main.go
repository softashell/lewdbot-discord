package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/softashell/lewdbot-discord/brain"
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

	log.Println("Filling brain")

	brain.LearnFileLines("./data/brain.txt", true)
	brain.LearnFileLines("./data/dump.txt", true)
	brain.LearnFileLines("./data/chatlog.txt", false)

	connectToDiscord()

	// Simple way to keep program running until any key press.
	var input string
	fmt.Scanln(&input)
}

func connectToDiscord() {
	log.Println("Connecting to discord")

	var err error

	c := config.Get()

	dg, err := discordgo.New(c.Email, c.Password)
	if err != nil {
		log.Println(err.Error())
		return
	}

	// Register messageCreate as a callback for the OnMessageCreate event.
	dg.AddHandler(messageCreate)

	// Retry after broken websocket
	dg.ShouldReconnectOnError = true

	// Open websocket connection
	err = dg.Open()
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Connected")
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	text := m.ContentWithMentionsReplaced()

	if m.Author.ID == s.State.User.ID {
		// Ignore self
		return
	}

	channel, _ := s.State.Channel(m.ChannelID)

	if channel.IsPrivate {
		channel.Name = "direct message"
	}

	if strings.HasPrefix(text, "!") || strings.HasPrefix(text, ".") || strings.HasPrefix(text, "bot.") {
		// Ignore shit meant for other bots
		return
	}

	isMentioned := isUserMentioned(s.State.User, m.Mentions) || m.MentionEveryone

	text = strings.Replace(text, "@everyone", "", -1)

	// Log cleaned up message
	fmt.Printf("%20s %20s %20s > %s\n", channel.Name, time.Now().Format(time.Stamp), m.Author.Username, text)

	if shouldIgnore(m.Author) {
		return
	}

	links_found, reply := lewd.ParseLinks(text)

	if links_found {
		s.ChannelMessageSend(m.ChannelID, reply)
		return
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
	} else {
		// Just learn
		brain.Learn(text, true)
	}
}

func shouldIgnore(user *discordgo.User) bool {
	c := config.Get()

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
