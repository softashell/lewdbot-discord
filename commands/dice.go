package commands

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

func dice(text string, author *discordgo.User) string {
	match := regexp.MustCompile(`(\d+)d(\d+)`).FindStringSubmatch(text)

	if len(match) < 3 {
		return fmt.Sprintf("%s, you fucked up~", author.Username)
	}

	dice, err := strconv.Atoi(match[1])
	if err != nil {
		return fmt.Sprintf("%s, you fucked up~", author.Username)
	}

	sides, err := strconv.Atoi(match[2])
	if err != nil {
		return fmt.Sprintf("%s, you fucked up~", author.Username)
	}

	if dice <= 0 || sides <= 0 || dice > 100 || sides > 100 {
		return fmt.Sprintf("%s, fuck off~", author.Username)
	}

	if sides == 1 {
		return fmt.Sprintf("%s, you rolled %d~ What else did you expect?~", author.Username, dice*sides)
	}

	roll := 0

	diceResults := make([]int, 0)
	var details string

	for dice > 0 {
		rollResult := rand.Intn(sides) + 1
		roll += rollResult
		diceResults = append(diceResults, roll)

		details += fmt.Sprintf("%d", rollResult)

		dice--
		if dice >= 1 {
			details += " + "
		}
	}

	if len(diceResults) <= 1 {
		return fmt.Sprintf("%s, you rolled %d~", author.Username, roll)
	}

	return fmt.Sprintf("%s, you rolled %d (%s)~", author.Username, roll, details)
}
