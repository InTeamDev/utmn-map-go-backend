package config

import (
	"errors"
	"os"

	"gopkg.in/yaml.v2"
)

type PublicAPI struct {
	Server struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		DSN string `yaml:"dsn"`
	} `yaml:"database"`
}

func LoadPublicAPI(path string) (*PublicAPI, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg PublicAPI
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if cfg.Server.Host == "" {
		return nil, errors.New("server.host is required")
	}
	if cfg.Server.Port == 0 {
		return nil, errors.New("server.port is required")
	}
	if cfg.Database.DSN == "" {
		return nil, errors.New("database.dsn is required")
	}

	return &cfg, nil
}
