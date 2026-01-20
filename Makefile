.PHONY: help run build test clean lint fmt install-deps

# Variáveis
APP_NAME=challenge-joao-barbosa
CMD_PATH=./cmd/api
BUILD_DIR=./bin

help: ## Mostra esta mensagem de ajuda
	@echo "Comandos disponíveis:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

install-deps: ## Instala as dependências do projeto
	@echo "📦 Instalando dependências..."
	go mod tidy
	go mod download

run: ## Executa a aplicação
	@echo "🚀 Iniciando aplicação..."
	@if [ -z "$$STARK_PROJECT_ID" ]; then \
		echo ""; \
		echo "❌ ERRO: Variável STARK_PROJECT_ID não definida"; \
		echo ""; \
		echo "Configure com:"; \
		echo "  export STARK_PROJECT_ID=\"seu-project-id\""; \
		echo ""; \
		echo "Ou copie env.example para .env e preencha os valores"; \
		echo ""; \
		exit 1; \
	fi
	go run $(CMD_PATH)/main.go

build: ## Compila a aplicação
	@echo "🔨 Compilando aplicação..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) $(CMD_PATH)/main.go
	@echo "✅ Binário gerado em: $(BUILD_DIR)/$(APP_NAME)"

test: ## Executa os testes
	@echo "🧪 Executando testes..."
	go test -v -cover ./...

test-coverage: ## Executa os testes com coverage
	@echo "🧪 Executando testes com coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report gerado: coverage.html"

lint: ## Executa o linter
	@echo "🔍 Executando linter..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run ./...; \
	else \
		echo "⚠️  golangci-lint não instalado. Use: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

fmt: ## Formata o código
	@echo "✨ Formatando código..."
	go fmt ./...
	gofmt -s -w .

clean: ## Remove binários e arquivos temporários
	@echo "🧹 Limpando arquivos..."
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html
	@echo "✅ Limpeza concluída"

dev: install-deps fmt ## Prepara o ambiente de desenvolvimento
	@echo "✅ Ambiente de desenvolvimento pronto!"

check: fmt lint test ## Executa todas as verificações
	@echo "✅ Todas as verificações passaram!"

webhook-setup: ## Configura webhook (uso: make webhook-setup URL=https://sua-url.ngrok-free.app)
	@if [ -z "$(URL)" ]; then \
		echo "❌ Erro: URL não fornecida"; \
		echo ""; \
		echo "Uso: make webhook-setup URL=https://sua-url.ngrok-free.app"; \
		echo ""; \
		echo "1️⃣  Execute em outro terminal: ngrok http 8080"; \
		echo "2️⃣  Copie a URL gerada (ex: https://abc123.ngrok-free.app)"; \
		echo "3️⃣  Execute: make webhook-setup URL=<url-copiada>"; \
		exit 1; \
	fi
	@echo "🔗 Configurando webhook..."
	go run ./scripts/setup_webhook.go $(URL)

ngrok: ## Inicia ngrok (expõe localhost:8080 para internet)
	@echo "🌐 Iniciando ngrok..."
	@echo "💡 Copie a URL que aparecer e use com: make webhook-setup URL=<url>"
	@echo ""
	ngrok http 8080

ngrok-url: ## Mostra a URL do ngrok que está rodando
	@echo "🔍 Buscando URL do ngrok..."
	@curl -s http://localhost:4040/api/tunnels | grep -o '"public_url":"[^"]*"' | grep https | cut -d'"' -f4 || echo "❌ Ngrok não está rodando"

test-webhook: ## Testa o webhook localmente
	@echo "🧪 Testando webhook..."
	@./scripts/test_webhook.sh

balance: ## Consulta o saldo da conta
	@echo "💰 Consultando saldo..."
	@curl -s http://localhost:8080/balance | python3 -m json.tool || curl -s http://localhost:8080/balance

health: ## Verifica status do servidor
	@echo "❤️  Verificando status..."
	@curl -s http://localhost:8080/health | python3 -m json.tool || curl -s http://localhost:8080/health
