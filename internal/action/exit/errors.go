package exit

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

const (
	Unknown = iota
	Config
	Hostsfile
)

// Error returns a user friendly CLI error.
func Error(exitCode int, err error, format string, args ...any) error {
	msg := fmt.Sprintf(format, args...)
	log.Error().Err(err).Msg(msg)

	return cli.Exit("", exitCode)
}
