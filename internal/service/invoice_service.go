package service

import (
	"log"
	"math/rand"

	"github.com/jpdsbarbosa/challenge-joao-barbosa/internal/domain"
)

// InvoiceService gerencia a lógica de negócio relacionada a invoices
type InvoiceService struct {
	repo domain.InvoiceRepository
}

// NewInvoiceService cria uma nova instância do serviço
func NewInvoiceService(repo domain.InvoiceRepository) *InvoiceService {
	return &InvoiceService{
		repo: repo,
	}
}

// GenerateRandomInvoices gera entre 8 e 12 invoices aleatórios
func (s *InvoiceService) GenerateRandomInvoices() ([]domain.Invoice, error) {
	count := rand.Intn(5) + 8 // 8-12
	log.Printf("📝 Gerando %d invoices...\n", count)

	invoices := make([]domain.Invoice, count)
	for i := 0; i < count; i++ {
		invoices[i] = s.generateRandomInvoice()
	}

	created, err := s.repo.Create(invoices)
	if err != nil {
		log.Printf("❌ Erro ao criar invoices: %v\n", err)
		return nil, err
	}

	// Log dos invoices criados
	for _, invoice := range created {
		log.Printf("✅ Invoice criado: ID=%s | Valor=R$%.2f | Nome=%s\n",
			invoice.ID,
			float64(invoice.Amount)/100,
			invoice.Name)
	}

	log.Printf("🎉 Total de %d invoices criados com sucesso!\n", len(created))
	return created, nil
}

// generateRandomInvoice gera um invoice com dados aleatórios
func (s *InvoiceService) generateRandomInvoice() domain.Invoice {
	names := []string{
		"João Silva", "Maria Santos", "Pedro Oliveira", "Ana Costa",
		"Carlos Souza", "Juliana Lima", "Fernando Alves", "Patricia Rocha",
		"Roberto Martins", "Camila Ferreira", "Lucas Pereira", "Fernanda Gomes",
	}

	// CPFs válidos para testes (SEM formatação - apenas números)
	// Valor aleatório entre R$ 100 e R$ 1000
	amount := rand.Intn(90000) + 10000 // R$ 100 a R$ 1000
	name := names[rand.Intn(len(names))]

	// GERAR CPF DINAMICAMENTE - sempre válido!
	taxId := GenerateCPF()

	invoice := domain.Invoice{
		Amount: amount,
		Name:   name,
		TaxID:  taxId,
	}

	log.Printf("🎲 Invoice: %s | CPF:%s | R$%.2f\n",
		name, taxId, float64(amount)/100)

	return invoice
}

// GetByID busca um invoice por ID
func (s *InvoiceService) GetByID(id string) (*domain.Invoice, error) {
	return s.repo.GetByID(id)
}

// List lista invoices
func (s *InvoiceService) List(limit int) ([]domain.Invoice, error) {
	return s.repo.List(limit)
}
