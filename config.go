package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

func LoadConfigFromFile(filename string) (string, string) {
	fileDump, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalln(err.Error())
	}

	type fileCredentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var creds = fileCredentials{}
	if err := json.Unmarshal(fileDump, &creds); err != nil {
		log.Fatalln(err.Error())
	}

	return creds.Email, creds.Password
}
