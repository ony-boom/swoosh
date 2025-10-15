package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Config struct {
	HideSink            []string `json:"hideSink"`
	PollIntervalSeconds int      `json:"pollIntervalSeconds"`
}

var defaultConfig = Config{
	HideSink:            []string{},
	PollIntervalSeconds: 2,
}

func New() Config {
	filePath, err := xdgConfigFile(defaultConfig)
	if err != nil {
		log.Println(err)
		return defaultConfig
	}

	file, err := os.Open(filePath)
	if err != nil {
		log.Println(fmt.Errorf("failed to open config file :%v", err))
		return defaultConfig
	}

	defer file.Close()

	config := defaultConfig

	decoder := json.NewDecoder(file)

	if err = decoder.Decode(&config); err != nil {
		log.Println(fmt.Errorf("failed to deconde config file: %v", err))
		return defaultConfig
	}

	return config
}
