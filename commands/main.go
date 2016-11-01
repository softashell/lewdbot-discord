package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/softashell/lewdbot-discord/config"
	"math/rand"
	"strings"
)

var eightballResponses = []string{
	"Most definitely yes",
	"For sure",
	"As I see it, yes",
	"My sources say yes",
	"Yes",
	"Most likely",
	"Perhaps",
	"Maybe",
	"Not sure",
	"It is uncertain",
	"Ask me again later",
	"Don't count on it",
	"Probably not",
	"Very doubtful",
	"Most likely no",
	"Nope",
	"No",
	"My sources say no",
	"Dont even think about it",
	"Definitely no",
	"NO - It may cause disease contraction",
}

func ParseMessage(s *discordgo.Session, m *discordgo.MessageCreate, text string) (bool, string) {
	var reply string

	channel, err := s.State.Channel(m.ChannelID)
	if err != nil {
		fmt.Printf("s.Channel(%s) >> %s\n", m.ChannelID, err.Error())
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
				reply = ":ok_hand:"
			} else {
				reply = ":clap:"
			}

			return true, reply
		} else if strings.HasPrefix(command, "dumb") {
			if config.GuildSetDumb(channel.GuildID) {
				reply = ":ok_hand:"
			} else {
				reply = ":clap:"
			}

			return true, reply

		} else if strings.HasPrefix(command, "roles") {
			if config.SetManageRoles(channel.GuildID) {
				reply = ":ok_hand:"
			} else {
				reply = ":clap:"
			}

			return true, reply
		}

		return false, ""
	}

	if strings.HasPrefix(command, "!8ball") {
		reply = eightball(text)

		return true, reply
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

	return false, reply
}

func eightball(text string) string {
	answer := eightballResponses[rand.Intn(len(eightballResponses)-1)]

	if len(text) > 7 {
		question := text[7:]

		return fmt.Sprintf("*%s* **%s**", question, answer)
	}

	return answer
}
