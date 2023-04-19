package agent

import (
	"github.com/eteissonniere/hercules/llms"

	"github.com/rs/zerolog/log"
)

func logMessage(msg llms.ChatMessage) {
	log.Info().
		Str("role", string(msg.Role)).
		Msg(msg.Content)
}
