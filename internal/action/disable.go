package action

import (
	"adless/internal/action/exit"
	"adless/internal/hostsfile"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func (a *Action) Disable(_ *cli.Context) error {
	hosts, err := hostsfile.New()
	if err != nil {
		return exit.Error(exit.HostsFile, err, "failed to process hosts file")
	}

	if hosts.Status() == hostsfile.Disabled {
		log.Info().Msg("domains blocking is already disabled")
		return nil
	}

	if err := hosts.Backup(); err != nil {
		return exit.Error(exit.HostsFile, err, "failed to backup hosts file")
	}

	if err := hosts.RemoveDomainsBlocking(); err != nil {
		return exit.Error(exit.HostsFile, err, "failed to disable domains blocking")
	}

	log.Info().Msg("domain blocking successfully disabled")

	return nil
}
