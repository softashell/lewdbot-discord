package config

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
)

type Config struct {
	loginCredentials `json:"login"`
	Brain            []brainFile              `json:"brain"`
	Blacklist        []string                 `json:"blacklist"`
	Guilds           map[string]guildSettings `json:"guilds"`
	Masters          []string                 `json:"masters"`
}

type loginCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

type brainFile struct {
	File   string `json:"file"`
	Simple bool   `json:"simple"`
}

type guildSettings struct {
	Channels    map[string]channelSettings `json:"channels"`
	Dumb        bool                       `json:"dumb"`
	ManageRoles bool                       `json:"roles"`
}

type channelSettings struct {
	Lewd   bool `json:"lewd"`
	Pso2eq bool `json:"pso2"`
}

var c Config

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

	if (len(c.Email) == 0 || len(c.Password) == 0) && len(c.Token) == 0 {
		log.Fatal("Unable to load login information, did you set it in config?")
	}

	Save()
}

func Get() *Config {
	return &c
}

func Save() {
	_json, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile("./data/config.json", []byte(_json), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func Print(c Config) {
	// Print out current config
	_json, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	log.Println(string(_json))
}

func loadConfigFromFile(filename string) Config {
	fileDump, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	var config = Config{}

	if err := json.Unmarshal(fileDump, &config); err != nil {
		log.Fatal(err)
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
