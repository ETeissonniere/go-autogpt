package executors

import (
	"bytes"
	"html/template"

	"github.com/eteissonniere/go-autogpt/prompt"
	"github.com/eteissonniere/go-autogpt/prompt/commands"
	"github.com/eteissonniere/go-autogpt/prompt/internal"
)

const basicTemplate = `You have been provided with the following task: "{{.Task}}". Please execute it using the previously defined commands andc terminate yourself when you deem it appropriate.

Prior agent instances and researches defined that the best plan to execute this task is the following:
{{range $index, $element := .Plan}}{{$index}}. {{$element}}
{{end}}`

type basicTemplateArguments struct {
	Task prompt.Task
	Plan Plan
}

type Basic struct{}

func NewBasic() *Basic {
	return &Basic{}
}

func (b *Basic) Execute(task prompt.Task, plan Plan, commands []commands.Command) (prompt.Prompt, error) {
	t := template.Must(template.New("basic").Parse(basicTemplate))
	p := bytes.NewBufferString("")
	t.Execute(p, basicTemplateArguments{Task: task, Plan: plan})

	final := internal.NewBasePromptArguments(commands, p.String())
	return prompt.Prompt(final.Build()), nil
}
