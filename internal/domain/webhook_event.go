package domain

// WebhookEvent representa um evento de webhook
type WebhookEvent struct {
	Subscription string
	EventType    string
	InvoiceID    string
	Amount       int64
	Fee          int64
	Status       string
}

// WebhookService define a interface para processar webhooks
type WebhookService interface {
	ProcessEvent(event WebhookEvent) error
	ValidateSignature(body, signature string) bool
}
