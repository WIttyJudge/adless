package config

import (
	"errors"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/rs/zerolog/log"

	"gopkg.in/yaml.v3"
)

// Config represents the entire configuration structure.
type Config struct {
	Blocklists []Blocklist `yaml:"blocklists"`
	HTTP       HTTP        `yaml:"http"`
}

// Blocklist is a strcture that represents where all these domains we need to
// block are located.
type Blocklist struct {
	Target string `yaml:"target"`
}

type HTTP struct {
	Timeout time.Duration `yaml:"timeout"`
}

// Load loads config file.
// If config file is located at filesystem, it merges its options with
// default one and returns the result.
// If config file simply doesn't present at filesystem, it returns default one.
func Load() (*Config, error) {
	location := location()

	config, err := read(location)
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

// LoadByUser loads config file at the location that was provided
// by the user via config-file CLI flag.
func LoadByUser(location string) (*Config, error) {
	config, err := read(location)
	if err != nil {
		return nil, err
	}

	if err := Validate(config); err != nil {
		return nil, err
	}

	log.Debug().Str("config", location).Msg("config file loaded")

	return config, nil
}

// location returns the location of the config file.
func location() string {
	username := os.Getenv("SUDO_USER")
	sudoUser, _ := user.Lookup(username)

	if bch := os.Getenv("BARRIER_CONFIG_HOME"); bch != "" {
		return filepath.Join(bch, "config.yml")
	}

	if xdgConfig := os.Getenv("XDG_CONFIG_HOME"); xdgConfig != "" {
		return filepath.Join(xdgConfig, "barrier", "config.yml")
	}

	return filepath.Join(sudoUser.HomeDir, ".config", "barrier", "config.yml")
}

// read reads config file by location in file system.
func read(location string) (*Config, error) {
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

// defaultConfig returns precreated configuration.
// In case if there is no config file stored on user's filesystem,
// default one will be used.
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
