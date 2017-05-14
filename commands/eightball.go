package commands

import (
	"fmt"
	"math/rand"
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

func eightball(text string) string {
	answer := eightballResponses[rand.Intn(len(eightballResponses))]

	if len(text) > 7 {
		question := text[7:]

		return fmt.Sprintf("*%s* **%s**", question, answer)
	}

	return answer
}
