package main

import (
	"barrier/internal/action"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/urfave/cli/v2"
)

var (
	version   string
	gitCommit string
	buildDate string
)

func main() {
	app := setupApp()

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func setupApp() *cli.App {
	action := action.New()

	app := cli.NewApp()
	app.Name = "barrier"
	app.Usage = "Local ad blocker writter in Go"
	app.UseShortOptionHandling = true
	app.Version = version

	cli.VersionPrinter = func(ctx *cli.Context) {
		fmt.Println("Version:\t", ctx.App.Version)
		fmt.Println("Git Commit:\t", gitCommit)
		fmt.Println("Build Date:\t", buildDate)
	}

	app.CommandNotFound = func(ctx *cli.Context, command string) {
		fmt.Printf("Error. Unknown command: '%s'\n\n", command)
		cli.ShowAppHelpAndExit(ctx, 1)
	}

	app.Before = action.BeforeAction
	app.Commands = action.GetCommands()
	app.Flags = action.GetFlags()

	sort.Sort(cli.CommandsByName(app.Commands))
	sort.Sort(cli.FlagsByName(app.Flags))

	return app
}
