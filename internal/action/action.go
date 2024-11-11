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

	if ctx.Bool("quite") {
		zerolog.SetGlobalLevel(zerolog.Disabled)
	}

	if err := a.loadConfig(ctx); err != nil {
		return exit.Error(exit.Config, err, "failed to load config file")
	}

	return nil
}

func (a *Action) GetCommands() []*cli.Command {
	return []*cli.Command{
		{
			Name:   "disable",
			Usage:  "Disable domains blocking",
			Action: a.Disable,
		},
		{
			Name:   "enable",
			Usage:  "Enable domains blocking",
			Action: a.Enable,
		},

		{
			Name:   "status",
			Usage:  "Check if domains blocking enabled or not",
			Action: a.Status,
		},
		{
			Name:   "update",
			Usage:  "Update the list of domains to be blocked",
			Action: a.Update,
		},
		{
			Name:  "restore",
			Usage: "Restore hosts file from backup to its previous state",
			Description: "" +
				"When a `enable`, `disable` or `update` command is invoked, it creates a backup of the " +
				"original hosts file by copying it a backup file (hosts.backup).\n" +
				"The `restore` command copies the backup file (hosts.backup) back to its " +
				"original location (hosts).\n" +
				"Backup file must already exist to perform a command successfully.",
			Action: a.Restore,
		},
	}
}

func (a *Action) GetFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:  "config-file",
			Usage: "Path to the configuration file",
		},
		&cli.BoolFlag{
			Name:               "verbose",
			Aliases:            []string{"v"},
			Usage:              "Enable debug mode",
			DisableDefaultText: true,
		},
		&cli.BoolFlag{
			Name:               "quite",
			Aliases:            []string{"q"},
			Usage:              "Enable quiet mode",
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
