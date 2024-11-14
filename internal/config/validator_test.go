package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidate(t *testing.T) {
	t.Run("config is valid", func(t *testing.T) {
		config := &Config{
			Blocklists: []Blocklist{
				{Target: "https://raw.githubusercontent.com/FadeMind/hosts.extras/master/add.Spam/hosts"},
			},
		}

		validate := Validate(config)
		require.NoError(t, validate)
	})

	t.Run("config has no blocklists", func(t *testing.T) {
		config := &Config{}
		assert.ErrorIs(t, Validate(config), ErrNoBlocklistsProvided)
	})

	t.Run("config has invalid target", func(t *testing.T) {
		config := &Config{
			Blocklists: []Blocklist{
				{Target: "https://example.com/test page?query=value&extra#section+details"},
			},
		}

		assert.ErrorContains(t, Validate(config), "invalid blocklist target provided")
	})
}

func TestHasInvalidURLSymbols(t *testing.T) {
	t.Run("invalid URL symbols aren't used", func(t *testing.T) {
		testURL := "https://raw.githubusercontent.com/FadeMind/hosts.extras/master/add.Spam/hosts"
		assert.False(t, hasInvalidURLSymbols(testURL))
	})

	t.Run("invalid URL symbols are used", func(t *testing.T) {
		testURL := "https://example.com/test page?query=value&extra#section+details"
		assert.True(t, hasInvalidURLSymbols(testURL))
	})
}
