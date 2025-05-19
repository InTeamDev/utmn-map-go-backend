package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Database Database `yaml:"database"`
	HTTP     HTTP     `yaml:"http"`
	TGBot    TGBot    `yaml:"tgbot"`
	JWT      JWT      `yaml:"jwt"`
}

type Database struct {
	DSN string `yaml:"dsn"`
}

type HTTP struct {
	Port int `yaml:"port"`
}

type TGBot struct {
	Token            string `yaml:"token"`
	DevelopersChatID int64  `yaml:"developers_chat_id"`
}

type JWT struct {
	Secret          string `yaml:"secret"`
	ExpirationHours int    `yaml:"expiration_hours"`
}

func New(path string) (*Config, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var c Config
	err = yaml.Unmarshal(bytes, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
