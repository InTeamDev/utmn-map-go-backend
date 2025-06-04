package config

import (
    "errors"
    "os"

    "gopkg.in/yaml.v2"
)

type Bot struct {
    Server struct {
        Host string `yaml:"host"`
        Port int    `yaml:"port"`
    } `yaml:"server"`
    Bot struct {
        Token string `yaml:"token"`
    } `yaml:"bot"`
    Backend struct {
        URL         string `yaml:"url"`
        ClientID    string `yaml:"client_id"`
        AccessToken string `yaml:"access_token"`
    } `yaml:"backend"`
    Auth struct {
        ClientID    string `yaml:"client_id"`
        AccessToken string `yaml:"access_token"`
    } `yaml:"auth"`
}

func LoadBot(path string) (*Bot, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }

    var cfg Bot
    if err := yaml.Unmarshal(data, &cfg); err != nil {
        return nil, err
    }

    if cfg.Server.Host == "" {
        return nil, errors.New("server.host is required")
    }
    if cfg.Server.Port == 0 {
        return nil, errors.New("server.port is required")
    }
    if cfg.Bot.Token == "" {
        return nil, errors.New("bot.token is required")
    }
    return &cfg, nil
}
