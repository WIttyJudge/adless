package action

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func (a *Action) Stop(ctx *cli.Context) error {
	fmt.Println("Stop")

	return nil
}
