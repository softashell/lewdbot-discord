package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func getRole(g *discordgo.Guild, name string) (*discordgo.Role, error) {
	name = strings.ToLower(name)

	for _, role := range g.Roles {
		if role.Name == "@everyone" {
			continue
		}

		if strings.ToLower(role.Name) == name {
			return role, nil
		}
	}

	return nil, fmt.Errorf("couldn't find role: %s", name)
}

func hasRole(s *discordgo.Session, guildID, userID, roleID string) (bool, error) {
	member, err := s.GuildMember(guildID, userID)
	if err != nil {
		return false, err
	}

	for _, role := range member.Roles {
		if role == roleID {
			return true, nil
		}
	}

	return false, nil
}
