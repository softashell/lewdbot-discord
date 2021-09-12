package brain

import (
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/zenthangplus/goccm"
)

func learnInBackground() {
	c := goccm.New(2000)
	for text := range textChan {
		c.Wait()

		go func(t learningText) {
			defer c.Done()

			Learn(t.text, t.saveText)  
        }(text)
	}
}

func backgroundMessageLogger() {
	t := time.Now()
	fileName := fmt.Sprintf("./data/chatlog-%d-%02d.txt", t.Year(), t.Month())

	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		log.Error("Unable to open log", err)
	}
	defer f.Close()

	for message := range textLogChan {
		if _, err = f.WriteString(fmt.Sprintf("%s\n\n", message)); err != nil {
			log.Error("Unable to write in log", err)
		}
	}
}