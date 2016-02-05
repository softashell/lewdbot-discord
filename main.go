package main

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/softashell/lewdbot-discord/regex"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var bots = [...]string{
	"141615500462522368", // DEE JAY-chang
	"142359333227724800", // SCI FI-chang
}

var (
	chat   *Chat
	client *http.Client
)

func main() {
	chat = NewChat()

	chat.learnFileLines("./data/brain.txt", true)
	chat.learnFileLines("./data/dump.txt", true)
	chat.learnFileLines("./data/chatlog.txt", false)

	email, pw, id, hash := LoadConfigFromFile("config.json")

	client = create_client(id, hash)

	d, err := discordgo.New(email, pw)
	if err != nil {
		fmt.Println(err)
		return
	}

	d.OnMessageCreate = messageCreate
	d.ShouldReconnectOnError = true

	// Open the websocket and begin listening.
	d.Open()

	// Simple way to keep program running until any key press.
	var input string
	fmt.Scanln(&input)
	return
}

func LoadConfigFromFile(filename string) (string, string, string, string) {
	fileDump, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err.Error())
	}

	type fileCredentials struct {
		Email         string `json:"email"`
		Password      string `json:"password"`
		Exhtentai_id  string `json:"exhentai_id"`
		Exhentai_hash string `json:"exhentai_hash"`
	}

	var creds = fileCredentials{}
	if err := json.Unmarshal(fileDump, &creds); err != nil {
		fmt.Println(err.Error())
	}

	return creds.Email, creds.Password, creds.Exhtentai_id, creds.Exhentai_hash
}

func messageCreate(s *discordgo.Session, m *discordgo.Message) {
	text := m.Content
	isMentioned := false

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

	for _, id := range bots {
		if id == m.Author.ID {
			// Fucking bot spam smh
			return
		}
	}

	// Replace internal mention strings with actual name
	for _, mention := range m.Mentions {
		if mention.ID == s.State.User.ID {
			isMentioned = true
		}

		mention_text := "<@" + mention.ID + ">"
		text = strings.Replace(text, mention_text, mention.Username, -1)
	}

	text = strings.Replace(text, "@everyone", "", -1)

	// Log cleaned up message
	fmt.Printf("%20s %20s %20s > %s\n", channel.Name, time.Now().Format(time.Stamp), m.Author.Username, text)

	links_found, reply := parse_links(text)

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
		reply := chat.generateReply(text)
		reply = regex.Lewdbot.ReplaceAllString(reply, m.Author.Username)

		// Log our reply
		fmt.Printf("%20s %20s %20s > %s\n", channel.Name, time.Now().Format(time.Stamp), s.State.User.Username, reply)

		s.ChannelMessageSend(m.ChannelID, reply)
	} else {
		// Just learn
		chat.learnMessage(text, true)
	}
}