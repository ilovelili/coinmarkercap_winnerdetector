// Package config configuration related
package config

import (
	"encoding/json"
	"os"
	"path"
)

// GetConfig get config defined in config.json
func GetConfig() (config *Config) {
	pwd, _ := os.Getwd()
	path := path.Join(pwd, "config.json")
	configFile, err := os.Open(path)
	defer configFile.Close()

	if err != nil {
		panic(err)
	}

	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&config); err != nil {
		panic(err)
	}

	return
}

// Sender mail sender
type Sender struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

// Mail Mail service config
type Mail struct {
	Receivers string `json:"receivers"`
	Sender    `json:"sender"`
}

// Threshold Percentage change threshold in hour
type Threshold struct {
	Max float64 `json:"maxpercentchangeinhour"`
	Min float64 `json:"minpercentchangeinhour"`
}

// Config config entry
type Config struct {
	Threshold `json:"threshold"`
	Mail      `json:"mail"`
}
