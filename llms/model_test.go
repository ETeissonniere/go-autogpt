package llms

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func openAiApiKey() string {
	f, err := os.Open("../.secrets/openai_sk")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	buf := make([]byte, 1024)
	n, err := f.Read(buf)
	if err != nil && err != io.EOF {
		panic(err)
	}

	return string(buf[:n])
}

func TestModels(t *testing.T) {
	models := []LLMChatModel{
		NewOpenAI(openAiApiKey(), &GPT3Point5Turbo{}),
	}

	conversation := ChatConversation{
		{
			Role:    ChatRoleSystem,
			Content: "You are a friendly chatbot. Reply to the user.",
		},
		{
			Role:    ChatRoleUser,
			Content: "Hello",
		},
		{
			Role:    ChatRoleAssistant,
			Content: "Hi",
		},
		{
			Role:    ChatRoleUser,
			Content: "How are you?",
		},
	}

	for _, model := range models {
		msg, err := model.Complete(conversation)
		if err != nil {
			assert.NoError(t, err, "failed to complete conversation")
		}

		t.Logf("model: %s", msg.Content)
	}
}
