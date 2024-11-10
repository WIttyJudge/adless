package action

import (
	"barrier/internal/action/exit"
	"barrier/internal/hostsfile"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func (a *Action) Enable(ctx *cli.Context) error {
	hosts, err := hostsfile.New()
	if err != nil {
		return exit.Error(exit.HostsFile, err, "failed to process hosts file")
	}

	if hosts.Status() == hostsfile.Enabled {
		log.Info().Msg("domains blocking is already enabled")
		return nil
	}

	if err := hosts.Backup(); err != nil {
		return exit.Error(exit.HostsFile, err, "failed to backup hosts file")
	}

	processor := hostsfile.NewProcessor(a.config)
	parsedBlocklists, err := processor.Process()
	if err != nil {
		return err
	}

	if err := hosts.Write(parsedBlocklists.FormatToHostsfile()); err != nil {
		return exit.Error(exit.HostsFile, err, "failed to write to hosts file")
	}

	log.Info().Msg("domain blocking successfully enabled")

	return nil
}
