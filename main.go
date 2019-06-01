package main

import (
	"os"
	"os/signal"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
	"github.com/softashell/lewdbot-discord/brain"
	"github.com/softashell/lewdbot-discord/config"
)

const maxConnectionFailures = 5

func main() {
	err := os.Mkdir("./data", 0700)
	if err != nil && !os.IsExist(err) {
		log.Errorln("Can't create data directory", err)
		return
	}

	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt)

	config.Init()
	brain.Init()

	go fillBrain()

	failures := 0

	go func() {
		for failures < maxConnectionFailures {
			if err := connectToDiscord(); err != nil {
				log.Error(err)
				time.Sleep(25 * time.Second)

				failures++
			} else {
				break
			}
		}

		if failures >= maxConnectionFailures-1 {
			log.Error("maximum failures reached while starting up")
			interrupt <- os.Interrupt
		}
	}()

	<-interrupt

	log.Info("Shutting down")
}

func fillBrain() {
	c := config.Get()

	start := time.Now()

	log.Println("Starting to fill brain")

	for _, b := range c.Brain {
		log.Println("Parsing", b.File)

		if err := brain.LearnFileLines(b.File, b.Simple); err != nil {
			log.WithFields(log.Fields{
				"file":   b.File,
				"simple": b.Simple,
			}).Warn(err)
		}
	}

	if logs, err := filepath.Glob("./data/chatlog-*.txt"); err != nil {
		log.Error(err)
	} else {
		for _, l := range logs {
			log.Println("Parsing", l)

			if err := brain.LearnFileLines(l, false); err != nil {
				log.WithFields(log.Fields{
					"file": l,
				}).Warn(err)
			}
		}
	}

	log.Println("Brain filled in", time.Since(start))
}

func connectToDiscord() error {
	log.Println("Connecting to discord")

	var err error

	c := config.Get()

	dg, err := discordgo.New("Bot " + c.Token)
	if err != nil {
		return errors.Wrap(err, "failed to create discordgo")
	}

	dg.AddHandler(ready)
	dg.AddHandler(messageCreate)
	dg.AddHandler(presenceUpdate)
	dg.AddHandler(guildMembersChunk)

	// Retry after broken websocket
	dg.ShouldReconnectOnError = true

	// Verify the Token is valid and grab user information
	dg.State.User, err = dg.User("@me")
	if err != nil {
		return errors.Wrapf(err, "error fetching user information")
	}

	// Open websocket connection
	err = dg.Open()
	if err != nil {
		return errors.Wrap(err, "error opening connection to Discord")
	}

	log.Println("Connected")

	return nil
}

func shouldIgnore(user *discordgo.User) bool {
	c := config.Get()

	if user.Bot {
		return true
	}

	for _, id := range c.Blacklist {
		if id == user.ID {
			return true
		}
	}

	return false
}

func isUserMentioned(user *discordgo.User, mentions []*discordgo.User) bool {
	for _, mention := range mentions {
		if mention.ID == user.ID {
			return true
		}
	}

	return false
}
