package main

import (
	"fmt"

	"github.com/eteissonniere/hercules/agents/executor"
	"github.com/eteissonniere/hercules/agents/mastermind"
	"github.com/eteissonniere/hercules/agents/planner"
	"github.com/eteissonniere/hercules/llms"
	"github.com/eteissonniere/hercules/misc/logging"
	"github.com/eteissonniere/hercules/prompt"
	"github.com/eteissonniere/hercules/prompt/commands"
	"github.com/eteissonniere/hercules/prompt/executors"
	"github.com/eteissonniere/hercules/prompt/planners"

	"github.com/urfave/cli/v2"
)

var cmdRun = cli.Command{
	Name:  "run",
	Usage: "Run the agent",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "apiKey",
			Usage:    "OpenAI API Key",
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
		exporter := logging.ExportToDebugLogs()
		if c.String("export") != "" {
			var err error
			exp, err := logging.ExportToFile(c.String("export"))
			if err != nil {
				return fmt.Errorf("failed to create file exporter: %w", err)
			}
			exporter = logging.ExportChain(
				exporter,
				exp,
			)
		}

		plannerPrompt := planners.NewBasic()
		executorPrompt := executors.NewBasic()

		llm := llms.NewOpenAI(c.String("apiKey"), "gpt-3.5-turbo")

		plannerAgent := planner.New(plannerPrompt, llm)
		executorAgent := executor.New(executorPrompt, llm)

		mastermindAgent := mastermind.New(plannerAgent, executorAgent)
		return mastermindAgent.Run(prompt.Task(c.String("task")), commands.DefaultCommands, exporter)
	},
}
