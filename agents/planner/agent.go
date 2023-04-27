// The planner agent make use of a `Planenr` prompt to create an execution
// plan which can then be passed to another agent for execution.
package planner

import (
	"fmt"

	"github.com/eteissonniere/hercules/agents/internal"
	"github.com/eteissonniere/hercules/llms"
	"github.com/eteissonniere/hercules/misc/logging"
	"github.com/eteissonniere/hercules/prompt"
	"github.com/eteissonniere/hercules/prompt/commands"
	"github.com/eteissonniere/hercules/prompt/executors"
	"github.com/eteissonniere/hercules/prompt/planners"
)

type Agent struct {
	prompter planners.Planner
	llm      llms.LLMChatModel
}

func New(prompter planners.Planner, llm llms.LLMChatModel) *Agent {
	return &Agent{prompter, llm}
}

func (a *Agent) Run(task prompt.Task, commands []commands.Command, exporter logging.Exporter) (executors.Plan, error) {
	prompt, err := a.prompter.Plan(task, commands)
	if err != nil {
		return nil, fmt.Errorf("failed to plan for task: %w", err)
	}

	conversation, err := internal.Evaluate(prompt, a.llm, commands, exporter)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate: %w", err)
	}

	plan, err := a.prompter.Convert(conversation)
	if err != nil {
		return nil, fmt.Errorf("failed to convert conversation to plan: %w", err)
	}

	return plan, nil
}
