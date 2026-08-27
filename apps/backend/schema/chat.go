package schema

import "github.com/lucassoaresfr/winthor-sql-agent.git/engine"

// Payload padronizado para requisições de Chat/Webhook
type ChatWebhookRequest struct {
	Event    string               `json:"event"`  // Ex: "message.create"
	Sender   string               `json:"sender"` // Ex: "5581999999999"
	Messages []engine.ChatMessage `json:"messages"`
}
