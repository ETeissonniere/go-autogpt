// The executor is an agent which makes use of a `Executor` prompt to execute
// a plan created by a planner agent.
package executor

import (
	"fmt"

	"github.com/eteissonniere/go-autogpt/agents/internal"
	"github.com/eteissonniere/go-autogpt/llms"
	"github.com/eteissonniere/go-autogpt/misc/logging"
	"github.com/eteissonniere/go-autogpt/prompt"
	"github.com/eteissonniere/go-autogpt/prompt/commands"
	"github.com/eteissonniere/go-autogpt/prompt/executors"
)

type Agent struct {
	prompter executors.Executor
	llm      llms.LLMChatModel
}

func New(prompter executors.Executor, llm llms.LLMChatModel) *Agent {
	return &Agent{prompter, llm}
}

func (a *Agent) Run(task prompt.Task, plan executors.Plan, commands []commands.Command, exporter logging.Exporter) error {
	prompt, err := a.prompter.Execute(task, plan, commands)
	if err != nil {
		return fmt.Errorf("failed to create execution prompt: %w", err)
	}

	_, err = internal.Evaluate(prompt, a.llm, commands, exporter)
	if err != nil {
		return fmt.Errorf("failed to evaluate: %w", err)
	}

	return nil
}
