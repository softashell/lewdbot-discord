package commands

import "github.com/bwmarrin/discordgo"

func pinMessage(s *discordgo.Session, m *discordgo.MessageCreate) string {
	err := s.ChannelMessagePin(m.ChannelID, m.ID)
	if err != nil {
		return "Sorry, I am unable to pin the message~"
	}
	return "Pinned the message~"
}
