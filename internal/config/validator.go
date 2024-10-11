package config

import (
	"errors"
	"fmt"
	"regexp"
)

var (
	ErrNoBlocklistsProvided = errors.New("no blocklists provided")
)

func Validate(config *Config) error {
	if len(config.Blocklists) == 0 {
		return ErrNoBlocklistsProvided
	}

	for _, blocklist := range config.Blocklists {
		url := blocklist.Target
		if !checkInvalidURLSymbolsUsage(url) {
			return fmt.Errorf("invalid blocklist target provided: %s", url)
		}
	}

	return nil
}

// checkInvalidURLSymbolsUsage checks for characters NOT allowed in URL
func checkInvalidURLSymbolsUsage(url string) bool {
	matched, err := regexp.MatchString("[^a-zA-Z0-9:/?&%=~._()-;]", url)
	if err != nil {
		return true
	}

	return !matched
}
