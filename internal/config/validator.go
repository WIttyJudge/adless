package config

import (
	"errors"
	"fmt"
	"regexp"
)

var ErrNoBlocklistsProvided = errors.New("no blocklists provided")

func Validate(config *Config) error {
	if len(config.Blocklists) == 0 {
		return ErrNoBlocklistsProvided
	}

	for _, blocklist := range config.Blocklists {
		url := blocklist.Target
		if hasInvalidURLSymbols(url) {
			return fmt.Errorf("invalid blocklist target provided: %s", url)
		}
	}

	return nil
}

// hasInvalidURLSymbols checks for characters NOT allowed in URL.
func hasInvalidURLSymbols(url string) bool {
	matched, _ := regexp.MatchString("[^a-zA-Z0-9:/?&%=~._()-;]", url)
	return matched
}
