package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	Balance "github.com/starkbank/sdk-go/starkbank/balance"
)

// BalanceHandler gerencia requisições de consulta de saldo
type BalanceHandler struct{}

// NewBalanceHandler cria uma nova instância do handler
func NewBalanceHandler() *BalanceHandler {
	return &BalanceHandler{}
}

// Handle processa requisições de consulta de saldo
func (h *BalanceHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	log.Println("💰 Consultando saldo da conta...")

	// Buscar saldo na StarkBank
	balance, err := Balance.Get(nil)
	if err.Errors != nil {
		log.Printf("❌ Erro ao consultar saldo: %v\n", err.Errors)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":   "Erro ao consultar saldo",
			"details": err.Errors,
		})
		return
	}

	log.Printf("✅ Saldo consultado com sucesso: R$%.2f\n", float64(balance.Amount)/100)

	// Preparar resposta
	response := map[string]interface{}{
		"amount":   balance.Amount,
		"currency": balance.Currency,
		"updated":  balance.Updated,
		"formatted": map[string]string{
			"amount":   fmt.Sprintf("R$ %.2f", float64(balance.Amount)/100),
			"currency": balance.Currency,
			"updated":  balance.Updated.Format("02/01/2006 15:04:05"),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
