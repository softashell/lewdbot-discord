package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/softashell/lewdbot-discord/brain"
	"github.com/softashell/lewdbot-discord/commands"
	"github.com/softashell/lewdbot-discord/config"
	"github.com/softashell/lewdbot-discord/lewd"
	"github.com/softashell/lewdbot-discord/regex"
)

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID || m.Author.Bot || len(m.Message.Content) < 1 {
		return
	}

	channel, err := s.State.Channel(m.ChannelID)
	if err != nil {
		log.Warn("s.State.Channel >> ", err)
		channel, err = s.Channel(m.ChannelID)
		if err != nil {
			log.Warn("s.Channel >> ", err)
			return
		}
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
		_, err := s.ChannelMessageSendReply(m.ChannelID, reply, m.Reference())
		if err != nil {
			log.Warn("s.ChannelMessageSendReply >> ", err)
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

	if config.ChannelIsMangadex(channel.GuildID, m.ChannelID) {
		if mangadexClient.ParseLinks(s, m, m.ChannelID, text) {
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

		user, err := s.State.Member(m.GuildID, m.Author.ID)
		log.Info(user.Nick)
		if err != nil {
			user, err = s.GuildMember(m.GuildID, m.Author.ID, discordgo.WithRestRetries(1))
			if err != nil {
				// We fucked up??
			}
		}

		username := m.Author.GlobalName // Display name
		if user != nil && user.Nick != "" {
			username = user.Nick
		}

		if username == "" {
			username = m.Author.Username
		}

		reply = regex.Lewdbot.ReplaceAllString(reply, username)

		// Log our reply
		fmt.Printf("%20s %20s %20s > %s\n", channel.Name, time.Now().Format(time.Stamp), s.State.User.Username, reply)

		// Fetch latest channel
		channel, err = s.Channel(m.ChannelID)
		if err != nil {
			log.Warn("s.Channel >> ", err)
		}

		if channel == nil || channel.LastMessageID != m.ID {
			_, err = s.ChannelMessageSendReply(m.ChannelID, reply, m.Reference())
			if err != nil {
				log.Warn("s.ChannelMessageSendReply >> ", err)
			}
		} else {
			_, err = s.ChannelMessageSend(m.ChannelID, reply)
			if err != nil {
				log.Warn("s.ChannelMessageSend >> ", err)
			}
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
		log.Errorf("presenceUpdate :: failed to get guild - %s", err)
		return
	}

	role, err := getRole(guild, "Streamer")
	if err != nil {
		log.Error(err)
		return
	}

	updateStreamerRole(s, &m.Presence, guild.ID, m.User.ID, role.ID)
}

func guildMembersChunk(s *discordgo.Session, g *discordgo.GuildMembersChunk) {
	log.Infof("GuildMembersChunk: Adding %d members to %s", len(g.Members), g.GuildID)
}

func ready(s *discordgo.Session, r *discordgo.Ready) {
	log.Info("websocket READY")

	for _, g := range r.Guilds {
		log.Info("requesting members ", g.ID)
		err := s.RequestGuildMembers(g.ID, "", 0, "benis", true)
		if err != nil {
			log.Errorf("failed to get members")
			continue
		}
	}

	time.Sleep(30 * time.Second)
	log.Info("fixing roles")

	fixStreamerRoles(s)
}
