package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockExecutor struct {
	err         error
	out         string
	executeArgs string
}

func (e *mockExecutor) Execute(command string) (string, error) {
	e.executeArgs = command
	return e.out, e.err
}

func TestExecuteShellCommand_passToExecutor(t *testing.T) {
	mock := &mockExecutor{err: nil}

	command := &ExecuteShellCommand{Executor: mock}
	_, err := command.Execute([]string{"ls", "-la"})
	assert.Nil(t, err)

	assert.Equal(t, "ls -la", mock.executeArgs)
}

func TestExecuteShellCommand_forwardError(t *testing.T) {
	mock := &mockExecutor{err: assert.AnError}

	command := &ExecuteShellCommand{Executor: mock}
	_, err := command.Execute([]string{"ls", "-la"})
	assert.Equal(t, NewAgentError(assert.AnError), err)
}

func TestExecuteShellCommand_forwardOutput(t *testing.T) {
	mock := &mockExecutor{err: nil, out: "output"}

	command := &ExecuteShellCommand{Executor: mock}
	out, err := command.Execute([]string{"ls", "-la"})
	assert.Nil(t, err)
	assert.Equal(t, "output: output", out)
}
