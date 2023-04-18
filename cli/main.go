package main

import (
	"github.com/avantgardists/hercules"
)

func main() {
	// TODO: select model interface and model of agent
	// TODO: memory store

	hercules.NewBrain(hercules.Config{
		Name:     "system",
		Task:     "system",
		Commands: []hercules.Command{},
	}).Start()
}
