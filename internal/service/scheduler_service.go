package service

import (
	"log"
	"time"
)

// SchedulerService gerencia tarefas agendadas
type SchedulerService struct {
	invoiceService *InvoiceService
	stopChan       chan bool
}

// NewSchedulerService cria uma nova instância do serviço
func NewSchedulerService(invoiceService *InvoiceService) *SchedulerService {
	return &SchedulerService{
		invoiceService: invoiceService,
		stopChan:       make(chan bool),
	}
}

// StartInvoiceGeneration inicia a geração periódica de invoices
func (s *SchedulerService) StartInvoiceGeneration() {
	log.Println("🚀 Iniciando gerador de invoices...")
	log.Println("📋 Configuração: 8-12 invoices a cada 3 horas durante 24 horas")

	// Gerar invoices imediatamente
	if _, err := s.invoiceService.GenerateRandomInvoices(); err != nil {
		log.Printf("❌ Erro ao gerar invoices iniciais: %v\n", err)
	}

	// Ticker para executar a cada 3 horas
	ticker := time.NewTicker(3 * time.Hour)
	defer ticker.Stop()

	// Timer para parar após 24 horas
	stopTimer := time.NewTimer(24 * time.Hour)
	defer stopTimer.Stop()

	for {
		select {
		case <-ticker.C:
			if _, err := s.invoiceService.GenerateRandomInvoices(); err != nil {
				log.Printf("❌ Erro ao gerar invoices: %v\n", err)
			}
		case <-stopTimer.C:
			log.Println("⏰ 24 horas completadas! Parando gerador de invoices...")
			return
		case <-s.stopChan:
			log.Println("🛑 Gerador de invoices interrompido manualmente")
			return
		}
	}
}

// Stop para o scheduler
func (s *SchedulerService) Stop() {
	close(s.stopChan)
}
