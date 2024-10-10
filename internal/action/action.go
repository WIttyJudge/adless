package action

import (
	"barrier/internal/action/exit"
	"barrier/internal/config"

	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
)

type Action struct {
	config *config.Config
}

func New() *Action {
	return &Action{}
}

func (a *Action) BeforeAction(ctx *cli.Context) error {
	if ctx.NArg() == 0 {
		return nil
	}

	if ctx.Bool("verbose") {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	if err := a.loadConfig(ctx); err != nil {
		return exit.Error(exit.Config, err, "failed to load config file")
	}

	return nil
}

func (a *Action) GetCommands() []*cli.Command {
	return []*cli.Command{
		{
			Name:        "start",
			Usage:       "Start blocking",
			Description: "",
			Action:      a.Start,
		},
		{
			Name:        "stop",
			Usage:       "Stop blocking",
			Description: "",
			Action:      a.Stop,
		},
		{
			Name:        "update",
			Usage:       "Update resources",
			Description: "",
			Action:      a.Update,
		},
	}
}

func (a *Action) GetFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:  "config-file",
			Usage: "Path to config file",
		},
		&cli.BoolFlag{
			Name:               "verbose",
			Aliases:            []string{"v"},
			Usage:              "Enable debug mode",
			DisableDefaultText: true,
		},
	}
}

func (a *Action) loadConfig(ctx *cli.Context) error {
	var (
		cfg *config.Config
		err error
	)

	providedConfigPath := ctx.String("config-file")

	if providedConfigPath != "" {
		cfg, err = config.LoadByUser(providedConfigPath)
	} else {
		cfg, err = config.Load()
	}

	if err != nil {
		return err
	}

	a.config = cfg

	return nil
}
