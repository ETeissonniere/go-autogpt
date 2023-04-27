package planners

import (
	"github.com/eteissonniere/hercules/llms"
	"github.com/eteissonniere/hercules/prompt"
	"github.com/eteissonniere/hercules/prompt/commands"
	"github.com/eteissonniere/hercules/prompt/executors"
)

type Planner interface {
	Plan(prompt.Task, []commands.Command) (prompt.Prompt, error)
	Convert(llms.ChatConversation) (executors.Plan, error)
}
