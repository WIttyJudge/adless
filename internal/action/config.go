package action

import (
	"adless/internal/action/exit"
	"adless/internal/config"

	"github.com/urfave/cli/v2"
)

func (a *Action) Config(_ *cli.Context) error {
	a.config.Print()

	return nil
}

func (a *Action) ConfigInit(_ *cli.Context) error {
	if err := config.Init(); err != nil {
		return exit.Error(exit.Config, err, "failed to initialize the config file")
	}

	return nil
}

func (a *Action) ConfigEdit(_ *cli.Context) error {
	if err := config.Init(); err != nil {
		return exit.Error(exit.Config, err, "failed to edit the config file")
	}

	if err := config.Edit(); err != nil {
		return exit.Error(exit.Config, err, "failed to edit the config file")
	}

	return nil
}
