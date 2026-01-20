package repository

import (
	"fmt"

	"github.com/jpdsbarbosa/challenge-joao-barbosa/internal/domain"
	Transfer "github.com/starkbank/sdk-go/starkbank/transfer"
)

// StarkBankTransferRepository implementa TransferRepository usando o SDK da StarkBank
type StarkBankTransferRepository struct{}

// NewStarkBankTransferRepository cria uma nova instância do repositório
func NewStarkBankTransferRepository() *StarkBankTransferRepository {
	return &StarkBankTransferRepository{}
}

// Create cria transferências na StarkBank
func (r *StarkBankTransferRepository) Create(transfers []domain.Transfer) ([]domain.Transfer, error) {
	// Converter domain.Transfer para Transfer.Transfer
	sdkTransfers := make([]Transfer.Transfer, len(transfers))
	for i, t := range transfers {
		sdkTransfers[i] = Transfer.Transfer{
			Amount:        t.Amount,
			BankCode:      t.BankCode,
			BranchCode:    t.BranchCode,
			AccountNumber: t.AccountNumber,
			Name:          t.Name,
			TaxId:         t.TaxID,
			AccountType:   t.AccountType,
			Description:   t.Description,
			ExternalId:    t.ExternalID, // ID único para idempotência (gerado no service)
		}
	}

	// Criar na StarkBank
	created, err := Transfer.Create(sdkTransfers, nil)
	if err.Errors != nil {
		return nil, fmt.Errorf("erro ao criar transferências: %v", err.Errors)
	}

	// Converter de volta para domain.Transfer
	result := make([]domain.Transfer, len(created))
	for i, t := range created {
		result[i] = domain.Transfer{
			ID:            t.Id,
			Amount:        t.Amount,
			BankCode:      t.BankCode,
			BranchCode:    t.BranchCode,
			AccountNumber: t.AccountNumber,
			Name:          t.Name,
			TaxID:         t.TaxId,
			AccountType:   t.AccountType,
			Description:   t.Description,
			Status:        t.Status,
			Fee:           t.Fee,
			Created:       t.Created,
		}
	}

	return result, nil
}

// GetByID busca uma transferência por ID
func (r *StarkBankTransferRepository) GetByID(id string) (*domain.Transfer, error) {
	t, err := Transfer.Get(id, nil)
	if err.Errors != nil {
		return nil, fmt.Errorf("erro ao buscar transferência: %v", err.Errors)
	}

	return &domain.Transfer{
		ID:            t.Id,
		Amount:        t.Amount,
		BankCode:      t.BankCode,
		BranchCode:    t.BranchCode,
		AccountNumber: t.AccountNumber,
		Name:          t.Name,
		TaxID:         t.TaxId,
		AccountType:   t.AccountType,
		Description:   t.Description,
		Status:        t.Status,
		Fee:           t.Fee,
		Created:       t.Created,
	}, nil
}

// List lista transferências
func (r *StarkBankTransferRepository) List(limit int) ([]domain.Transfer, error) {
	params := map[string]interface{}{
		"limit": limit,
	}

	transfers, errChan := Transfer.Query(params, nil)
	result := []domain.Transfer{}

	for {
		select {
		case t, ok := <-transfers:
			if !ok {
				return result, nil
			}
			result = append(result, domain.Transfer{
				ID:            t.Id,
				Amount:        t.Amount,
				BankCode:      t.BankCode,
				BranchCode:    t.BranchCode,
				AccountNumber: t.AccountNumber,
				Name:          t.Name,
				TaxID:         t.TaxId,
				AccountType:   t.AccountType,
				Description:   t.Description,
				Status:        t.Status,
				Fee:           t.Fee,
				Created:       t.Created,
			})
		case err, ok := <-errChan:
			if ok && err.Errors != nil {
				return result, fmt.Errorf("erro ao listar transferências: %v", err.Errors)
			}
			return result, nil
		}
	}
}
