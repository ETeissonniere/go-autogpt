package llms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const openAiUrl = "https://api.openai.com/v1/chat/completions"

type OpenAI struct {
	ApiKey string
	Model  string
}

func NewOpenAI(apiKey string, model string) *OpenAI {
	return &OpenAI{ApiKey: apiKey, Model: model}
}

func (o *OpenAI) Complete(conversation ChatConversation) (ChatMessage, error) {
	body, err := json.Marshal(map[string]interface{}{
		"model":    o.Model,
		"messages": conversation,
	})
	if err != nil {
		return ChatMessage{}, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	req, err := http.NewRequest("POST", openAiUrl, bytes.NewBuffer(body))
	if err != nil {
		return ChatMessage{}, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", o.ApiKey))

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
