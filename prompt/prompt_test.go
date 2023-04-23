package prompt

import (
	"testing"

	"github.com/eteissonniere/hercules/prompt/commands"

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

func TestCommands_sharesInitialCommands(t *testing.T) {
	commands := []commands.Command{
		&testCommand{name: "test", usage: "test usage"},
	}

	prompt := New("test", "test task", commands)
	assert.Equal(t, commands, prompt.Commands())
}
