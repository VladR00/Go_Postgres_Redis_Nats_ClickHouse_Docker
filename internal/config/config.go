package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	User     string `json:"user"`
	Password string `json:"pass"`
	Host     string `json:"host"`
	DbName   string `json:"dbname"`
	PoolSize int32  `json:"poolsize"`
}

var (
	cfgPATH string = "../internal/config/config.json"
)

func LoadConfig() (Config, error) {
	var config Config
	file, err := os.Open(cfgPATH)
	if err != nil {
		return config, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	return config, err
}
