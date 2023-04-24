package main

import (
	"os"

	"github.com/eteissonniere/hercules/misc/logging"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func main() {
	// TODO: select model interface and model of agent
	// TODO: memory store
	// TODO: should be multi agent

	logging.Init(true)

	app := &cli.App{
		Name: os.Args[0],
		Commands: []*cli.Command{
			&cmdRun,
			&cmdReplay,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("failed to run app")
	}
}
