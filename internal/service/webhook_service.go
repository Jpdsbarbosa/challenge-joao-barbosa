package service

import (
	"log"

	"github.com/jpdsbarbosa/challenge-joao-barbosa/internal/domain"
)

// WebhookServiceImpl implementa a lógica de processamento de webhooks
type WebhookServiceImpl struct {
	transferService *TransferService
}

// NewWebhookService cria uma nova instância do serviço
func NewWebhookService(transferService *TransferService) *WebhookServiceImpl {
	return &WebhookServiceImpl{
		transferService: transferService,
	}
}

// ProcessEvent processa um evento de webhook
func (s *WebhookServiceImpl) ProcessEvent(event domain.WebhookEvent) error {
	log.Printf("📋 Processando evento: Tipo=%s | Status=%s\n", event.EventType, event.Status)

	// Processar apenas invoices creditados
	if event.Subscription != "invoice" {
		log.Printf("⏭️  Evento ignorado: subscription=%s\n", event.Subscription)
		return nil
	}

	// IMPORTANTE: Processar APENAS 'credited', NÃO 'paid'
	// O desafio pede: "Receives the webhook callback of the Invoice credit"
	// Se processar 'paid' também, cria transferências duplicadas!
	if event.EventType != "credited" {
		log.Printf("⏭️  Invoice ainda não creditado: tipo=%s (aguardando 'credited')\n", event.EventType)
		return nil
	}

	log.Printf("💰 Invoice creditado detectado! ID: %s | Valor: R$%.2f | Taxa: R$%.2f\n",
		event.InvoiceID,
		float64(event.Amount)/100,
		float64(event.Fee)/100)

	// Criar transferência com o valor recebido menos as taxas
	_, err := s.transferService.CreateFromInvoicePayment(
		event.InvoiceID,
		event.Amount,
		event.Fee,
	)

	return err
}

// ValidateSignature valida a assinatura digital de um webhook
func (s *WebhookServiceImpl) ValidateSignature(body, signature string) bool {
	// TODO: Implementar validação de assinatura usando Event.Parse
	// Por enquanto, retorna true para desenvolvimento
	return true
}
