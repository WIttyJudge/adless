package config

import (
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	td, err := os.MkdirTemp("", "barrier-config")
	defer os.RemoveAll(td)

	require.NoError(t, err)

	t.Run("return default config if there is no local", func(t *testing.T) {
		os.Setenv("BARRIER_CONFIG_HOME", "test")
		defer os.Unsetenv("BARRIER_CONFIG_HOME")

		config, err := Load()

		assert.NotNil(t, config)
		assert.NoError(t, err)
		assert.Exactly(t, config, defaultConfig())
	})

	t.Run("returns errors if yml is invalid", func(t *testing.T) {
		testConfig, err := os.CreateTemp(td, "config.yml")
		require.NoError(t, err)
		defer os.Remove(testConfig.Name())

		_, err = testConfig.WriteString("invalid yml content")
		require.NoError(t, err)

		os.Setenv("BARRIER_CONFIG_PATH", testConfig.Name())
		defer os.Unsetenv("BARRIER_CONFIG_PATH")

		config, err := Load()

		assert.Nil(t, config)
		assert.ErrorContains(t, err, "yaml: unmarshal errors")
	})

	t.Run("return error if there are zero blocklists", func(t *testing.T) {
		testConfig, err := os.CreateTemp(td, "config.yml")
		require.NoError(t, err)
		defer os.Remove(testConfig.Name())

		_, err = testConfig.WriteString("blocklists:\n")
		require.NoError(t, err)

		os.Setenv("BARRIER_CONFIG_PATH", testConfig.Name())
		defer os.Unsetenv("BARRIER_CONFIG_PATH")

		config, err := Load()

		assert.Nil(t, config)
		assert.ErrorIs(t, err, ErrNoBlocklistsProvided)
	})

	t.Run("returns config successfully", func(t *testing.T) {
		testConfig, err := os.CreateTemp(td, "config.yml")
		require.NoError(t, err)
		defer os.Remove(testConfig.Name())

		_, err = testConfig.WriteString("blocklists:\n- target: https://test.com")
		require.NoError(t, err)

		os.Setenv("BARRIER_CONFIG_PATH", testConfig.Name())
		defer os.Unsetenv("BARRIER_CONFIG_PATH")

		config, err := Load()

		assert.NotNil(t, config)
		assert.NoError(t, err)
	})
}

func TestLoadByUser(t *testing.T) {
	td, err := os.MkdirTemp("", "barrier-config")
	defer os.RemoveAll(td)

	require.NoError(t, err)

	t.Run("config not found", func(t *testing.T) {
		config, err := LoadByUser(".test_config")

		assert.Nil(t, config)
		assert.NotNil(t, err)
		assert.ErrorIs(t, err, os.ErrNotExist)
	})

	t.Run("config is invalid", func(t *testing.T) {
		testConfig, err := os.CreateTemp(td, "")
		require.NoError(t, err)

		_, err = testConfig.WriteString("blocklists:\n")
		require.NoError(t, err)

		config, err := LoadByUser(testConfig.Name())

		assert.Nil(t, config)
		assert.ErrorIs(t, err, ErrNoBlocklistsProvided)
	})

	t.Run("config is loaded successfully", func(t *testing.T) {
		testConfig, err := os.CreateTemp(td, "")
		require.NoError(t, err)

		_, err = testConfig.WriteString("blocklists:\n- target: https://test.com")
		require.NoError(t, err)

		config, err := LoadByUser(testConfig.Name())

		assert.NotNil(t, config)
		assert.NoError(t, err, ErrNoBlocklistsProvided)
	})
}

func TestLocation(t *testing.T) {
	t.Run("BARRIER_CONFIG_PATH environment variable", func(t *testing.T) {
		bcp := path.Join(homeDir(), ".config", "test_barrier", "config.yml")

		os.Setenv("BARRIER_CONFIG_PATH", bcp)
		defer os.Unsetenv("BARRIER_CONFIG_PATH")

		assert.Equal(t, location(), bcp)
	})

	t.Run("BARRIER_CONFIG_HOME environment variable", func(t *testing.T) {
		bch := path.Join(homeDir(), ".config", "test_barrier")

		os.Setenv("BARRIER_CONFIG_HOME", bch)
		defer os.Unsetenv("BARRIER_CONFIG_HOME")

		expected := filepath.Join(bch, "config.yml")
		assert.Equal(t, location(), expected)
	})

	t.Run("XDG_CONFIG_HOME environment variable", func(t *testing.T) {
		xdgConfig := filepath.Join(homeDir(), ".test_config")

		os.Setenv("XDG_CONFIG_HOME", xdgConfig)
		defer os.Unsetenv("XDG_CONFIG_HOME")

		expected := filepath.Join(xdgConfig, "barrier", "config.yml")
		assert.Equal(t, location(), expected)
	})

	t.Run("default location", func(t *testing.T) {
		expected := filepath.Join(homeDir(), ".config", "barrier", "config.yml")
		assert.Equal(t, location(), expected)
	})
}

func TestRead(t *testing.T) {
	td, err := os.MkdirTemp("", "barrier-config")
	defer os.RemoveAll(td)

	require.NoError(t, err)

	t.Run("successfully read config", func(t *testing.T) {
		testConfig, err := os.CreateTemp(td, "")
		require.NoError(t, err)

		_, err = testConfig.WriteString("blocklists:\n- target: https://test.com")
		require.NoError(t, err)

		config, err := read(testConfig.Name())
		require.NoError(t, err)

		assert.Len(t, config.Blocklists, 1)
	})

	t.Run("config not found", func(t *testing.T) {
		config, err := read("test_config.yml")

		require.Nil(t, config)
		require.NotNil(t, err)

		assert.ErrorIs(t, err, os.ErrNotExist)
	})

	t.Run("failed to unmarshal yaml config", func(t *testing.T) {
		testConfig, err := os.CreateTemp(td, "")
		require.NoError(t, err)

		_, err = testConfig.WriteString("invalid yml content")
		require.NoError(t, err)

		config, err := read(testConfig.Name())
		require.Nil(t, config)
		require.NotNil(t, err)

		assert.ErrorContains(t, err, "yaml: unmarshal errors")
	})
}

func TestHomeDir(t *testing.T) {
	t.Run("SUDO_USER environment variable", func(t *testing.T) {
		os.Setenv("SUDO_USER", "root")
		defer os.Unsetenv("SUDO_USER")

		expected := "/root/.config/barrier/config.yml"
		assert.Equal(t, location(), expected)
	})
}

func TestDefaultConfig(t *testing.T) {
	t.Run("return default config", func(t *testing.T) {
		config := defaultConfig()

		require.NotNil(t, config)

		assert.IsType(t, config, &Config{})
	})
}
