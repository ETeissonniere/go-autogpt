package agent

import (
	"github.com/avantgardists/hercules/llms"
	"github.com/rs/zerolog/log"
)

func logMessage(msg llms.ChatMessage) {
	switch msg.Role {
	case llms.ChatRoleAssistant:
		log.Info().Msgf("Assistant: %s", msg.Content)
	case llms.ChatRoleSystem:
		log.Info().Msgf("System: %s", msg.Content)
	case llms.ChatRoleUser:
		log.Info().Msgf("User: %s", msg.Content)
	}
}
