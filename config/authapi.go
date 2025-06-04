package config

import (
	"errors"
	"os"

	"gopkg.in/yaml.v2"
)

type AuthAPI struct {
	Server struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"server"`
	JWTSecret string `yaml:"jwt_secret"`
	BotClient struct {
		URL         string `yaml:"url"`
		ClientID    string `yaml:"client_id"`
		AccessToken string `yaml:"access_token"`
	} `yaml:"bot_client"`
	Auth struct {
		ClientID    string `yaml:"client_id"`
		AccessToken string `yaml:"access_token"`
	} `yaml:"auth"`
}

func LoadAuthAPI(path string) (*AuthAPI, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg AuthAPI
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	if cfg.Server.Host == "" || cfg.Server.Port == 0 {
		return nil, errors.New("server config required")
	}
	if cfg.JWTSecret == "" {
		return nil, errors.New("jwt_secret required")
	}
	return &cfg, nil
}
