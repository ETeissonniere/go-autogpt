package commands

type Command interface {
	Name() string
	Usage() string
	Execute(args []string) (string, error)
}
