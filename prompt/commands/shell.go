package commands

import (
	"os/exec"
	"strings"

	"github.com/rs/zerolog/log"
)

type ShellCommandExecutor interface {
	Execute(command string) (string, error)
}

type ExecuteShellCommand struct {
	Executor ShellCommandExecutor
}

func (c *ExecuteShellCommand) Name() string {
	return "shell"
}

func (c *ExecuteShellCommand) Usage() string {
	return "execute shell command with the provided arguments. Example: shell ls -la"
}

func (c *ExecuteShellCommand) Execute(args []string) (string, error) {
	return c.Executor.Execute(strings.Join(args, " "))
}

type ShellCommandExecutorWithNoGatekeeping struct{}

func (e *ShellCommandExecutorWithNoGatekeeping) Execute(command string) (string, error) {
	log.Debug().Str("command", command).Msg("executing shell command")

	cmd := exec.Command("sh", "-c", command)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	log.Debug().Str("output", string(output)).Msg("shell command output")

	return string(output), nil
}

func init() {
	DefaultCommands = append(DefaultCommands, &ExecuteShellCommand{Executor: &ShellCommandExecutorWithNoGatekeeping{}})
}
