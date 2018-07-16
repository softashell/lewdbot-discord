package main

import (
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/bwmarrin/discordgo"
	"github.com/softashell/lewdbot-discord/config"
)

func fixStreamerRoles(s *discordgo.Session) {
	for _, g := range s.State.Guilds {
		if !config.GuildHasStreamerRoleEnabled(g.ID) {
			continue
		}

		// Fetch data manually
		g, err := s.Guild(g.ID)
		if err != nil {
			log.Error(err)
			continue
		}

		role, err := getRole(g, "Streamer")
		if err != nil {
			log.Error(err, " ", g.ID, " ", g.Name)
			continue
		}

		for _, p := range g.Presences {
			updateStreamerRole(s, p, g.ID, p.User.ID, role.ID)
		}
	}
}

func updateStreamerRole(s *discordgo.Session, p *discordgo.Presence, guildID, userID, roleID string) error {
	roleAdded, err := hasRole(s, guildID, userID, roleID)
	if err != nil {
		log.Errorf("updateStreamerRole: failed to get member roles - %s", err)
		return err
	}

	if p.Game == nil && roleAdded {
		log.Infof("updateStreamerRole: removing  streamer group from %s", userID)
		err = s.GuildMemberRoleRemove(guildID, userID, roleID)
		if err != nil {
			log.Errorf("updateStreamerRole: failed to remove streamer role - %s", err)
			return err
		}
	} else if p.Game != nil && p.Game.Type == discordgo.GameTypeStreaming && !roleAdded {
		log.Infof("updateStreamerRole: adding streamer group from %s", userID)
		err = s.GuildMemberRoleAdd(guildID, userID, roleID)
		if err != nil {
			log.Errorf("updateStreamerRole: failed to add streamer role - %s", err)
			return err
		}
	}

	return nil
}

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
	member, err := s.State.Member(guildID, userID)
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
