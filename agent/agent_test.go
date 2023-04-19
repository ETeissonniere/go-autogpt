package agent

import (
	"fmt"
	"testing"

	"github.com/avantgardists/hercules/llms"
	"github.com/avantgardists/hercules/prompt"
	"github.com/avantgardists/hercules/prompt/commands"
	"github.com/stretchr/testify/assert"
)

var errMockCommand = fmt.Errorf("mock error")

type mockErrorCommand struct{}

func (c *mockErrorCommand) Name() string {
	return "error"
}

func (c *mockErrorCommand) Usage() string {
	return "error"
}

func (c *mockErrorCommand) Execute(args []string) (string, error) {
	return "", errMockCommand
}

type mockChatModel string

func (m mockChatModel) Complete(conversation llms.ChatConversation) (llms.ChatMessage, error) {
	return llms.ChatMessage{
		Role:    llms.ChatRoleAssistant,
		Content: string(m),
	}, nil
}

func testPrompt() prompt.Prompt {
	return prompt.New("test", "test", []commands.Command{
		&commands.DoNothingCommand{},
		&commands.ShutdownCommand{},
		&mockErrorCommand{},
	})
}

type mockChatModelWithTest func(conversation llms.ChatConversation) (llms.ChatMessage, error)

func (m mockChatModelWithTest) Complete(conversation llms.ChatConversation) (llms.ChatMessage, error) {
	return m(conversation)
}

func TestRun_llmError(t *testing.T) {
	a := New(testPrompt(), mockChatModelWithTest(func(conversation llms.ChatConversation) (llms.ChatMessage, error) {
		return llms.ChatMessage{}, fmt.Errorf("mock error")
	}))
	assert.Equal(t, fmt.Errorf("failed to get next message in conversation: %w", fmt.Errorf("mock error")), a.Run())
}

func TestRun_exitOnShutdown(t *testing.T) {
	a := New(testPrompt(), mockChatModel((&commands.ShutdownCommand{}).Name()))
	assert.Nil(t, a.Run())
}

func TestRun_forwardError(t *testing.T) {
	a := New(testPrompt(), mockChatModel((&mockErrorCommand{}).Name()))
	assert.Equal(t, fmt.Errorf("failed to execute command: %w", errMockCommand), a.Run())
}

func TestRun_handleCommandNotFound(t *testing.T) {
	called := false
	a := New(testPrompt(), mockChatModelWithTest(func(conversation llms.ChatConversation) (llms.ChatMessage, error) {
		if !called {
			called = true
			return llms.ChatMessage{
				Role:    llms.ChatRoleAssistant,
				Content: "not a command",
			}, nil
		} else {
			assert.Equal(t, "command not found", conversation[len(conversation)-1].Content)
			return llms.ChatMessage{}, fmt.Errorf("done")
		}
	}))
	assert.Error(t, a.Run())
}
