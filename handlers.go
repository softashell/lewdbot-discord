package main

import (
	"fmt"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/bwmarrin/discordgo"
	"github.com/softashell/lewdbot-discord/brain"
	"github.com/softashell/lewdbot-discord/commands"
	"github.com/softashell/lewdbot-discord/config"
	"github.com/softashell/lewdbot-discord/lewd"
	"github.com/softashell/lewdbot-discord/regex"
)

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID || len(m.Message.Content) < 1 {
		return
	}

	channel, err := s.State.Channel(m.ChannelID)
	if err != nil {
		log.Warn("s.State.Channel >> ", err)
	}

	if channel.Type == discordgo.ChannelTypeDM {
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
			log.Warn("s.ChannelMessageSend >> ", err)
		}

		return
	} else if strings.HasPrefix(text, "!") || strings.HasPrefix(text, ".") || strings.HasPrefix(text, "bot.") {
		// Ignore shit meant for other bots
		return
	}

	if config.ChannelIsLewd(channel.GuildID, m.ChannelID) {
		if lewd.ParseLinks(s, m.ChannelID, text) {
			return
		}
	}

	// Accept the legacy mention as well and trim it off from text
	if strings.HasPrefix(strings.ToLower(text), "lewdbot, ") {
		text = text[9:]
		isMentioned = true
	}

	if channel.Type == discordgo.ChannelTypeDM || isMentioned || config.ChannelShouldSpam(channel.GuildID, m.ChannelID) {
		err := s.ChannelTyping(m.ChannelID)
		if err != nil {
			log.Warn("s.ChannelTyping >> ", err)
		}

		reply := brain.Reply(text)
		reply = regex.Lewdbot.ReplaceAllString(reply, m.Author.Username)

		// Log our reply
		fmt.Printf("%20s %20s %20s > %s\n", channel.Name, time.Now().Format(time.Stamp), s.State.User.Username, reply)

		_, err = s.ChannelMessageSend(m.ChannelID, reply)
		if err != nil {
			log.Warn("s.ChannelMessageSend >> ", err)
		}

	} else if !config.GuildIsDumb(channel.GuildID) {
		// Just learn
		brain.Learn(text, true)
	}
}

func presenceUpdate(s *discordgo.Session, m *discordgo.PresenceUpdate) {
	if !config.GuildHasStreamerRoleEnabled(m.GuildID) {
		return
	}

	guild, err := s.State.Guild(m.GuildID)
	if err != nil {
		log.Error("presenceUpdate :: failed to get guild - %s", err)
		return
	}

	role, err := getRole(guild, "Streamer")
	if err != nil {
		log.Error(err)
		return
	}

	roleAdded, err := hasRole(s, m.GuildID, m.User.ID, role.ID)
	if err != nil {
		log.Error("presenceUpdate :: failed to get member roles - %s", err)
		return
	}

	if m.Presence.Game == nil && roleAdded {
		log.Infof("presenceUpdate :: removing  streamer group from %s (%s || %s)", m.User.ID, m.User.Username, m.Nick)
		err = s.GuildMemberRoleRemove(m.GuildID, m.User.ID, role.ID)
		if err != nil {
			log.Errorf("presenceUpdate :: failed to remove streamer role - %s", err)
			return
		}
	} else if m.Presence.Game != nil && m.Presence.Game.Type == discordgo.GameTypeStreaming && !roleAdded {
		log.Infof("presenceUpdate :: adding streamer group from %s (%s || %s)", m.User.ID, m.User.Username, m.Nick)
		err = s.GuildMemberRoleAdd(m.GuildID, m.User.ID, role.ID)
		if err != nil {
			log.Errorf("presenceUpdate :: failed to add streamer role - %s", err)
			return
		}
	}
}
