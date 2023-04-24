package main

import (
	"os"

	"github.com/eteissonniere/hercules/agent"
	"github.com/eteissonniere/hercules/llms"
	"github.com/eteissonniere/hercules/misc/logging"
	"github.com/eteissonniere/hercules/prompt"
	"github.com/eteissonniere/hercules/prompt/commands"

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
			{
				Name:  "run",
				Usage: "Run the agent",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "apiKey",
						Usage:    "OpenAI API Key",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "name",
						Usage:    "Name of the agent",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "task",
						Usage:    "Task of the agent",
						Required: true,
					},
					&cli.StringFlag{
						Name:  "export",
						Usage: "Path to export the agent logs to - ignored if not specified",
					},
				},
				Action: func(c *cli.Context) error {
					llm := llms.NewOpenAI(c.String("apiKey"), "gpt-3.5-turbo")
					exporter := agent.DoNotExport()
					if c.String("export") != "" {
						var err error
						exporter, err = agent.ExportToFile(c.String("export"))
						if err != nil {
							return err
						}
					}
					agentPrompt := prompt.New(c.String("name"), c.String("task"), commands.DefaultCommands)
					agent := agent.New(agentPrompt, llm, exporter)

					return agent.Run()
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("failed to run app")
	}
}
