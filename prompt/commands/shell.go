package commands

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

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
	return "execute shell command with the provided arguments with a timeout of at most 1 second. Example: shell ls -la"
}

func (c *ExecuteShellCommand) Execute(args []string) (string, error) {
	output, err := c.Executor.Execute(strings.Join(args, " "))
	if err != nil {
		return "", NewAgentError(err)
	}

	return fmt.Sprintf("output: %s", output), nil
}

type ShellCommandExecutorWithNoGatekeeping struct{}

func (e *ShellCommandExecutorWithNoGatekeeping) Execute(command string) (string, error) {
	log.Debug().Str("command", command).Msg("executing shell command")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "sh", "-c", command)
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
