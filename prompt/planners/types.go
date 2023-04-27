package planners

import (
	"github.com/eteissonniere/go-autogpt/llms"
	"github.com/eteissonniere/go-autogpt/prompt"
	"github.com/eteissonniere/go-autogpt/prompt/commands"
	"github.com/eteissonniere/go-autogpt/prompt/executors"
)

type Planner interface {
	Plan(prompt.Task, []commands.Command) (prompt.Prompt, error)
	Convert(llms.ChatConversation) (executors.Plan, error)
}
