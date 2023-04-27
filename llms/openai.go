package llms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/tiktoken-go/tokenizer"
)

const openAiUrl = "https://api.openai.com/v1/chat/completions"

type OpenAI struct {
	apiKey    string
	model     OpenAiModel
	tokenizer tokenizer.Codec
}

type OpenAiModel interface {
	Name() string
	Tokenizer() tokenizer.Codec
	TokenLimit() int
}

type GPT3Point5Turbo struct{}

func (m *GPT3Point5Turbo) Name() string {
	return "gpt-3.5-turbo"
}

func (m *GPT3Point5Turbo) Tokenizer() tokenizer.Codec {
	enc, err := tokenizer.ForModel(tokenizer.GPT35Turbo)
	if err != nil {
		panic(err)
	}

	return enc
}

func (m *GPT3Point5Turbo) TokenLimit() int {
	// max is 4096, we need to keep some for reply
	return 2048
}

func NewOpenAI(apiKey string, model OpenAiModel) *OpenAI {
	return &OpenAI{apiKey: apiKey, model: model, tokenizer: model.Tokenizer()}
}

func (o *OpenAI) nbTokens(content string) (int, error) {
	ids, _, err := o.tokenizer.Encode(content)
	if err != nil {
		return 0, fmt.Errorf("failed to compute tokens: %w", err)
	}

	return len(ids), nil
}

func (o *OpenAI) Complete(conversation ChatConversation) (ChatMessage, error) {
	strippedConversation := ChatConversation{conversation[0]}
	tokensUsed, err := o.nbTokens(conversation[0].Content)
	if err != nil {
		return ChatMessage{}, fmt.Errorf("failed to compute tokens of message: %w", err)
	}
	if tokensUsed > o.model.TokenLimit() {
		return ChatMessage{}, fmt.Errorf("first message is too long")
	}

	// start by the end, so we can stop when we reach the token limit
	for i := len(conversation) - 1; i >= 1; i-- {
		msg := conversation[i]
		tokens, err := o.nbTokens(msg.Content)
		if err != nil {
			return ChatMessage{}, fmt.Errorf("failed to compute tokens of message: %w", err)
		}

		if tokensUsed+tokens > o.model.TokenLimit() {
			break
		}

		tokensUsed += tokens
		strippedConversation = append(ChatConversation{msg}, strippedConversation...)
	}

	log.Debug().
		Int("tokens", tokensUsed).
		Int("limit", o.model.TokenLimit()).
		Int("messages", len(strippedConversation)).
		Msg("prepared stripped conversation")

	body, err := json.Marshal(map[string]interface{}{
		"model":    o.model.Name(),
		"messages": strippedConversation,
	})
	if err != nil {
		return ChatMessage{}, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	req, err := http.NewRequest("POST", openAiUrl, bytes.NewBuffer(body))
	if err != nil {
		return ChatMessage{}, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", o.apiKey))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ChatMessage{}, fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ChatMessage{}, fmt.Errorf("OpenAI returned status code %d", resp.StatusCode)
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return ChatMessage{}, fmt.Errorf("failed to decode JSON: %w", err)
	}

	if len(result.Choices) > 0 {
		return ChatMessage{
			Role:    ChatRoleAssistant,
			Content: result.Choices[0].Message.Content,
		}, nil
	}

	return ChatMessage{}, nil
}
