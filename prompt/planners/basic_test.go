package planners

import (
	"testing"

	"github.com/eteissonniere/hercules/llms"
	"github.com/eteissonniere/hercules/prompt/executors"
	"github.com/stretchr/testify/assert"
)

func TestBasic_Convert(t *testing.T) {
	conversation := llms.ChatConversation{
		llms.ChatMessage{
			Role:    llms.ChatRoleSystem,
			Content: "This is an example system prompt",
		},
		llms.ChatMessage{
			Role:    llms.ChatRoleAssistant,
			Content: "This is an example assistant thought",
		},
		llms.ChatMessage{
			Role:    llms.ChatRoleSystem,
			Content: "This is an example system reply",
		},
		llms.ChatMessage{
			Role: llms.ChatRoleAssistant,
			Content: `This is an example assistant plan
1. Do first part of the plan
2. Do something else
3. More stuff
4. etc...

shutdown`,
		},
	}
	b := NewBasic()
	plan, err := b.Convert(conversation)
	assert.NoError(t, err)
	assert.Equal(t, executors.Plan{
		"Do first part of the plan",
		"Do something else",
		"More stuff",
		"etc...",
	}, plan)
}
