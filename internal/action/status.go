package action

import (
	"barrier/internal/action/exit"
	"barrier/internal/hostsfile"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func (a *Action) Status(ctx *cli.Context) error {
	hosts, err := hostsfile.New()
	if err != nil {
		return exit.Error(exit.HostsFile, err, "failed to process hosts file")
	}

	status := hosts.Status()
	if status == hostsfile.Enabled {
		log.Info().Msg("domains blocking enabled")
		return nil
	}

	log.Info().Msg("domains blocking disabled")

	return nil
}
