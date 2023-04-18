package hercules

type Command interface {
	Name() string
	Usage() string
	Execute(args []string) error
}
