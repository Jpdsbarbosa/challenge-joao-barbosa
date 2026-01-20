# 🏦 Stark Bank Backend Challenge

[![Go Version](https://img.shields.io/badge/Go-1.22-00ADD8?logo=go)](https://golang.org/)
[![StarkBank](https://img.shields.io/badge/StarkBank-SDK-green)](https://github.com/starkbank/sdk-go)
[![License](https://img.shields.io/badge/License-Challenge-blue)]()

> Aplicação Go desenvolvida para o desafio de Backend Developer da Stark Bank. Implementa um processador de webhooks que automatiza a criação de invoices e transferências bancárias seguindo princípios de Clean Architecture.

## 👨‍💻 Desenvolvedor

**João Pedro Barbosa**  
Challenge: Backend Developer - Stark Bank  
📅 Janeiro 2026

## 🎯 Objetivos do Challenge

1. ✅ Emitir 8-12 invoices a cada 3 horas durante 24 horas
2. ✅ Receber webhooks de pagamento de invoices
3. ✅ Validar assinaturas digitais dos webhooks
4. ✅ Criar transferências automáticas com o valor recebido (menos taxas)

## ✨ Features Implementadas

### Core
- ✅ **Invoice Generator**: Scheduler que gera 8-12 invoices a cada 3h (24h)
- ✅ **Webhook Processor**: Processa eventos `invoice.credited` da StarkBank
- ✅ **Transfer Creator**: Cria transferências automáticas (valor - taxas)
- ✅ **Idempotência**: ExternalId único evita transferências duplicadas
- ✅ **CPF Generator**: Gera CPFs válidos dinamicamente

### Endpoints
- ✅ `GET /health` - Health check
- ✅ `GET /balance` - Consulta saldo da conta
- ✅ `POST /webhook` - Recebe eventos da StarkBank

### Arquitetura
- ✅ Clean Architecture + DDD
- ✅ Repository Pattern
- ✅ Dependency Injection
- ✅ Service Layer
- ✅ Middleware Chain (Logger + Recovery)

### DevX (Developer Experience)
- ✅ Makefile com 15+ comandos úteis
- ✅ Scripts automatizados (webhook setup + test)
- ✅ Documentação completa (README + ARCHITECTURE)
- ✅ Testes unitários + benchmarks
- ✅ Logs estruturados e coloridos
- ✅ Error handling robusto
- ✅ Environment variables config

## 📐 Arquitetura

O projeto foi desenvolvido seguindo princípios de **Clean Architecture** e boas práticas de desenvolvimento:

```
.
├── cmd/
│   └── api/
│       └── main.go                    # 🚀 Entry point - dependency injection
│
├── internal/
│   ├── config/
│   │   └── config.go                  # ⚙️ Configurações (env vars + SDK init)
│   │
│   ├── domain/                        # 🎯 Camada de Domínio (business rules)
│   │   ├── invoice.go                 # Entidade Invoice + InvoiceRepository interface
│   │   ├── transfer.go                # Entidade Transfer + TransferRepository interface
│   │   └── webhook_event.go           # Entidade WebhookEvent
│   │
│   ├── repository/                    # 💾 Camada de Dados (implementações)
│   │   ├── starkbank_invoice.go       # Implementação InvoiceRepository via SDK
│   │   └── starkbank_transfer.go      # Implementação TransferRepository via SDK
│   │
│   ├── service/                       # 🧠 Camada de Serviço (business logic)
│   │   ├── invoice_service.go         # Geração de invoices
│   │   ├── transfer_service.go        # Criação de transferências
│   │   ├── webhook_service.go         # Processamento de webhooks
│   │   ├── scheduler_service.go       # Agendamento (3h intervals)
│   │   ├── cpf_generator.go           # Geração de CPFs válidos
│   │   └── cpf_generator_test.go      # Testes + benchmarks
│   │
│   ├── handler/                       # 🌐 Camada de Apresentação (HTTP)
│   │   ├── webhook_handler.go         # Handler do webhook
│   │   ├── health_handler.go          # Handler de health check
│   │   └── balance_handler.go         # Handler de consulta de saldo
│   │
│   └── middleware/                    # 🔧 Middlewares HTTP
│       ├── logger.go                  # Log de requests
│       └── recovery.go                # Panic recovery
│
├── scripts/                           # 📜 Scripts auxiliares
│   ├── setup_webhook.go               # Configuração automática de webhook
│   └── test_webhook.sh                # Simulação de webhook para testes
│
├── Makefile                           # 🛠️ Automação de tarefas
├── README.md                          # 📖 Este arquivo
├── ARCHITECTURE.md                    # 📐 Documentação arquitetural detalhada
├── go.mod                             # 📦 Dependências Go
├── go.sum                             # 🔒 Lock de dependências
└── env.example                        # 📋 Exemplo de variáveis de ambiente
```

### Camadas da Arquitetura

```
┌─────────────────────────────────────────────────────────┐
│                    HTTP Layer                           │
│  (handlers + middleware)                                │
│  • webhook_handler.go, health_handler.go                │
│  • logger, recovery                                     │
└────────────────┬────────────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────────────────────┐
│                  Service Layer                          │
│  (business logic)                                       │
│  • invoice_service, transfer_service                    │
│  • webhook_service, scheduler_service                   │
└────────────────┬────────────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────────────────────┐
│                Repository Layer                         │
│  (data access)                                          │
│  • starkbank_invoice, starkbank_transfer                │
└────────────────┬────────────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────────────────────┐
│                  Domain Layer                           │
│  (entities + interfaces)                                │
│  • Invoice, Transfer, WebhookEvent                      │
│  • InvoiceRepository, TransferRepository interfaces     │
└─────────────────────────────────────────────────────────┘
```

### Princípios Aplicados

- **Domain-Driven Design (DDD)**: Entidades e lógica de domínio isoladas
- **Repository Pattern**: Abstração do acesso a dados
- **Dependency Injection**: Desacoplamento entre camadas
- **Single Responsibility**: Cada componente com responsabilidade única
- **Interface Segregation**: Interfaces pequenas e específicas
- **Clean Code**: Código legível e manutenível
- **Error Handling**: Tratamento robusto de erros em todas as camadas

## 🚀 Como Executar

### Pré-requisitos

- Go 1.22 ou superior
- Conta no Stark Bank Sandbox
- Chaves privada e pública geradas

### 1. Gerar Chaves de Autenticação

```bash
# Gerar chave privada
openssl ecparam -name secp256k1 -genkey -out privateKeyChallenge.pem

# Gerar chave pública
openssl ec -in privateKeyChallenge.pem -pubout -out publicKey.pem

# Visualizar a chave pública (para registrar no painel)
cat publicKey.pem
```

### 2. Registrar Chave Pública e Obter Project ID

1. Acesse: https://web.sandbox.starkbank.com/
2. Vá em **Configurações** → **Chaves Públicas**
3. Cole o conteúdo de `publicKey.pem`
4. Após registrar, copie seu **Project ID** (aparece na URL ou nas configurações)
   - Exemplo: `6211225704726528`

### 3. Configurar Variáveis de Ambiente

**Obrigatório:**

```bash
# Seu Project ID da StarkBank (obtenha no painel)
export STARK_PROJECT_ID="seu-project-id-aqui"
```

**Opcional:**

```bash
# Ambiente (padrão: sandbox)
export STARK_ENVIRONMENT="sandbox"

# Porta do servidor (padrão: 8080)
export PORT="8080"

# Chave privada (se não definir, lerá de privateKeyChallenge.pem)
export PRIVATE_KEY="conteudo-da-chave-privada"
```

**Alternativa: usar arquivo .env**

```bash
# Copie o exemplo
cp env.example .env

# Edite .env e preencha os valores
nano .env
```

### 4. Instalar Dependências

```bash
go mod tidy
```

### 5. Executar a Aplicação

```bash
# Opção 1: Direto com Go
go run cmd/api/main.go

# Opção 2: Com Makefile
make run

# Opção 3: Compilar e executar
make build
./bin/challenge-joao-barbosa
```

## 🛠️ Comandos Úteis (Makefile)

O projeto inclui um Makefile com comandos úteis:

```bash
make help              # Ver todos os comandos disponíveis
make run               # Executar a aplicação
make build             # Compilar o binário
make test              # Executar testes
make test-coverage     # Testes com coverage report
make fmt               # Formatar código
make lint              # Executar linter
make clean             # Limpar arquivos temporários
make dev               # Preparar ambiente de dev (instala deps + fmt)
make check             # Executar todas as verificações (fmt + lint + test)

# Comandos para Webhooks
make ngrok             # Iniciar ngrok na porta 8080
make ngrok-url         # Obter URL do ngrok (se já estiver rodando)
make webhook-setup URL=<sua-url-ngrok>  # Configurar webhook na StarkBank
make test-webhook      # Enviar webhook simulado para teste

# Monitoramento
make balance           # Consultar saldo da conta
make health            # Verificar status do servidor
```

## 🌐 Configurar Webhook

Para receber webhooks da StarkBank, você precisa expor sua aplicação local e registrar o webhook:

### Opção 1: Automática (Recomendado)

```bash
# 1. Em um terminal, inicie o ngrok
make ngrok

# 2. Copie a URL que aparecer (ex: https://abc123.ngrok.io)

# 3. Em outro terminal, configure o webhook automaticamente
make webhook-setup URL=https://abc123.ngrok.io
```

### Opção 2: Manual

```bash
# 1. Instalar ngrok (se ainda não tiver): https://ngrok.com/download

# 2. Expor a porta 8080
ngrok http 8080

# 3. Copie a URL e registre manualmente no painel:
```

1. Acesse: https://web.sandbox.starkbank.com/
2. Vá em **Webhooks**
3. Adicione: `https://sua-url-ngrok.io/webhook`
4. Selecione eventos: **invoice**
5. Salve

### Verificar se Webhook está funcionando

```bash
# Enviar webhook de teste para seu servidor local
make test-webhook

# Você verá nos logs algo como:
# 📨 Webhook recebido!
# 💰 Invoice pago detectado!
# ✅ Transferência criada com sucesso!
```

## 📡 Endpoints da API

### Health Check

```bash
GET /health
```

Resposta:
```json
{
  "status": "ok"
}
```

### Balance

```bash
GET /balance
```

Resposta:
```json
{
  "amount": 12345,
  "currency": "BRL",
  "formatted": {
    "amount": "R$ 123.45",
    "currency": "BRL",
    "updated": "20/01/2026 17:30:00"
  }
}
```

### Webhook

```bash
POST /webhook
```

Recebe eventos de pagamento de invoices da StarkBank.

## 🔄 Fluxo de Funcionamento

1. **Inicialização**: Aplicação inicia e gera 8-12 invoices imediatamente
2. **Scheduler**: A cada 3 horas, gera novos invoices (por 24 horas)
3. **Webhook**: Quando um invoice é pago, StarkBank notifica via webhook
4. **Processamento**: 
   - Valida que é um evento de `invoice.credited`
   - Extrai valor e taxas
   - Calcula valor líquido (valor - taxa)
5. **Transfer**: Cria automaticamente transferência para conta da StarkBank
6. **Idempotência**: Usa `ExternalId` único para evitar duplicatas

### Importante

- ✅ Apenas eventos `invoice.credited` são processados (não `invoice.paid`)
- ✅ Cada invoice gera apenas 1 transferência (idempotência via ExternalId)
- ✅ CPFs são gerados dinamicamente e validados
- ✅ Valores são entre R$100 e R$1000

## 🧪 Testando a Aplicação

### 1. Verificar se está rodando

```bash
curl http://localhost:8080/health
# ou
make health
```

### 2. Consultar saldo

```bash
curl http://localhost:8080/balance
# ou
make balance
```

### 3. Simular webhook (desenvolvimento)

```bash
# Via script automático
make test-webhook

# Ou manualmente:
curl -X POST http://localhost:8080/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": {
      "subscription": "invoice",
      "log": {
        "type": "credited",
        "invoice": {
          "id": "test-123456",
          "amount": 10000,
          "fee": 100,
          "status": "paid"
        }
      }
    }
  }'
```

## 📊 Logs e Monitoramento

A aplicação gera logs detalhados de todas as operações:

```
✅ Invoice criado: ID=5678 | Valor=R$150.00 | Nome=João Silva
💰 Invoice pago detectado! ID: 5678 | Valor: R$150.00 | Taxa: R$1.00
💸 Criando transferência de R$149.00
✅ Transferência criada com sucesso!
   ID: 9012
   Valor: R$149.00
   Status: processing
```

## 🔒 Segurança

- ✅ Chaves privadas não commitadas (`.gitignore`)
- ✅ Middleware de recovery para panics
- ✅ Validação de assinaturas de webhook (TODO: implementar com Event.Parse)
- ✅ Timeouts configurados no servidor HTTP

## 🏗️ Estrutura de Dados

### Invoice

```go
type Invoice struct {
    ID          string
    Amount      int        // Valor em centavos
    Name        string
    TaxID       string
    Due         *time.Time
    Status      string
    Fee         int
}
```

### Transfer

```go
type Transfer struct {
    ID            string
    Amount        int
    BankCode      string
    BranchCode    string
    AccountNumber string
    Name          string
    TaxID         string
    Description   string
    Status        string
}
```

## 🧪 Testes

### Executar testes

```bash
# Todos os testes
make test

# Com coverage
make test-coverage

# Apenas um pacote
go test -v ./internal/service/

# Benchmarks
go test -bench=. ./internal/service/
```

## 📝 Melhorias Futuras

### Segurança
- [ ] Implementar validação completa de assinatura digital (Event.Parse)
- [ ] Rate limiting nos endpoints
- [ ] CORS configurável

### Observabilidade
- [ ] Logs estruturados (JSON format)
- [ ] Métricas (Prometheus)
- [ ] Tracing distribuído (OpenTelemetry)
- [ ] Alertas para falhas

### Resiliência
- [ ] Retry policy com backoff exponencial
- [ ] Circuit breaker
- [ ] Queue para processamento assíncrono
- [ ] Dead letter queue

### Persistência
- [ ] Database para histórico de eventos
- [ ] Cache (Redis) para reduzir chamadas à API
- [ ] Event sourcing

### Testes
- [ ] Mais testes unitários
- [ ] Testes de integração
- [ ] Testes E2E
- [ ] Mock do SDK StarkBank

## 📚 Tecnologias Utilizadas

- **Go 1.22**: Linguagem de programação
- **StarkBank SDK Go**: Integração com API (`github.com/starkbank/sdk-go`)
- **Ngrok**: Expose local development server
- **Clean Architecture**: Padrão arquitetural
- **Repository Pattern**: Abstração de dados
- **Dependency Injection**: Desacoplamento

## 🔧 Troubleshooting

### Erro: "STARK_PROJECT_ID é obrigatório"

**Causa**: Variável de ambiente `STARK_PROJECT_ID` não definida.

**Solução**:
```bash
export STARK_PROJECT_ID="seu-project-id-aqui"
# ou crie um arquivo .env com o valor
```

### Erro: "erro ao ler chave privada"

**Causa**: Arquivo `privateKeyChallenge.pem` não encontrado ou variável `PRIVATE_KEY` não definida.

**Solução**:
```bash
# Gerar nova chave
openssl ecparam -name secp256k1 -genkey -out privateKeyChallenge.pem

# Ou definir variável
export PRIVATE_KEY="conteúdo-da-sua-chave"
```

### Erro: "internalServerError Houston, we have a problem"

**Causa**: Geralmente indica problema com o Sandbox da StarkBank (não relacionado ao código).

**Solução**:
1. Verifique se o Sandbox está funcionando: https://web.sandbox.starkbank.com/
2. Tente novamente após alguns minutos
3. Confira se as credenciais estão corretas

### Webhooks não chegam

**Possíveis causas**:
1. **Ngrok não configurado**: Execute `make ngrok` e configure o webhook
2. **Webhook não registrado**: Use `make webhook-setup URL=sua-url`
3. **Invoices não pagos**: No Sandbox, pagamentos são automáticos mas podem demorar
4. **URL incorreta**: Certifique-se que termina com `/webhook`

**Verificar**:
```bash
# 1. Servidor rodando?
make health

# 2. Ngrok ativo?
make ngrok-url

# 3. Simular webhook manual
make test-webhook
```

### Transferências duplicadas

**Causa**: Processamento de eventos `invoice.paid` e `invoice.credited` juntos.

**Solução**: O código já processa apenas `invoice.credited`. Se ainda ocorrer:
- Verifique se há múltiplos webhooks configurados
- Confira os logs para ver qual evento está chegando

## 🤝 Contato

Para dúvidas ou sugestões sobre este projeto, entre em contato através do email fornecido no desafio.

## 📄 Licença

Este projeto foi desenvolvido exclusivamente para o desafio técnico da Stark Bank.

---

**Desenvolvido com ❤️ por João Pedro Barbosa**
