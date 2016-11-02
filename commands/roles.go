package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
)

func getMentionableRoles(s *discordgo.Session, GuildID string) []*discordgo.Role {
	g, err := s.State.Guild(GuildID)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	m, err := s.GuildMember(GuildID, s.State.User.ID)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	roleID := ""
	rolePos := 0

	if len(m.Roles) >= 1 {
		roleID = m.Roles[0]
	}

	for _, role := range g.Roles {
		if role.ID == roleID {
			rolePos = role.Position
			break
		}
	}

	var roles []*discordgo.Role

	for _, role := range g.Roles {
		if role.Name == "@everyone" || !role.Mentionable || role.Position > rolePos {
			continue
		}

		roles = append(roles, role)
	}

	return roles
}

func roleExists(g *discordgo.Guild, name string) (bool, *discordgo.Role) {
	for _, role := range g.Roles {
		if role.Name == "@everyone" {
			continue
		}

		if role.Name == name {
			return true, role
		}

	}

	return false, nil
}

func listRoles(s *discordgo.Session, GuildID string) string {
	fmt.Println("listRoles")

	var reply string

	roles := getMentionableRoles(s, GuildID)
	fmt.Println("Found", len(roles), "mentionable roles")

	roleCount := len(roles) - 1

	for i, role := range roles {
		fmt.Println(role)

		if i >= roleCount {
			reply += fmt.Sprintf("%s~", role.Name)
		} else {
			reply += fmt.Sprintf("%s, ", role.Name)
		}
	}

	if len(reply) <= 0 {
		return "I couldn't find any mentionable roles you could subscribe to~"
	}

	return reply
}

func listRoleMembers(s *discordgo.Session, GuildID string, arg string) string {
	fmt.Println("listRoleMembers")

	var reply string

	g, err := s.State.Guild(GuildID)
	if err != nil {
		fmt.Println(err)
		return "You fucking broke it~"
	}

	exists, role := roleExists(g, arg)

	if !exists {
		return "I can't find that~"
	}

	var members []*discordgo.Member

	for _, m := range g.Members {
		for _, r := range m.Roles {
			if r == role.ID {
				members = append(members, m)
				break
			}
		}
	}

	fmt.Println("Found", len(members), "members in", role.Name)

	memberCount := len(members) - 1

	for i, m := range members {
		var name string

		if len(m.Nick) > 0 {
			name = m.Nick
		} else {
			name = m.User.Username
		}

		reply += name

		if i >= memberCount {
			reply += "~"
		} else {
			reply += ", "
		}
	}

	if len(reply) <= 0 {
		return "The role seems empty~"
	}

	return reply
}

func addRole(s *discordgo.Session, GuildID string, UserID string, arg string) string {
	g, err := s.State.Guild(GuildID)
	if err != nil {
		fmt.Println(err)
		return "You fucking broke it~"
	}

	exists, role := roleExists(g, arg)

	if !exists {
		if !strings.HasPrefix(arg, "yes") {
			return fmt.Sprintf("I can't find such group~ Are you sure you didn't mistype it? Say **!subscribe yes %s** to create a new one~", arg)
		}

		if len(arg) < 5 {
			return "Are you sure you actually typed a name?~"
		}

		arg = arg[4:]

		exists, role = roleExists(g, arg)

		if !exists {
			newRole, err := s.GuildRoleCreate(GuildID)

			if err != nil {
				fmt.Println(err)
				return "Failed to create role"
			}

			role, err = s.GuildRoleEdit(GuildID, newRole.ID, arg, newRole.Color, newRole.Hoist, 37080064, true)
			if err != nil {
				fmt.Println(err)

				err = s.GuildRoleDelete(GuildID, newRole.ID)

				if err != nil {
					fmt.Println(err)
				}

				return "You fucking broke it~"
			}
			fmt.Println(role)
		} else {
			return "Why are you trying to recreate that group?"
		}
	}

	member, err := s.GuildMember(GuildID, UserID)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	if exists {
		for _, _role := range member.Roles {
			if _role == role.ID {
				return fmt.Sprintf("You're already subscribed to %s ~", arg)
			}
		}
	}

	newRoles := append(member.Roles, role.ID)

	err = s.GuildMemberEdit(GuildID, UserID, newRoles)
	if err != nil {
		fmt.Println(err)
		return "I can't touch that group dude, do it yourself~"
	}

	if exists {
		return fmt.Sprintf("You're now subscribed to %s~", arg)
	}

	return fmt.Sprintf("Created and subscribed to %s", arg)
}

func removeRole(s *discordgo.Session, GuildID string, UserID string, arg string) string {
	fmt.Println("removeRole", arg)

	g, err := s.State.Guild(GuildID)
	if err != nil {
		fmt.Println(err)
		return "You fucking broke it~"
	}

	exists, role := roleExists(g, arg)
	if !exists {
		return "I can't find such group~"
	}

	fmt.Println("Found?", exists, role)

	member, err := s.GuildMember(GuildID, UserID)
	if err != nil {
		fmt.Println(err)
	}

	pos := -1
	for i, r := range member.Roles {
		if r == role.ID {
			pos = i
			break
		}
	}

	if pos < 0 {
		return fmt.Sprintf("You're already not subscribed to %s~", arg)
	}

	member.Roles = append(member.Roles[:pos], member.Roles[pos+1:]...)
	err = s.GuildMemberEdit(GuildID, UserID, member.Roles)
	if err != nil {
		fmt.Println(err)
		return "I can't touch that group dude, do it yourself~"
	}

	delete := true
	for _, member := range g.Members {
		if member.User.ID == UserID {
			continue // Ignore self since it's not updated here yet
		}

		for _, r := range member.Roles {
			if r == role.ID {
				delete = false
				break
			}
		}
	}

	fmt.Println("Should delete it?", delete)

	if delete {
		err := s.GuildRoleDelete(GuildID, role.ID)
		if err != nil {
			fmt.Println(err)
			return fmt.Sprintf("Unsubscribed from but failed to delete %s~", arg)
		}

		fmt.Println("Unsubscribed and deleted")
		return fmt.Sprintf("Unsubscribed from and deleted %s~", arg)
	}

	fmt.Println("Unsubscribed")
	return fmt.Sprintf("Unsubscribed from %s~", arg)
}
