package config

import (
	"os"
	"github.com/BurntSushi/toml"
	"path/filepath"
	"time"
	"log"
)

type ConfigData struct {
	Discord_token string

	DB_ip   string
	DB_name string
	DB_user string
	DB_pass string

	Devmode   bool
	Timestamp time.Time
}

var appConfig ConfigData

func LoadConfig() {

	newConf, err := readConfig()
	if err != nil {
		log.Fatal("Bad or No Config File!")
	}
	appConfig = newConf
}

func Config() ConfigData {

	if time.Since(appConfig.Timestamp) > time.Hour*12 {
		newConf, err := readConfig()
		if err == nil {
			appConfig = newConf
		}
	}

	return appConfig
}

func isDevmode() bool {
	return appConfig.Devmode
}

func readConfig() (ConfigData, error) {
	path, _ := os.Getwd()
	var configfile = filepath.Join(path, "config.toml")
	_, err := os.Stat(configfile)
	if err != nil {
		return ConfigData{}, err
	}

	cfg := ConfigData{}
	if _, err := toml.DecodeFile(configfile, &cfg); err != nil {
		return ConfigData{}, err
	}

	cfg.Timestamp = time.Now()
	return cfg, nil
}
