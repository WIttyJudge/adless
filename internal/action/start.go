package action

import (
	"barrier/internal/http"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func (a *Action) Start(ctx *cli.Context) error {
	httpClient := http.New(a.config)

	for _, blocklist := range a.config.Blocklists {
		target := blocklist.Target

		log.Info().Str("target", target).Msg("Processing blocklist..")

		resp, err := httpClient.Get(target)
		if err != nil {
			return err
		}

		fmt.Println(resp)
	}

	return nil
}
