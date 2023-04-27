package planners

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"text/template"

	"github.com/eteissonniere/go-autogpt/llms"
	"github.com/eteissonniere/go-autogpt/prompt"
	"github.com/eteissonniere/go-autogpt/prompt/commands"
	"github.com/eteissonniere/go-autogpt/prompt/executors"
	"github.com/eteissonniere/go-autogpt/prompt/internal"
)

const basicTemplate = `Given a task defined as the following: "{{.Task}}" and the previously defined instructions, please create the outline of a plan to execute the previously described task.

Your plan should be written in the form of numeroted bullet points followed by a command to shutdown or terminate yourself. You may use the provided commands to do some preliminary research or further think by yourself. Only your last message will be taken to form the plan. Here is an example:
1. Do first part of the plan
2. Do something else
3. More stuff
4. etc...
shutdown

Each step of your plan should fit on exactly one line. Lines not conforming to the previous format will be ignored.`

type basicTemplateArguments struct {
	Task prompt.Task
}

type Basic struct{}

func NewBasic() *Basic {
	return &Basic{}
}

func (b *Basic) Plan(task prompt.Task, commands []commands.Command) (prompt.Prompt, error) {
	t := template.Must(template.New("basic").Parse(basicTemplate))
	p := bytes.NewBufferString("")
	t.Execute(p, basicTemplateArguments{Task: task})

	final := internal.NewBasePromptArguments(commands, p.String())
	return prompt.Prompt(final.Build()), nil
}

func (b *Basic) Convert(conversation llms.ChatConversation) (executors.Plan, error) {
	lastMessage := conversation[len(conversation)-1]
	if lastMessage.Role != llms.ChatRoleAssistant {
		return executors.Plan{}, fmt.Errorf("last message is not an assistant message")
	}

	pattern := regexp.MustCompile(`^\d+\. `)
	plan := executors.Plan{}
	for _, line := range strings.Split(lastMessage.Content, "\n") {
		if !pattern.MatchString(line) {
			continue
		}

		line = pattern.ReplaceAllString(line, "")
		plan = append(plan, line)
	}
	return plan, nil
}
