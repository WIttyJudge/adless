package action

import (
	"barrier/internal/config"
	"log"

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

	if err := a.loadConfig(ctx); err != nil {
		log.Fatal(err)
	}

	return nil
}

func (a *Action) GetCommands() []*cli.Command {
	return []*cli.Command{
		{
			Name:        "start",
			Usage:       "start blocking",
			Description: "",
			Action:      a.Start,
		},
		{
			Name:        "stop",
			Usage:       "stop blocking",
			Description: "",
			Action:      a.Stop,
		},
		{
			Name:        "update",
			Usage:       "update resources",
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
