package internal

import (
	"fmt"
	"strings"

	"github.com/eteissonniere/go-autogpt/agents/internal/helpers"
	"github.com/eteissonniere/go-autogpt/llms"
	"github.com/eteissonniere/go-autogpt/misc/logging"
	"github.com/eteissonniere/go-autogpt/prompt"
	"github.com/eteissonniere/go-autogpt/prompt/commands"
	"github.com/rs/zerolog/log"
)

func Evaluate(prompt prompt.Prompt, llm llms.LLMChatModel, cmdsList []commands.Command, exporter logging.Exporter) (llms.ChatConversation, error) {
	cmds := map[string]commands.Command{}
	for _, cmd := range cmdsList {
		cmds[cmd.Name()] = cmd
		log.Debug().Str("command", cmd.Name()).Msg("registered command")
	}

	conversation := llms.ChatConversation{
		{Role: llms.ChatRoleSystem, Content: prompt.String()},
	}
	exporter.Export(conversation[0])

	for {
		resp, err := llm.Complete(conversation)
		if err != nil {
			return nil, fmt.Errorf("failed to get next message in conversation: %w", err)
		}
		exporter.Export(resp)
		conversation = append(conversation, resp)

		lines := strings.Split(resp.Content, "\n")
		lastLine := lines[len(lines)-1]
		words := helpers.SplitCommand(helpers.WithEscapeCharacters(lastLine))
		command := words[0]
		args := words[1:]

		reply := llms.ChatMessage{
			Role:    llms.ChatRoleSystem,
			Content: "command not found",
		}
		if cmd, ok := cmds[command]; ok {
			ret, err := cmd.Execute(args)
			if err == commands.ErrShutdown {
				return conversation, nil
			} else if aerr, ok := err.(*commands.AgentError); ok {
				reply.Content = aerr.AgentExplainer()
			} else if err != nil {
				return nil, fmt.Errorf("failed to execute command: %w", err)
			} else {
				reply.Content = ret
			}
		}

		conversation = append(conversation, reply)
		exporter.Export(reply)
	}
}
