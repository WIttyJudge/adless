package action

import (
	"github.com/WIttyJudge/adless/internal/action/exit"
	"github.com/WIttyJudge/adless/internal/hostsfile"

	"github.com/rs/zerolog/log"

	"github.com/urfave/cli/v2"
)

func (a *Action) Restore(_ *cli.Context) error {
	log.Info().Msg("restoring hosts file from backup..")

	hosts, err := hostsfile.New()
	if err != nil {
		return exit.Error(exit.HostsFile, err, "failed to process hosts file")
	}

	if err := hosts.Restore(); err != nil {
		return exit.Error(exit.HostsFile, err, "failed to restore hosts file")
	}

	log.Info().Msg("restoring operation is successful")

	return nil
}
