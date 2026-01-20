package service

import (
	"fmt"
	"log"
	"time"

	"github.com/jpdsbarbosa/challenge-joao-barbosa/internal/config"
	"github.com/jpdsbarbosa/challenge-joao-barbosa/internal/domain"
)

// TransferService gerencia a lógica de negócio relacionada a transferências
type TransferService struct {
	repo        domain.TransferRepository
	destination config.DestinationAccount
}

// NewTransferService cria uma nova instância do serviço
func NewTransferService(repo domain.TransferRepository, dest config.DestinationAccount) *TransferService {
	return &TransferService{
		repo:        repo,
		destination: dest,
	}
}

// CreateFromInvoicePayment cria uma transferência a partir de um pagamento de invoice
func (s *TransferService) CreateFromInvoicePayment(invoiceID string, amount, fee int64) (*domain.Transfer, error) {
	// Calcular valor líquido (valor recebido - taxas)
	netAmount := amount - fee

	if netAmount <= 0 {
		return nil, fmt.Errorf("valor líquido inválido: R$%.2f", float64(netAmount)/100)
	}

	log.Printf("💸 Criando transferência de R$%.2f (bruto: R$%.2f - taxa: R$%.2f)\n",
		float64(netAmount)/100,
		float64(amount)/100,
		float64(fee)/100)

	// Gerar ExternalID único e curto para idempotência
	externalID := fmt.Sprintf("inv-%s-%d", invoiceID, time.Now().Unix())

	transfer := domain.Transfer{
		Amount:        int(netAmount),
		BankCode:      s.destination.BankCode,
		BranchCode:    s.destination.BranchCode,
		AccountNumber: s.destination.AccountNumber,
		Name:          s.destination.Name,
		TaxID:         s.destination.TaxID,
		AccountType:   s.destination.AccountType,
		Description:   fmt.Sprintf("Transferência referente ao invoice %s", invoiceID),
		ExternalID:    externalID,
	}

	created, err := s.repo.Create([]domain.Transfer{transfer})
	if err != nil {
		log.Printf("❌ Erro ao criar transferência: %v\n", err)
		return nil, err
	}

	if len(created) == 0 {
		return nil, fmt.Errorf("nenhuma transferência criada")
	}

	result := &created[0]
	log.Printf("✅ Transferência criada com sucesso!\n")
	log.Printf("   ID: %s\n", result.ID)
	log.Printf("   Valor: R$%.2f\n", float64(result.Amount)/100)
	log.Printf("   Status: %s\n", result.Status)
	log.Printf("   Destinatário: %s\n", result.Name)
	log.Printf("   Invoice Origem: %s\n", invoiceID)

	return result, nil
}

// GetByID busca uma transferência por ID
func (s *TransferService) GetByID(id string) (*domain.Transfer, error) {
	return s.repo.GetByID(id)
}

// List lista transferências
func (s *TransferService) List(limit int) ([]domain.Transfer, error) {
	return s.repo.List(limit)
}
