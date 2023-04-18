package hercules

type Brain struct {
	config Config
}

func NewBrain(config Config) *Brain {
	return &Brain{config: config}
}

func (b *Brain) Start() {
}
