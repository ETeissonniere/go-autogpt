package prompt

import (
	"bytes"
	"text/template"

	"github.com/avantgardists/hercules/prompt/commands"
)

const PromptTemplate = `You are {{.Name}}. Your task is to {{.Task}}.

You should accomplish your task autonomously. The user is not allowed
to and cannot interfere with your actions.

You can use the following commands:
{{range .Commands}}{{.Name}}: {{.Usage}}{{end}}

When replying, you can include any context, description or thoughts in
your answer. However, you must ensure that the last line of your
answer is the command you want to execute along with its arguments.
`

type promptVariables struct {
	Name     string
	Task     string
	Commands []struct {
		Name  string
		Usage string
	}
}

func New(name string, task string, commands []commands.Command) string {
	vars := promptVariables{
		Name: name,
		Task: task,
	}
	for _, command := range commands {
		vars.Commands = append(vars.Commands, struct {
			Name  string
			Usage string
		}{
			Name:  command.Name(),
			Usage: command.Usage(),
		})
	}

	t := template.Must(template.New("prompt").Parse(PromptTemplate))
	prompt := bytes.NewBufferString("")
	t.Execute(prompt, vars)
	return prompt.String()
}
