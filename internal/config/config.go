package config

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/charmbracelet/x/editor"
	"github.com/rs/zerolog/log"

	"gopkg.in/yaml.v3"
)

// Config represents the entire configuration structure.
type Config struct {
	Blocklists []Blocklist `yaml:"blocklists"`
	Whitelists []Whitelist `yaml:"whitelists"`
}

// Blocklist is a structure that represents where all these domains that needs
// to be blocked are located.
type Blocklist struct {
	Target string `yaml:"target"`
}

// Whitelists contains domains that must not be blocked.
// Adding some domains to whitelist may fix many problems like YouTube
// watch history, videos on news sites and so on.
type Whitelist struct {
	Target string `yaml:"target"`
}

// Load loads config file.
// If config file is located at filesystem, it merges its options with
// default one and returns the result.
// If config file simply doesn't present at filesystem, it returns default one.
func Load() (*Config, error) {
	location := location()

	config, err := read(location)
	if errors.Is(err, os.ErrNotExist) {
		log.Debug().Msg("local config file doesn't exist. Default config loaded")
		return defaultConfig(), nil
	}
	if err != nil {
		return nil, err
	}

	if err := Validate(config); err != nil {
		return nil, err
	}

	log.Debug().Str("location", location).Msg("config file loaded")

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

	log.Debug().Str("location", location).Msg("config file loaded")

	return config, nil
}

// Print prints config to stdout.
func (c *Config) Print() {
	out, _ := yaml.Marshal(c)
	fmt.Print(string(out))
}

// Init saves default configuration file locally in case if
// it doesn't exist yet.
func Init() error {
	location := location()
	dirs := filepath.Dir(location)

	if _, err := os.Stat(location); !os.IsNotExist(err) {
		log.Debug().Str("location", location).Msg("config file has already been initialized")
		return nil
	}

	config, err := yaml.Marshal(defaultConfig())
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dirs, 0o700); err != nil {
		return err
	}

	if err := os.Chmod(dirs, 0o777); err != nil {
		return err
	}

	file, err := os.Create(location)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.Write(config); err != nil {
		return err
	}

	if err := os.Chmod(location, 0o777); err != nil {
		return err
	}

	log.Info().Str("location", location).Msg("config file has been initialized successfully")

	return nil
}

// Edit opens config in text editor.
func Edit() error {
	location := location()

	c, err := editor.Cmd("adless", location)
	if err != nil {
		return err
	}

	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	if err := c.Run(); err != nil {
		return err
	}

	return nil
}

// location returns the location of the config file.
func location() string {
	if bcp := os.Getenv("ADLESS_CONFIG_PATH"); bcp != "" {
		return bcp
	}

	if bch := os.Getenv("ADLESS_CONFIG_HOME"); bch != "" {
		return filepath.Join(bch, "config.yml")
	}

	if xdgConfig := os.Getenv("XDG_CONFIG_HOME"); xdgConfig != "" {
		return filepath.Join(xdgConfig, "adless", "config.yml")
	}

	return filepath.Join(homeDir(), ".config", "adless", "config.yml")
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

func homeDir() string {
	username := os.Getenv("SUDO_USER")
	if username == "" {
		return os.Getenv("HOME")
	}

	sudoUser, _ := user.Lookup(username)
	return sudoUser.HomeDir
}

// defaultConfig returns precreated configuration.
// In case if there is no config file stored on user's filesystem,
// default one will be used.
func defaultConfig() *Config {
	return &Config{
		Blocklists: []Blocklist{
			{Target: "https://raw.githubusercontent.com/StevenBlack/hosts/master/data/StevenBlack/hosts"},
		},
		Whitelists: []Whitelist{
			{Target: "https://raw.githubusercontent.com/anudeepND/whitelist/master/domains/whitelist.txt"},
		},
	}
}
