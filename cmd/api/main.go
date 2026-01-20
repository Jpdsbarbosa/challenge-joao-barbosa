package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/starkbank/sdk-go/starkbank"
	"github.com/starkinfra/core-go/starkcore/user/project"

	"github.com/jpdsbarbosa/challenge-joao-barbosa/internal/config"
	"github.com/jpdsbarbosa/challenge-joao-barbosa/internal/handler"
	"github.com/jpdsbarbosa/challenge-joao-barbosa/internal/middleware"
	"github.com/jpdsbarbosa/challenge-joao-barbosa/internal/repository"
	"github.com/jpdsbarbosa/challenge-joao-barbosa/internal/service"
)

func main() {
	// Inicializar seed para números aleatórios
	rand.Seed(time.Now().UnixNano())

	// Banner
	printBanner()

	// Carregar configurações
	log.Println("⚙️  Carregando configurações...")
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("❌ Erro ao carregar configurações: %v\n", err)
	}

	// Inicializar SDK da StarkBank
	log.Println("🔐 Inicializando SDK da StarkBank...")
	starkbank.User = project.Project{
		Environment: cfg.StarkBank.Environment,
		Id:          cfg.StarkBank.ProjectID,
		PrivateKey:  cfg.StarkBank.PrivateKey,
	}
	log.Println("✅ SDK inicializado com sucesso!")

	// Inicializar repositórios
	invoiceRepo := repository.NewStarkBankInvoiceRepository()
	transferRepo := repository.NewStarkBankTransferRepository()

	// Inicializar serviços
	invoiceService := service.NewInvoiceService(invoiceRepo)
	transferService := service.NewTransferService(transferRepo, cfg.Destination)
	webhookService := service.NewWebhookService(transferService)
	schedulerService := service.NewSchedulerService(invoiceService)

	// Inicializar handlers
	webhookHandler := handler.NewWebhookHandler(webhookService)
	healthHandler := handler.NewHealthHandler()
	balanceHandler := handler.NewBalanceHandler()

	// Configurar rotas
	mux := http.NewServeMux()
	mux.HandleFunc("/webhook", webhookHandler.Handle)
	mux.HandleFunc("/health", healthHandler.Handle)
	mux.HandleFunc("/balance", balanceHandler.Handle)

	// Aplicar middlewares
	handlerWithMiddleware := middleware.Recovery(middleware.Logger(mux))

	// Iniciar scheduler em background
	go schedulerService.StartInvoiceGeneration()

	// Configurar servidor HTTP
	server := &http.Server{
		Addr:         cfg.Server.Host + ":" + cfg.Server.Port,
		Handler:      handlerWithMiddleware,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Canal para capturar sinais de interrupção
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Iniciar servidor em goroutine
	go func() {
		log.Printf("🌐 Servidor HTTP iniciado em %s\n", server.Addr)
		log.Printf("📡 Endpoint webhook: http://localhost:%s/webhook\n", cfg.Server.Port)
		log.Printf("❤️  Endpoint health: http://localhost:%s/health\n", cfg.Server.Port)
		log.Printf("💰 Endpoint balance: http://localhost:%s/balance\n", cfg.Server.Port)
		log.Println("💡 Dica: Use ngrok para expor localmente: ngrok http", cfg.Server.Port)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Erro ao iniciar servidor: %v\n", err)
		}
	}()

	// Aguardar sinal de interrupção
	<-sigChan
	log.Println("\n🛑 Recebido sinal de interrupção. Encerrando aplicação...")
	schedulerService.Stop()
	log.Println("👋 Aplicação encerrada!")
}

func printBanner() {
	banner := `
╔═══════════════════════════════════════════════════════════╗
║                                                           ║
║      🏦  STARK BANK CHALLENGE - WEBHOOK PROCESSOR  🏦      ║
║                                                           ║
║  Arquitetura: Clean Architecture + Repository Pattern    ║
║  Desenvolvido por: João Pedro                            ║
║  Challenge: Backend Developer                            ║
║                                                           ║
╚═══════════════════════════════════════════════════════════╝

🎯 Objetivos:
  ✓ Gerar 8-12 invoices a cada 3 horas (24h)
  ✓ Receber webhooks de pagamento
  ✓ Validar assinaturas digitais
  ✓ Criar transferências automáticas

📐 Arquitetura:
  ✓ Domain-Driven Design
  ✓ Repository Pattern
  ✓ Service Layer
  ✓ Dependency Injection
  ✓ Middleware Chain
  ✓ Clean Code Principles

`
	log.Println(banner)
}
