package prompt

import (
	"testing"

	"github.com/avantgardists/hercules/prompt/commands"
	"github.com/stretchr/testify/assert"
)

type testCommand struct {
	name  string
	usage string
}

func (c *testCommand) Name() string {
	return c.name
}

func (c *testCommand) Usage() string {
	return c.usage
}

func (c *testCommand) Execute(args []string) (string, error) {
	return "", nil
}

func TestString_generatePrompt(t *testing.T) {
	commands := []commands.Command{
		&testCommand{name: "test1", usage: "test usage"},
		&testCommand{name: "test2", usage: "test usage"},
	}

	prompt := New("test", "test task", commands)
	expected := `You are test. Your task is to test task.

You should accomplish your task autonomously. The user is not allowed
to and cannot interfere with your actions.

You can use the following commands:
test1: test usage
test2: test usage

When replying, you can include any context, description or thoughts in
your answer. However, you must ensure that the last line of your
answer is the command you want to execute along with its arguments.

You are allowed only one command per reply.`

	assert.Equal(t, expected, prompt.String())
}

func TestCommands_sharesInitialCommands(t *testing.T) {
	commands := []commands.Command{
		&testCommand{name: "test", usage: "test usage"},
	}

	prompt := New("test", "test task", commands)
	assert.Equal(t, commands, prompt.Commands())
}
