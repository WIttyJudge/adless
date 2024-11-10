package main

import (
	"barrier/internal/action"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

var (
	version   string
	gitCommit string
	buildDate string
)

func main() {
	setupLogger()
	app := setupApp()

	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(err)
	}
}

func setupLogger() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.TimeOnly})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}

func setupApp() *cli.App {
	action := action.New()

	app := cli.NewApp()
	app.Name = "barrier"
	app.Usage = "Local ad blocker writter in Go"
	app.UseShortOptionHandling = true
	app.Version = version

	cli.VersionFlag = &cli.BoolFlag{
		Name:               "version",
		Aliases:            []string{"V"},
		Usage:              "Print the version",
		DisableDefaultText: true,
	}

	cli.HelpFlag = &cli.BoolFlag{
		Name:               "help",
		Aliases:            []string{"h"},
		Usage:              "Show help",
		DisableDefaultText: true,
	}

	cli.VersionPrinter = func(ctx *cli.Context) {
		fmt.Println("Version:\t", ctx.App.Version)
		fmt.Println("Git Commit:\t", gitCommit)
		fmt.Println("Build Date:\t", buildDate)
	}

	app.CommandNotFound = func(ctx *cli.Context, command string) {
		fmt.Printf("error: unrecognized command: '%s'\n\n", command)
		fmt.Println("for more information, try '--help'.")
	}

	app.Before = action.BeforeAction
	app.Commands = action.GetCommands()
	app.Flags = action.GetFlags()

	sort.Sort(cli.CommandsByName(app.Commands))
	sort.Sort(cli.FlagsByName(app.Flags))

	return app
}
