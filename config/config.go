package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	log "github.com/Sirupsen/logrus"
)

type Config struct {
	LoginCredentials `json:"login"`
	Brain            []brainFile              `json:"brain"`
	Blacklist        []string                 `json:"blacklist"`
	Guilds           map[string]guildSettings `json:"guilds"`
	Masters          []string                 `json:"masters"`
	Lastfm           `json:"lastfm"`
}

type LoginCredentials struct {
	Token string `json:"token"`
}

type Lastfm struct {
	lock      sync.RWMutex
	Key       string            `json:"api_key"`
	Usernames map[string]string `json:"usernames"`
}

type brainFile struct {
	File   string `json:"file"`
	Simple bool   `json:"simple"`
}

type guildSettings struct {
	Channels    map[string]channelSettings `json:"channels"`
	Dumb        bool                       `json:"dumb"`
	ManageRoles bool                       `json:"roles"`
	Lastfm      bool                       `json:"lastfm"`
}

type channelSettings struct {
	Lewd bool `json:"lewd"`
	Spam bool `json:"spam"`
}

var c *Config

func Init() {
	c = loadConfigFromFile("./data/config.json")

	if c.Guilds == nil {
		c.Guilds = make(map[string]guildSettings)
	}

	for _, g := range c.Guilds {
		if g.Channels == nil {
			g.Channels = make(map[string]channelSettings)
		}
	}

	c.Lastfm.lock = sync.RWMutex{}
	if c.Lastfm.Usernames == nil {
		c.Lastfm.Usernames = make(map[string]string)
	}

	if len(c.Token) == 0 {
		log.Fatal("Unable to load login information, did you set it in config?")
	}

	Save()
}

func Get() *Config {
	return c
}

func Save() {
	_json, err := json.MarshalIndent(&c, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile("./data/config.json", []byte(_json), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func Print(c *Config) {
	// Print out current config
	_json, err := json.MarshalIndent(&c, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	log.Println(string(_json))
}

func loadConfigFromFile(filename string) *Config {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err.Error())
	}

	var config = &Config{}

	if err := json.NewDecoder(bufio.NewReader(f)).Decode(&config); err != nil {
		log.Fatal(err.Error())
	}

	return config
}

func IsMaster(id string) bool {
	for _, u := range c.Masters {
		if u == id {
			return true
		}
	}

	return false
}

func GuildSetDumb(guild string) bool {
	g := c.Guilds[guild]

	g.Dumb = !g.Dumb

	c.Guilds[guild] = g

	Save()

	return g.Dumb
}

func GuildIsDumb(guild string) bool {
	return c.Guilds[guild].Dumb
}

func GuildSetLastfm(guild string) bool {
	g := c.Guilds[guild]

	g.Lastfm = !g.Lastfm

	c.Guilds[guild] = g

	Save()

	return g.Lastfm
}

func GuildHasLastfmEnabled(guild string) bool {
	return c.Guilds[guild].Lastfm
}

func ChannelSetLewd(guild string, channel string) bool {
	g := c.Guilds[guild]

	if g.Channels == nil {
		g.Channels = make(map[string]channelSettings)
	}

	ch := g.Channels[channel]

	ch.Lewd = !ch.Lewd

	g.Channels[channel] = ch
	c.Guilds[guild] = g

	Save()

	return ch.Lewd
}

func ChannelIsLewd(guild string, channel string) bool {
	return c.Guilds[guild].Channels[channel].Lewd
}

func ChannelSetSpam(guild string, channel string) bool {
	g := c.Guilds[guild]

	if g.Channels == nil {
		g.Channels = make(map[string]channelSettings)
	}

	ch := g.Channels[channel]

	ch.Spam = !ch.Spam

	g.Channels[channel] = ch
	c.Guilds[guild] = g

	Save()

	return ch.Spam
}

func ChannelShouldSpam(guild string, channel string) bool {
	return c.Guilds[guild].Channels[channel].Spam
}

func SetManageRoles(guild string) bool {
	g := c.Guilds[guild]

	g.ManageRoles = !g.ManageRoles

	c.Guilds[guild] = g

	Save()

	return g.ManageRoles
}

func ShouldManageRoles(guild string) bool {
	return c.Guilds[guild].ManageRoles
}

func SetLastfmUsername(UserID string, username string) {
	c.Lastfm.lock.Lock()
	defer c.Lastfm.lock.Unlock()

	log.Println("Setting username for", UserID, "to", username)

	c.Lastfm.Usernames[UserID] = username

	Save()
}

func GetLastfmUsername(UserID string) (string, error) {
	c.Lastfm.lock.RLock()
	defer c.Lastfm.lock.RUnlock()

	username := c.Lastfm.Usernames[UserID]
	if len(username) < 1 {
		return "", fmt.Errorf("Couldn't find saved last.fm username")
	}

	return username, nil
}

func RemoveLastfmUsername(UserID string) {
	c.Lastfm.lock.Lock()
	defer c.Lastfm.lock.Unlock()

	log.Println("Removed associated username for", UserID)

	delete(c.Lastfm.Usernames, UserID)

	Save()
}

func GetLastfmKey() string {
	return c.Lastfm.Key
}
