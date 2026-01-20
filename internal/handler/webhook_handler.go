package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/jpdsbarbosa/challenge-joao-barbosa/internal/domain"
)

// WebhookHandler gerencia requisições de webhook
type WebhookHandler struct {
	webhookService domain.WebhookService
}

// NewWebhookHandler cria uma nova instância do handler
func NewWebhookHandler(webhookService domain.WebhookService) *WebhookHandler {
	return &WebhookHandler{
		webhookService: webhookService,
	}
}

// Handle processa requisições de webhook
func (h *WebhookHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	log.Println("📨 Webhook recebido!")

	// Ler o corpo da requisição
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("❌ Erro ao ler corpo da requisição: %v\n", err)
		http.Error(w, "Erro ao ler requisição", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Obter a assinatura digital do header
	signature := r.Header.Get("Digital-Signature")

	// Validar assinatura (em produção, deve rejeitar sem assinatura válida)
	if signature == "" {
		log.Println("⚠️  Webhook sem assinatura digital")
	}

	// Parse do evento
	event, err := h.parseEvent(body)
	if err != nil {
		log.Printf("❌ Erro ao fazer parse do evento: %v\n", err)
		http.Error(w, "Evento inválido", http.StatusBadRequest)
		return
	}

	// Processar o evento
	if err := h.webhookService.ProcessEvent(*event); err != nil {
		log.Printf("❌ Erro ao processar evento: %v\n", err)
		http.Error(w, "Erro ao processar evento", http.StatusInternalServerError)
		return
	}

	// Responder com 200 OK
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Event processed",
	})
}

// parseEvent faz o parse do JSON do webhook para domain.WebhookEvent
func (h *WebhookHandler) parseEvent(body []byte) (*domain.WebhookEvent, error) {
	// LOG DETALHADO PARA DEBUG
	log.Printf("🔍 Body recebido: %s\n", string(body))

	var eventData map[string]interface{}
	if err := json.Unmarshal(body, &eventData); err != nil {
		return nil, err
	}

	// Tentar formato com wrapper "event" (formato real da StarkBank)
	var eventLog map[string]interface{}
	if wrapper, ok := eventData["event"].(map[string]interface{}); ok {
		eventLog = wrapper
		log.Println("✅ Formato com wrapper 'event' detectado")
	} else {
		// Formato direto (nosso teste)
		eventLog = eventData
		log.Println("✅ Formato direto detectado")
	}

	subscription, _ := eventLog["subscription"].(string)
	log.Printf("📋 Subscription: %s\n", subscription)

	// FILTRAR: Processar apenas webhooks de invoice
	if subscription != "invoice" {
		log.Printf("⏭️  Webhook ignorado: subscription=%s (esperado: invoice)\n", subscription)
		return &domain.WebhookEvent{
			Subscription: subscription,
		}, nil
	}

	// Extrair o log do evento
	logData, ok := eventLog["log"].(map[string]interface{})
	if !ok {
		log.Printf("❌ Campo 'log' não encontrado ou inválido. Dados: %+v\n", eventLog)
		return nil, fmt.Errorf("campo 'log' não encontrado no webhook")
	}

	eventType, _ := logData["type"].(string)
	log.Printf("📋 Event Type: %s\n", eventType)

	// Extrair dados do invoice
	invoiceData, ok := logData["invoice"].(map[string]interface{})
	if !ok {
		log.Printf("❌ Campo 'invoice' não encontrado. LogData: %+v\n", logData)
		return nil, fmt.Errorf("campo 'invoice' não encontrado no log do webhook")
	}

	invoiceID, _ := invoiceData["id"].(string)
	status, _ := invoiceData["status"].(string)
	amount, _ := invoiceData["amount"].(float64)
	fee, _ := invoiceData["fee"].(float64)

	log.Printf("💰 Invoice: ID=%s, Amount=%.2f, Fee=%.2f, Status=%s\n",
		invoiceID, amount/100, fee/100, status)

	return &domain.WebhookEvent{
		Subscription: subscription,
		EventType:    eventType,
		InvoiceID:    invoiceID,
		Amount:       int64(amount),
		Fee:          int64(fee),
		Status:       status,
	}, nil
}
