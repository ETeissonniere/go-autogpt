package llms

type ChatRole string

const (
	ChatRoleAssistant ChatRole = "assistant"
	ChatRoleUser      ChatRole = "user"
	ChatRoleSystem    ChatRole = "system"
)

type ChatMessage struct {
	Role    ChatRole `json:"role"`
	Content string   `json:"content"`
}

type ChatConversation []ChatMessage

type LLMChatModel interface {
	Complete(ChatConversation) (ChatMessage, error)
}
