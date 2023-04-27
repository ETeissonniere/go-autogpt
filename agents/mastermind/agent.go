// A special agent which uses multiple sub-agents such as a planner and an
// executor to achieve a specified task.
package mastermind

import (
	"fmt"

	"github.com/eteissonniere/hercules/misc/logging"
	"github.com/eteissonniere/hercules/prompt"
	"github.com/eteissonniere/hercules/prompt/commands"
	"github.com/eteissonniere/hercules/prompt/executors"
	"github.com/rs/zerolog/log"
)

type Planner interface {
	Run(prompt.Task, []commands.Command, logging.Exporter) (executors.Plan, error)
}

type Executor interface {
	Run(prompt.Task, executors.Plan, []commands.Command, logging.Exporter) error
}

type Agent struct {
	planner  Planner
	executor Executor
}

func New(planner Planner, executor Executor) *Agent {
	return &Agent{planner, executor}
}

func (a *Agent) Run(task prompt.Task, commands []commands.Command, exporter logging.Exporter) error {
	plan, err := a.planner.Run(task, commands, exporter)
	if err != nil {
		return fmt.Errorf("failed to get plan: %w", err)
	}

	log.Info().
		Str("task", string(task)).
		Interface("plan", plan).
		Msg("plan generated")

	err = a.executor.Run(task, plan, commands, exporter)
	if err != nil {
		return fmt.Errorf("failed to execute plan: %w", err)
	}

	return nil
}
