package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	loginCredentials `json:"login"`
	Blacklist        []string `json:"blacklist"`
}

type loginCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var c Config

func Init() {
	c = loadConfigFromFile("./data/config.json")

	if len(c.Email) == 0 || len(c.Password) == 0 {
		log.Println("Unable to load login information, did you set it in config?")
	}
}

func Get() *Config {
	return &c
}

func Print(c Config) {
	// Print out current config
	_json, _ := json.MarshalIndent(c, "", "\t")
	log.Println(string(_json))
}

func loadConfigFromFile(filename string) Config {
	fileDump, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalln(err.Error())
	}

	var config = Config{}

	if err := json.Unmarshal(fileDump, &config); err != nil {
		log.Fatalln(err.Error())
	}

	return config
}
