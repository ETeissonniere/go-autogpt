package executors

import (
	"github.com/eteissonniere/go-autogpt/prompt"
	"github.com/eteissonniere/go-autogpt/prompt/commands"
)

// A list of natural languages steps which describe a plan.
type Plan []string

type Executor interface {
	Execute(prompt.Task, Plan, []commands.Command) (prompt.Prompt, error)
}
