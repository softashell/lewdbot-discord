package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
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

	if strings.HasPrefix(command, "!8ball") {
		reply = eightball(text)

		return true, reply
	} else if channel.GuildID == "111928847846367232" || channel.GuildID == "135827109485608960" {
		if strings.HasPrefix(command, "!list") {

			reply = listRoles(s, channel.GuildID)

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

func listRoles(s *discordgo.Session, GuildID string) string {
	g, err := s.Guild(GuildID)
	if err != nil {
		fmt.Println(err)
	}

	u, err := s.User("@me")
	if err != nil {
		fmt.Println("error obtaining account details,", err)
	}

	m, err := s.GuildMember(GuildID, u.ID)
	if err != nil {
		fmt.Println(err)
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

	var reply string
	roles := len(g.Roles) - 1

	for i, role := range g.Roles {
		fmt.Println(role)

		if role.Name == "@everyone" || !role.Mentionable || role.Position > rolePos {
			continue
		}

		if i >= roles {
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

func addRole(s *discordgo.Session, GuildID string, UserID string, arg string) string {
	g, err := s.Guild(GuildID)
	if err != nil {
		fmt.Println(err)
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
	g, err := s.Guild(GuildID)
	if err != nil {
		fmt.Println(err)
	}

	exists, role := roleExists(g, arg)

	fmt.Println(arg, exists, role)

	if !exists {
		return "I can't find such group~"
	}

	member, err := s.GuildMember(GuildID, UserID)
	if err != nil {
		fmt.Println(err)
	}

	found := false
	pos := 0

	for i, _role := range member.Roles {
		if _role == role.ID {
			found = true
			pos = i
		}
	}

	if !found {
		return fmt.Sprintf("You're already not subscribed to %s~", arg)
	}

	member.Roles = append(member.Roles[:pos], member.Roles[pos+1:]...)

	err = s.GuildMemberEdit(GuildID, UserID, member.Roles)
	if err != nil {
		fmt.Println(err)
		return "I can't touch that group dude, do it yourself~"
	}

	members, err := s.GuildMembers(GuildID, 0, 1000)
	if err != nil {
		fmt.Println(err)
	}

	delete := true

	for _, member := range members {
		for _, _role := range member.Roles {
			if _role == role.ID {
				delete = false
				break
			}
		}
	}
	fmt.Println(delete, role)
	if delete {
		err := s.GuildRoleDelete(GuildID, role.ID)
		if err != nil {
			fmt.Println(err)
			return fmt.Sprintf("Unsubscribed from but failed to delete %s~", arg)
		}

		return fmt.Sprintf("Unsubscribed from and deleted %s~", arg)
	}

	return fmt.Sprintf("Unsubscribed from %s~", arg)
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
