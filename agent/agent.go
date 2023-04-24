package agent

import (
	"fmt"
	"strings"

	"github.com/eteissonniere/hercules/llms"
	"github.com/eteissonniere/hercules/prompt"
	"github.com/eteissonniere/hercules/prompt/commands"
	"github.com/rs/zerolog/log"
)

type Agent struct {
	prompt   prompt.Prompt
	llm      llms.LLMChatModel
	exporter Exporter
}

func New(prompt prompt.Prompt, llm llms.LLMChatModel, exporter Exporter) *Agent {
	return &Agent{prompt, llm, exporter}
}

func (a *Agent) Run() error {
	cmds := map[string]commands.Command{}
	for _, cmd := range a.prompt.Commands() {
		cmds[cmd.Name()] = cmd
	}

	conversation := llms.ChatConversation{
		{Role: llms.ChatRoleSystem, Content: a.prompt.String()},
	}
	a.onMessage(conversation[0])

	for {
		resp, err := a.llm.Complete(conversation)
		if err != nil {
			return fmt.Errorf("failed to get next message in conversation: %w", err)
		}
		a.onMessage(resp)

		lines := strings.Split(resp.Content, "\n")
		lastLine := lines[len(lines)-1]

		words := strings.Split(lastLine, " ")
		command := words[0]
		args := words[1:]

		reply := llms.ChatMessage{
			Role:    llms.ChatRoleSystem,
			Content: "command not found",
		}
		if cmd, ok := cmds[command]; ok {
			ret, err := cmd.Execute(args)
			if err == commands.ErrShutdown {
				return nil
			} else if aerr, ok := err.(*commands.AgentError); ok {
				reply.Content = aerr.AgentExplainer()
			} else if err != nil {
				return fmt.Errorf("failed to execute command: %w", err)
			} else {
				reply.Content = ret
			}
		}
		a.onMessage(reply)

		conversation = append(conversation, resp)
		conversation = append(conversation, reply)
	}
}

func (a *Agent) onMessage(msg llms.ChatMessage) error {
	log.Info().
		Str("role", string(msg.Role)).
		Msg(msg.Content)
	return a.exporter.Export(msg)
}
