package main

import (
	"flag"

	"github.com/eteissonniere/hercules/agent"
	"github.com/eteissonniere/hercules/llms"
	"github.com/eteissonniere/hercules/misc/logging"
	"github.com/eteissonniere/hercules/prompt"
	"github.com/eteissonniere/hercules/prompt/commands"

	"github.com/rs/zerolog/log"
)

func main() {
	// TODO: select model interface and model of agent
	// TODO: memory store
	// TODO: should be multi agent

	logging.Init(true)

	apiKey := flag.String("apiKey", "", "OpenAI API Key")
	name := flag.String("name", "", "Name of the agent")
	task := flag.String("task", "", "Task of the agent")

	flag.Parse()

	if *apiKey == "" || *name == "" || *task == "" {
		log.Fatal().Msg("missing required flag")
	}

	llm := llms.NewOpenAI(*apiKey, "gpt-3.5-turbo")
	agentPrompt := prompt.New(*name, *task, commands.DefaultCommands)
	agent := agent.New(agentPrompt, llm)

	if err := agent.Run(); err != nil {
		log.Error().Err(err).Msg("agent errored")
	}
}
