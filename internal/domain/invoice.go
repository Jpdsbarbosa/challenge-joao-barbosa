package domain

import "time"

// Invoice representa uma fatura no domínio da aplicação
type Invoice struct {
	ID         string
	Amount     int
	Name       string
	TaxID      string
	Due        *time.Time
	Expiration int
	Status     string
	Fee        int
	Created    *time.Time
}

// InvoiceRepository define a interface para operações com invoices
type InvoiceRepository interface {
	Create(invoices []Invoice) ([]Invoice, error)
	GetByID(id string) (*Invoice, error)
	List(limit int) ([]Invoice, error)
}
