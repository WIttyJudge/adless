package config

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog/log"

	"gopkg.in/yaml.v3"
)

// Config represents the entire configuration structure
type Config struct {
	Blocklists []Blocklist `yaml:"blocklists"`
	HTTP       HTTP        `yaml:"http"`
}

type Blocklist struct {
	Target string `yaml:"target"`
}

type HTTP struct {
	Timeout time.Duration `yaml:"timeout"`
}

func Load() (*Config, error) {
	location := configLocation()

	config, err := readConfig(location)
	if errors.Is(err, os.ErrNotExist) {
		log.Debug().Msg("local config file doesn't exist. Default configuration loaded")
		return defaultConfig(), nil
	}
	if err != nil {
		return nil, err
	}

	if err := Validate(config); err != nil {
		return nil, err
	}

	log.Debug().Str("config", location).Msg("config file loaded")

	return config, nil
}

func LoadByUser(location string) (*Config, error) {
	config, err := readConfig(location)
	if err != nil {
		return nil, err
	}

	if err := Validate(config); err != nil {
		return nil, err
	}

	log.Debug().Str("config", location).Msg("config file loaded")

	return config, nil
}

// configLocation returns the location of the config file
func configLocation() string {
	if bch := os.Getenv("BARRIER_CONFIG_HOME"); bch != "" {
		return filepath.Join(bch, "config.yml")
	}

	if xdgConfig := os.Getenv("XDG_CONFIG_HOME"); xdgConfig != "" {
		return filepath.Join(xdgConfig, "barrier", "config.yml")
	}

	return filepath.Join(os.Getenv("HOME"), "config", "barrier", "config.yml")
}

func readConfig(location string) (*Config, error) {
	data, err := os.ReadFile(location)
	if err != nil {
		return nil, err
	}

	config := defaultConfig()
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, err
	}

	return config, nil
}

func defaultConfig() *Config {
	httpTimeout, _ := time.ParseDuration("10s")

	return &Config{
		Blocklists: []Blocklist{
			{Target: "https://raw.githubusercontent.com/StevenBlack/hosts/master/data/StevenBlack/hosts"},
		},
		HTTP: HTTP{
			Timeout: httpTimeout,
		},
	}
}
