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
		log.Infof("Guild: %s (%s) Members: %d ", g.ID, g.Name, g.MemberCount)
		if !config.GuildHasStreamerRoleEnabled(g.ID) {
			continue
		}

		role, err := getRole(g, "Streamer")
		if err != nil {
			log.Error(err, " ", g.ID, " ", g.Name)
			continue
		}

		for _, m := range g.Members {
			//log.Info("Member:", m.User.ID, " ", m.User.Username)

			p, err := s.State.Presence(g.ID, m.User.ID)
			if err != nil {
				log.Warnf("failed to get presence for %s ( %s || %s) - %s", m.User.ID, m.User.Username, m.Nick, err)
				continue
			}

			updateStreamerRole(s, p, g.ID, m.User.ID, role.ID)
		}
	}
}

func updateStreamerRole(s *discordgo.Session, p *discordgo.Presence, guildID, userID, roleID string) error {
	roleAdded, err := hasRole(s, guildID, userID, roleID)
	if err != nil {
		log.Errorf("updateStreamerRole: failed to get member roles - %s", err)
		return err
	}

	if p != nil && p.Game != nil && p.Game.Type == discordgo.GameTypeStreaming && !roleAdded {
		log.Infof("updateStreamerRole: adding streamer group from %s", userID)
		err = s.GuildMemberRoleAdd(guildID, userID, roleID)
		if err != nil {
			log.Errorf("updateStreamerRole: failed to add streamer role - %s", err)
			return err
		}

		return nil
	}

	if roleAdded {
		log.Infof("updateStreamerRole: removing  streamer group from %s", userID)
		err = s.GuildMemberRoleRemove(guildID, userID, roleID)
		if err != nil {
			log.Errorf("updateStreamerRole: failed to remove streamer role - %s", err)
			return err
		}

		return nil
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
