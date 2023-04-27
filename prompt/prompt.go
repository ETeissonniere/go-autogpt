package prompt

type Prompt string
type Task string

func (p Prompt) String() string {
	return string(p)
}
