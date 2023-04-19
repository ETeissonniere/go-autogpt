package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/eteissonniere/hercules/agent"
	"github.com/eteissonniere/hercules/llms"
	"github.com/eteissonniere/hercules/misc/logging"
	"github.com/eteissonniere/hercules/prompt"
	"github.com/eteissonniere/hercules/prompt/commands"

	"github.com/rs/zerolog/log"
)

func askUserInput(prompt string) string {
	fmt.Print(prompt)

	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')

	text = text[:len(text)-1]
	return text
}

func main() {
	logging.Init(true)

	// TODO: select model interface and model of agent
	// TODO: memory store
	// TODO: should be multi agent

	apiKey := askUserInput("OpenAI API Key: ")
	name := askUserInput("Name: ")
	task := askUserInput("Task: ")

	fmt.Println(task)

	llm := llms.NewOpenAI(apiKey, "gpt-3.5-turbo")
	agentPrompt := prompt.New(name, task, commands.DefaultCommands)
	agent := agent.New(agentPrompt, llm)

	if err := agent.Run(); err != nil {
		log.Error().Err(err).Msg("agent errored")
	}
}
