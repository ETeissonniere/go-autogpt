package agent

import (
	"fmt"
	"strings"

	"github.com/eteissonniere/hercules/llms"
	"github.com/eteissonniere/hercules/prompt"
	"github.com/eteissonniere/hercules/prompt/commands"
)

type Agent struct {
	prompt prompt.Prompt
	llm    llms.LLMChatModel
}

func New(prompt prompt.Prompt, llm llms.LLMChatModel) *Agent {
	return &Agent{prompt, llm}
}

func (a *Agent) Run() error {
	cmds := map[string]commands.Command{}
	for _, cmd := range a.prompt.Commands() {
		cmds[cmd.Name()] = cmd
	}

	conversation := llms.ChatConversation{
		{Role: llms.ChatRoleSystem, Content: a.prompt.String()},
	}
	logMessage(conversation[0])

	for {
		resp, err := a.llm.Complete(conversation)
		if err != nil {
			return fmt.Errorf("failed to get next message in conversation: %w", err)
		}
		logMessage(resp)

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
			} else if err != nil {
				return fmt.Errorf("failed to execute command: %w", err)
			}
			reply.Content = ret
		}
		logMessage(reply)

		conversation = append(conversation, resp)
		conversation = append(conversation, reply)
	}
}
