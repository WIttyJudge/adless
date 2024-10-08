package action

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func (a *Action) Start(ctx *cli.Context) error {
	fmt.Println("Start")

	return nil
}
