package main

import (
	"fmt"
	"os"

	"github.com/starkbank/sdk-go/starkbank"
	Webhook "github.com/starkbank/sdk-go/starkbank/webhook"
	"github.com/starkinfra/core-go/starkcore/user/project"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("❌ Uso: go run scripts/setup_webhook.go <URL_DO_NGROK>")
		fmt.Println()
		fmt.Println("Exemplo:")
		fmt.Println("  go run scripts/setup_webhook.go https://abc123.ngrok-free.app")
		fmt.Println()
		fmt.Println("💡 Dica:")
		fmt.Println("  1. Execute: ngrok http 8080")
		fmt.Println("  2. Copie a URL gerada (https://...)")
		fmt.Println("  3. Execute este script com a URL")
		return
	}

	ngrokURL := os.Args[1]
	webhookURL := ngrokURL + "/webhook"

	// Carregar chave privada
	content, err := os.ReadFile("privateKeyChallenge.pem")
	if err != nil {
		fmt.Printf("❌ Erro ao ler chave: %v\n", err)
		return
	}

	// Configurar StarkBank
	starkbank.User = project.Project{
		Environment: "sandbox",
		Id:          "6211225704726528",
		PrivateKey:  string(content),
	}

	fmt.Println("╔════════════════════════════════════════════════════════╗")
	fmt.Println("║                                                        ║")
	fmt.Println("║        🔗 CONFIGURAÇÃO DE WEBHOOK - STARKBANK 🔗        ║")
	fmt.Println("║                                                        ║")
	fmt.Println("╚════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Printf("📡 URL do Webhook: %s\n", webhookURL)
	fmt.Println()

	// Primeiro, listar webhooks existentes
	fmt.Println("🔍 Verificando webhooks existentes...")
	webhooks, errChan := Webhook.Query(nil, nil)

	existingWebhooks := []Webhook.Webhook{}
	for {
		select {
		case webhook, ok := <-webhooks:
			if !ok {
				webhooks = nil
			} else {
				existingWebhooks = append(existingWebhooks, webhook)
			}
		case err, ok := <-errChan:
			if ok && err.Errors != nil {
				fmt.Printf("⚠️  Erro ao listar webhooks: %v\n", err.Errors)
			}
			errChan = nil
		}
		if webhooks == nil && errChan == nil {
			break
		}
	}

	fmt.Printf("   Encontrados: %d webhooks\n", len(existingWebhooks))

	// Deletar webhooks antigos para invoice
	for _, webhook := range existingWebhooks {
		hasInvoiceSubscription := false
		for _, sub := range webhook.Subscriptions {
			if sub == "invoice" {
				hasInvoiceSubscription = true
				break
			}
		}

		if hasInvoiceSubscription {
			fmt.Printf("   🗑️  Deletando webhook antigo: %s\n", webhook.Id)
			_, delErr := Webhook.Delete(webhook.Id, nil)
			if delErr.Errors != nil {
				fmt.Printf("      ⚠️  Erro ao deletar: %v\n", delErr.Errors)
			} else {
				fmt.Println("      ✅ Deletado")
			}
		}
	}

	fmt.Println()
	fmt.Println("📤 Criando novo webhook...")

	// Criar novo webhook
	created, errResp := Webhook.Create(
		Webhook.Webhook{
			Url:           webhookURL,
			Subscriptions: []string{"invoice"},
		}, nil)

	if errResp.Errors != nil {
		fmt.Println()
		fmt.Println("❌ ERRO AO CRIAR WEBHOOK:")
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		for i, e := range errResp.Errors {
			fmt.Printf("   %d. [%s] %s\n", i+1, e.Code, e.Message)
		}
		fmt.Println()
		fmt.Println("💡 Verifique:")
		fmt.Println("   - URL está acessível (teste: curl " + webhookURL + ")")
		fmt.Println("   - ngrok está rodando")
		fmt.Println("   - Servidor da aplicação está rodando (make run)")
		return
	}

	fmt.Println()
	fmt.Println("✅ WEBHOOK CONFIGURADO COM SUCESSO!")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("   ID: %s\n", created.Id)
	fmt.Printf("   URL: %s\n", created.Url)
	fmt.Printf("   Subscriptions: %v\n", created.Subscriptions)
	fmt.Println()
	fmt.Println("🎯 PRÓXIMOS PASSOS:")
	fmt.Println("   1. ✅ Webhook configurado")
	fmt.Println("   2. ✅ Invoices sendo criados")
	fmt.Println("   3. ⏳ Aguarde os invoices serem pagos (automático no Sandbox)")
	fmt.Println("   4. 📨 Você receberá webhooks quando os pagamentos acontecerem")
	fmt.Println("   5. 💸 Transfers serão criadas automaticamente!")
	fmt.Println()
	fmt.Println("📊 Monitore os logs do servidor para ver os webhooks chegando:")
	fmt.Println("   make run")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}
