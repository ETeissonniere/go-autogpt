package prompt

import (
	"bytes"
	"text/template"

	"github.com/eteissonniere/hercules/prompt/commands"
)

const PromptTemplate = `You are {{.Name}}. Your task is to {{.Task}}.

You should accomplish your task autonomously. The user is not allowed
to and cannot interfere with your actions.

You can use the following commands:
{{range .Cmds}}{{.Name}}: {{.Usage}}
{{end}}
When replying, you can include any context, description or thoughts in
your answer. However, you must ensure that the last line of your
answer is the command you want to execute along with its arguments.

You are allowed only one command per reply.`

type Prompt struct {
	Name string
	Task string
	Cmds []struct {
		Name  string
		Usage string
	}

	commands []commands.Command
}

func New(name string, task string, commands []commands.Command) Prompt {
	vars := Prompt{
		Name: name,
		Task: task,
		Cmds: []struct {
			Name  string
			Usage string
		}{},
		commands: commands,
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

func (p Prompt) String() string {
	t := template.Must(template.New("prompt").Parse(PromptTemplate))
	prompt := bytes.NewBufferString("")
	t.Execute(prompt, p)
	return prompt.String()
}

func (p Prompt) Commands() []commands.Command {
	return p.commands
}
