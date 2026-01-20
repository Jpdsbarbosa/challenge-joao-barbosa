package domain

import "time"

// Transfer representa uma transferência no domínio da aplicação
type Transfer struct {
	ID            string
	Amount        int
	BankCode      string
	BranchCode    string
	AccountNumber string
	Name          string
	TaxID         string
	AccountType   string
	Description   string
	ExternalID    string // ID único para idempotência
	Status        string
	Fee           int
	Created       *time.Time
}

// TransferRepository define a interface para operações com transferências
type TransferRepository interface {
	Create(transfers []Transfer) ([]Transfer, error)
	GetByID(id string) (*Transfer, error)
	List(limit int) ([]Transfer, error)
}
