package config

import "errors"

var (
	ErrNoBlocklistsProvided = errors.New("no blocklists provided")
)

func Validate(config *Config) error {
	if len(config.Blocklists) == 0 {
		return ErrNoBlocklistsProvided
	}

	return nil
}
