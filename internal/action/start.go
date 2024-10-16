package action

import (
	"barrier/internal/action/exit"
	"barrier/internal/hostsfile"
	"fmt"

	"github.com/urfave/cli/v2"
)

func (a *Action) Start(ctx *cli.Context) error {
	processor := hostsfile.NewProcessor(a.config)

	hosts, err := hostsfile.New()
	if err != nil {
		return exit.Error(exit.Hostsfile, err, "failed to process hosts file")
	}

	if err := hosts.Backup(); err != nil {
		return exit.Error(exit.Hostsfile, err, "failed to backup hosts file")
	}

	result, err := processor.Process()
	if err != nil {
		return err
	}

	fmt.Println(result.FormatToHostsfile())

	return nil
}
