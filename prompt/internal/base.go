package internal

import (
	"bytes"
	"text/template"
	"time"

	"github.com/eteissonniere/go-autogpt/prompt/commands"
)

const BasePromptTemplate = `Today is {{.Today.Format "Jan 02, 2006 15:04:05 UTC"}}. You are an autonomous agent which is in charge of planning and later executing the provided tasks. To do so, you have been provided with a set of tools and commands for you to use:
{{range .Cmds}}{{.Name}}: {{.Usage}}
{{end}}
When replying, you can include any context, description or thoughts in your answer. However, you must ensure that the last line of your answer is the command you want to execute along with its arguments. A command should fit on exactly one line, if necessary, try to use escape characters if necessary. You cannot combine multiple commands in one message. Your command should start with the command name followed by its arguments:
command arg1 arg2 arg3

You are allowed only one command per reply. Your reply should always finish with a command. If you need to execute multiple commands do so with multiple different messages.

{{.AdditionalPrompt}}`

type BasePromptArguments struct {
	Cmds []struct {
		Name  string
		Usage string
	}
	Today            time.Time
	AdditionalPrompt string
}

func NewBasePromptArguments(commands []commands.Command, additionalPrompt string) BasePromptArguments {
	vars := BasePromptArguments{
		Cmds: []struct {
			Name  string
			Usage string
		}{},
		Today:            time.Now(),
		AdditionalPrompt: additionalPrompt,
	}
	for _, command := range commands {
		vars.Cmds = append(vars.Cmds, struct {
			Name  string
			Usage string
		}{
			Name:  command.Name(),
			Usage: command.Usage(),
		})
	}

	return vars
}

func (b BasePromptArguments) Build() string {
	t := template.Must(template.New("basePrompt").Parse(BasePromptTemplate))
	prompt := bytes.NewBufferString("")
	err := t.Execute(prompt, b)
	if err != nil {
		panic(err)
	}
	return prompt.String()
}
