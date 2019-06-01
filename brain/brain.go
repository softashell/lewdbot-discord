package brain

import (
	"bufio"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/pteichman/fate"
	"github.com/softashell/lewdbot-discord/regex"
	"github.com/tebeka/snowball"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"math"
	"os"
	"strings"
	"time"
	"unicode"
)

var lewdbrain *fate.Model

// Init Sets the global fate model
func Init() {
	lewdbrain = fate.NewModel(fate.Config{Stemmer: newStemmer()})
}

type stemmer struct {
	tran     transform.Transformer
	snowball *snowball.Stemmer
}

func newStemmer() stemmer {
	isRemovable := func(r rune) bool {
		return unicode.Is(unicode.Mn, r) || unicode.IsPunct(r)
	}

	stem, err := snowball.New("english")
	if err != nil {
		log.Fatal("Unable to create new stemmer", err)
	}

	return stemmer{
		tran:     transform.Chain(norm.NFD, transform.RemoveFunc(isRemovable), norm.NFC),
		snowball: stem,
	}
}

func (s stemmer) Stem(word string) string {
	str, _, _ := transform.String(s.tran, word)
	return squish(s.snowball.Stem(strings.ToLower(str)), 2)
}

// Squish 2+ consecutive characters together during stemming
func squish(s string, max int) string {
	var (
		ret []rune
		cur rune
		n   int
	)

	emit := func(r rune, n int) {
		if n > max {
			n = max
		}
		for i := 0; i < n; i++ {
			ret = append(ret, r)
		}
	}

	for _, r := range s {
		if r == cur {
			n++
		} else {
			emit(cur, n)
			cur = r
			n = 1
		}
	}

	emit(cur, n)
	return string(ret)
}

// Learn Attempts to learn and log input text
func Learn(text string, log bool) bool {
	text = cleanMessage(text)

	if len(text) < 5 ||
		len(text) > 1000 ||
		getWordCount(text) < 3 ||
		regex.JustPunctuation.MatchString(text) ||
		regex.LeadingNumbers.MatchString(text) ||
		generateEntropy(text) < 3.0 {
		return false // Text doesn't contain enough information
	}

	lewdbrain.Learn(text)

	if log {
		logMessage(text)
	}

	return true
}

// LearnFileLines Attempts to learn input file
func LearnFileLines(path string, simple bool) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	defer f.Close()

	var text string

	s := bufio.NewScanner(bufio.NewReader(f))

	for s.Scan() {
		line := s.Text()
		if !simple { //Learn all lines between empty lines
			if line == "" {
				Learn(text, false)
				text = ""
			} else {
				text += " " + line
			}
		} else { // Learn every line
			Learn(line, false)
		}
	}

	return s.Err()
}

// Reply Returns generated reply to input and learns it
func Reply(message string) string {
	reply := lewdbrain.Reply(message)

	reply = strings.TrimSpace(reply)
	reply = regex.TrailingPunctuation.ReplaceAllString(reply, "")
	reply = fmt.Sprintf("%s~", reply)

	Learn(message, true)

	return reply
}

func logMessage(message string) {
	t := time.Now()
	fileName := fmt.Sprintf("./data/chatlog-%d-%02d.txt", t.Year(), t.Month())

	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		log.Error("Unable to open log", err)
	}

	defer f.Close()

	if _, err = f.WriteString(fmt.Sprintf("%s\n\n", message)); err != nil {
		log.Error("Unable to write in log", err)
	}
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
