package commands

import (
	"github.com/bwmarrin/discordgo"
)

func digits(s *discordgo.Session, m *discordgo.MessageCreate) string {
	if len(m.Mentions) == 0 || m.MessageReference != nil {
		return m.ID
	}

	messages, err := s.ChannelMessages(m.ChannelID, 100, "", "", "")
	if err != nil {
		return "Sorry, I am unable to read any messages in this channel~"
	}

	for _, message := range messages {
		if message.Author.ID == m.Mentions[0].ID && message.ID != m.ID {
			return "He got " + message.ID
		}
	}
	return "He needs to speak up first~"
}
