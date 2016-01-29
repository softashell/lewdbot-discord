package main

import (
	"bufio"
	"fmt"
	"github.com/pteichman/fate"
	"github.com/softashell/lewdbot-discord/regex"
	"log"
	"math"
	"os"
	"strings"
)

type Chat struct {
	brain *fate.Model
}

func NewChat() *Chat {
	return &Chat{
		fate.NewModel(fate.Config{}),
	}
}

func (c *Chat) learnFileLines(path string, simple bool) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	text := ""

	s := bufio.NewScanner(bufio.NewReader(f))
	for s.Scan() {
		line := s.Text()
		if !simple { //Learn all lines between empty lines
			if line == "" {
				c.learnMessage(text, false)
				text = ""
			} else {
				text += " " + line
			}
		} else { // Learn every line
			c.learnMessage(line, false)
		}
	}

	return s.Err()
}

func cleanMessage(message string) string {
	message = regex.Link.ReplaceAllString(message, "")
	message = regex.Emoticon.ReplaceAllString(message, "")
	message = regex.Junk.ReplaceAllString(message, "")
	message = regex.WikipediaCitations.ReplaceAllString(message, "")
	message = regex.Actions.ReplaceAllString(message, " ")
	message = regex.Russian.ReplaceAllString(message, "")
	message = regex.RepeatedWhitespace.ReplaceAllString(message, " ")

	return strings.TrimSpace(message)
}

func (c *Chat) learnMessage(text string, log bool) bool {
	text = cleanMessage(text)

	if len(text) < 5 ||
		len(text) > 1000 ||
		getWordCount(text) < 3 ||
		regex.JustPunctuation.MatchString(text) ||
		regex.LeadingNumbers.MatchString(text) ||
		generateEntropy(text) < 3.0 {
		return false // Text doesn't contain enough information
	}

	c.brain.Learn(text)

	if log {
		c.logMessage(text)
	}

	return true
}

func (c *Chat) generateReply(message string) string {
	reply := c.brain.Reply(message)
	reply = strings.TrimSpace(reply)

	reply = regex.TrailingPunctuation.ReplaceAllString(reply, "")
	reply = fmt.Sprintf("%s~", reply)

	c.learnMessage(message, true)

	return reply
}

func generateEntropy(s string) (e float64) {
	m := make(map[rune]bool)
	for _, r := range s {
		if m[r] {
			continue
		}
		m[r] = true
		n := strings.Count(s, string(r))
		p := float64(n) / float64(len(s))
		e += p * math.Log(p) / math.Log(2)
	}
	return math.Abs(e)
}

func getWordCount(s string) int {
	strs := strings.Fields(s)
	res := make(map[string]int)

	for _, str := range strs {
		res[strings.ToLower(str)]++
	}

	return len(res)
}

func (c *Chat) logMessage(message string) {
	f, err := os.OpenFile("./data/chatlog.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if _, err = f.WriteString(fmt.Sprintf("%s\n\n", message)); err != nil {
		log.Println(err)
	}
}
