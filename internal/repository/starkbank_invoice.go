package repository

import (
	"fmt"

	"github.com/jpdsbarbosa/challenge-joao-barbosa/internal/domain"
	Invoice "github.com/starkbank/sdk-go/starkbank/invoice"
)

// StarkBankInvoiceRepository implementa InvoiceRepository usando o SDK da StarkBank
type StarkBankInvoiceRepository struct{}

// NewStarkBankInvoiceRepository cria uma nova instância do repositório
func NewStarkBankInvoiceRepository() *StarkBankInvoiceRepository {
	return &StarkBankInvoiceRepository{}
}

// Create cria invoices na StarkBank
func (r *StarkBankInvoiceRepository) Create(invoices []domain.Invoice) ([]domain.Invoice, error) {
	// Converter domain.Invoice para Invoice.Invoice
	sdkInvoices := make([]Invoice.Invoice, len(invoices))
	for i, inv := range invoices {
		sdkInvoices[i] = Invoice.Invoice{
			Amount:   inv.Amount,
			Name:     inv.Name,
			TaxId:    inv.TaxID,
			Fine:     2.5, // 2.5% multa após vencimento
			Interest: 1.3, // 1.3% juros mensal
		}
		fmt.Printf("✅ Invoice %d: R$%.2f | %s | CPF:%s\n",
			i+1, float64(inv.Amount)/100, inv.Name, inv.TaxID)
	}

	fmt.Printf("📤 Enviando %d invoices para StarkBank API...\n", len(sdkInvoices))

	// Criar na StarkBank
	created, err := Invoice.Create(sdkInvoices, nil)
	if err.Errors != nil {
		fmt.Printf("❌ Resposta da API: %+v\n", err)
		return nil, fmt.Errorf("erro ao criar invoices: %v", err.Errors)
	}

	// Converter de volta para domain.Invoice
	result := make([]domain.Invoice, len(created))
	for i, inv := range created {
		result[i] = domain.Invoice{
			ID:         inv.Id,
			Amount:     inv.Amount,
			Name:       inv.Name,
			TaxID:      inv.TaxId,
			Due:        inv.Due,
			Expiration: inv.Expiration,
			Status:     inv.Status,
			Fee:        inv.Fee,
			Created:    inv.Created,
		}
	}

	return result, nil
}

// GetByID busca um invoice por ID
func (r *StarkBankInvoiceRepository) GetByID(id string) (*domain.Invoice, error) {
	inv, err := Invoice.Get(id, nil)
	if err.Errors != nil {
		return nil, fmt.Errorf("erro ao buscar invoice: %v", err.Errors)
	}

	return &domain.Invoice{
		ID:         inv.Id,
		Amount:     inv.Amount,
		Name:       inv.Name,
		TaxID:      inv.TaxId,
		Due:        inv.Due,
		Expiration: inv.Expiration,
		Status:     inv.Status,
		Fee:        inv.Fee,
		Created:    inv.Created,
	}, nil
}

// List lista invoices
func (r *StarkBankInvoiceRepository) List(limit int) ([]domain.Invoice, error) {
	params := map[string]interface{}{
		"limit": limit,
	}

	invoices, errChan := Invoice.Query(params, nil)
	result := []domain.Invoice{}

	for {
		select {
		case inv, ok := <-invoices:
			if !ok {
				return result, nil
			}
			result = append(result, domain.Invoice{
				ID:         inv.Id,
				Amount:     inv.Amount,
				Name:       inv.Name,
				TaxID:      inv.TaxId,
				Due:        inv.Due,
				Expiration: inv.Expiration,
				Status:     inv.Status,
				Fee:        inv.Fee,
				Created:    inv.Created,
			})
		case err, ok := <-errChan:
			if ok && err.Errors != nil {
				return result, fmt.Errorf("erro ao listar invoices: %v", err.Errors)
			}
			return result, nil
		}
	}
}
