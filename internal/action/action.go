package action

import "github.com/urfave/cli/v2"

type Action struct{}

func New() *Action {
	return &Action{}
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
	return []cli.Flag{}
}
