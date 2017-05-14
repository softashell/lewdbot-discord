package commands

import (
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/bwmarrin/discordgo"
	"github.com/softashell/lewdbot-discord/config"
)

const (
	msgOn  = ":ok_hand:"
	msgOff = ":clap:"
)

func ParseMessage(s *discordgo.Session, m *discordgo.MessageCreate, text string) (bool, string) {
	var reply string

	channel, err := s.State.Channel(m.ChannelID)
	if err != nil {
		log.Errorf("s.Channel(%s) >> %s\n", m.ChannelID, err)
		return false, reply
	}

	command := strings.ToLower(text)

	if strings.HasPrefix(command, "!set") {
		if !config.IsMaster(m.Author.ID) {
			reply = "Who the **HELL** are you, dude??"
			return true, reply
		}

		if len(command) < 6 {
			return false, ""
		}

		command = command[5:]

		if strings.HasPrefix(command, "lewd") {
			if config.ChannelSetLewd(channel.GuildID, channel.ID) {
				reply = msgOn
			} else {
				reply = msgOff
			}

			return true, reply
		} else if strings.HasPrefix(command, "dumb") {
			if config.GuildSetDumb(channel.GuildID) {
				reply = msgOn
			} else {
				reply = msgOff
			}

			return true, reply

		} else if strings.HasPrefix(command, "roles") {
			if config.SetManageRoles(channel.GuildID) {
				reply = msgOn
			} else {
				reply = msgOff
			}

			return true, reply
		}

		return false, ""
	} else if strings.HasPrefix(command, "!8ball") {
		reply = eightball(text)

		return true, reply
	} else if strings.HasPrefix(command, "!roll") {
		reply = dice(text[5:], m.Author)

		return true, reply
	} else if config.ShouldManageRoles(channel.GuildID) {
		if strings.HasPrefix(command, "!list") {
			if len(text) > 6 {
				reply = listRoleMembers(s, channel.GuildID, text[6:])
			} else {
				reply = listRoles(s, channel.GuildID)
			}
			return true, reply
		} else if strings.HasPrefix(command, "!subscribe") {
			if len(text) > 11 {
				reply = addRole(s, channel.GuildID, m.Author.ID, text[11:])
			} else {
				reply = "What are you subscribing to?~"
			}
			return true, reply
		} else if strings.HasPrefix(command, "!unsubscribe") {
			if len(text) > 13 {
				reply = removeRole(s, channel.GuildID, m.Author.ID, text[13:])
			} else {
				reply = "What are you unsubscribing from?~"
			}
			return true, reply
		}
	}

	return false, reply
}
