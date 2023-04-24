package prompt

import (
	"bytes"
	"text/template"
	"time"

	"github.com/eteissonniere/hercules/prompt/commands"
	"github.com/rs/zerolog/log"
)

const PromptTemplate = `Today is {{.Today.Format "Jan 02, 2006 15:04:05 UTC"}}. You are {{.Name}}. Your task follows:
{{.Task}}

You should accomplish your task autonomously. The user is not allowed to and cannot interfere with your actions. To do so, you can use the following commands:
{{range .Cmds}}{{.Name}}: {{.Usage}}
{{end}}
When replying, you can include any context, description or thoughts in your answer. However, you must ensure that the last line of your answer is the command you want to execute along with its arguments. A command should fit on exactly one line, if necessary, try to use escape characters if necessary. You cannot combine multiple commands in one message. Your command should start with the command name followed by its arguments:
command arg1 arg2 arg3

You are allowed only one command per reply. Your reply should always finish with a command.`

type Prompt struct {
	Name string
	Task string
	Cmds []struct {
		Name  string
		Usage string
	}
	Today time.Time

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
		Today:    time.Now(),
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

	log.Debug().
		Str("name", name).
		Str("task", task).
		Interface("commands", vars.Cmds).
		Msg("new prompt created")

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
