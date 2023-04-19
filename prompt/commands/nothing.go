package commands

type DoNothingCommand struct{}

func (c *DoNothingCommand) Name() string {
	return "nothing"
}

func (c *DoNothingCommand) Usage() string {
	return "do nothing"
}

func (c *DoNothingCommand) Execute(args []string) (string, error) {
	return "nothing done", nil
}

func init() {
	DefaultCommands = append(DefaultCommands, &DoNothingCommand{})
}
