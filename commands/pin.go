package commands

import "github.com/bwmarrin/discordgo"

func pinMessage(s *discordgo.Session, m *discordgo.MessageCreate) string {
	messages, err := s.ChannelMessages(m.ChannelID, 100, "", "", "")
	if err != nil {
		return "Sorry, I am unable to read any messages in this channel~"
	}

	for _, message := range messages {
		if message.Author.ID == m.Author.ID && message.ID != m.ID {
			err := s.ChannelMessagePin(m.ChannelID, message.ID)
			if err != nil {
				return "Sorry, I am unable to pin the message~"
			}
			return "Pinned the message~"
		}
	}
	return "Sorry, I was unable to find a recent message to pin~"
}
