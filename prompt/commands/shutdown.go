package commands

import "errors"

var ErrShutdown = errors.New("shutdown")

type ShutdownCommand struct{}

func (c *ShutdownCommand) Name() string {
	return "shutdown"
}

func (c *ShutdownCommand) Usage() string {
	return "terminate yourself along with the conversation thread, useful when you are done with the task"
}

func (c *ShutdownCommand) Execute(args []string) (string, error) {
	return "", ErrShutdown
}

func init() {
	DefaultCommands = append(DefaultCommands, &ShutdownCommand{})
}
