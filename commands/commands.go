package commands

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
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
		} else if strings.HasPrefix(command, "spam") {
			if config.ChannelSetSpam(channel.GuildID, channel.ID) {
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
		} else if strings.HasPrefix(command, "lastfm") {
			if config.GuildSetLastfm(channel.GuildID) {
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
		} else if strings.HasPrefix(command, "streamer") {
			if config.GuildSetStreamerRole(channel.GuildID) {
				reply = msgOn
			} else {
				reply = msgOff
			}

			return true, reply
		} else if strings.HasPrefix(command, "mangadex") {
			if config.ChannelSetMangadex(channel.GuildID, channel.ID) {
				reply = msgOn
			} else {
				reply = msgOff
			}

			return true, reply
		}

		return false, ""
	}

	if config.ShouldManageRoles(channel.GuildID) {
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

	if config.GuildHasLastfmEnabled(channel.GuildID) {
		if strings.HasPrefix(command, "!np") {
			err = s.ChannelTyping(m.ChannelID)
			if err != nil {
				log.Warn("s.ChannelTyping >> ", err)
			}

			if strings.HasPrefix(command, "!np set") && len(text) > 8 {
				reply = registerLastfmProfile(m.Author.ID, text[8:])
			} else if strings.HasPrefix(command, "!np remove") {
				reply = removeLastfmProfile(m.Author.ID)
			} else {
				reply = spamNowPlayingUser(m.Author.ID)
			}

			return true, reply
		} else if strings.HasPrefix(command, "!wp") {
			err = s.ChannelTyping(m.ChannelID)
			if err != nil {
				log.Warn("s.ChannelTyping >> ", err)
			}

			reply = spamNowPlayingServer(s, channel.GuildID)

			return true, reply
		}
	}

	if strings.HasPrefix(command, "!8ball") {
		reply = eightball(text)

		return true, reply
	} else if strings.HasPrefix(command, "!roll") {
		reply = dice(text[5:], m.Author)

		return true, reply
	} else if strings.HasPrefix(command, "!pin") {
		if perms, err := s.State.UserChannelPermissions(m.Author.ID, m.ChannelID); err == nil {
			if (perms&discordgo.PermissionManageWebhooks) != 0 || (perms&discordgo.PermissionAdministrator) != 0 {
				return true, pinMessage(s, m)
			}
			return true, "You don't have permission to do that~"
		}
		return true, "Sorry, something went wrong~"
	} else if strings.HasPrefix(command, "!digits") {
		return true, digits(s, m)
	}

	return false, reply
}
