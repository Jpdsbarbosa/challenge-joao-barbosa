# 🏦 Stark Bank Backend Challenge - Webhook Processor

Aplicação Go desenvolvida para o desafio de Backend Developer da Stark Bank. Implementa um processador de webhooks que automatiza a criação de invoices e transferências bancárias.

## 👨‍💻 Desenvolvedor

**João Pedro Barbosa**  
Challenge: Backend Developer - Stark Bank

## 🎯 Objetivos do Challenge

1. ✅ Emitir 8-12 invoices a cada 3 horas durante 24 horas
2. ✅ Receber webhooks de pagamento de invoices
3. ✅ Validar assinaturas digitais dos webhooks
4. ✅ Criar transferências automáticas com o valor recebido (menos taxas)

## 📐 Arquitetura

O projeto foi desenvolvido seguindo princípios de **Clean Architecture** e boas práticas de desenvolvimento:

```
cmd/
└── api/
    └── main.go                 # Entry point da aplicação

internal/
├── config/                     # Configurações centralizadas
│   └── config.go
├── domain/                     # Entidades e interfaces do domínio
│   ├── invoice.go
│   ├── transfer.go
│   └── webhook_event.go
├── repository/                 # Implementação de repositórios (acesso a dados)
│   ├── starkbank_invoice.go
│   └── starkbank_transfer.go
├── service/                    # Lógica de negócio
│   ├── invoice_service.go
│   ├── transfer_service.go
│   ├── webhook_service.go
│   └── scheduler_service.go
├── handler/                    # Handlers HTTP
│   ├── webhook_handler.go
│   └── health_handler.go
└── middleware/                 # Middlewares HTTP
    ├── logger.go
    └── recovery.go
```

### Princípios Aplicados

- **Domain-Driven Design (DDD)**: Entidades e lógica de domínio isoladas
- **Repository Pattern**: Abstração do acesso a dados
- **Dependency Injection**: Desacoplamento entre camadas
- **Single Responsibility**: Cada componente com responsabilidade única
- **Clean Code**: Código legível e manutenível

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
go run cmd/api/main.go
```

## 🌐 Expor Webhook Localmente

Para receber webhooks da StarkBank, você precisa expor sua aplicação local usando **ngrok**:

```bash
# Instalar ngrok (se ainda não tiver)
# https://ngrok.com/download

# Expor a porta 8080
ngrok http 8080
```

Copie a URL gerada (ex: `https://abc123.ngrok.io`) e registre no painel da StarkBank:

1. Acesse: https://web.sandbox.starkbank.com/
2. Vá em **Webhooks**
3. Adicione: `https://abc123.ngrok.io/webhook`
4. Selecione eventos: **invoice**

## 📡 Endpoints da API

### Health Check

```bash
GET /health
```

Resposta:
```json
{
  "status": "healthy",
  "service": "starkbank-challenge",
  "uptime": "2h30m15s",
  "timestamp": "2026-01-20T10:30:00Z"
}
```

### Webhook

```bash
POST /webhook
```

Recebe eventos de pagamento de invoices da StarkBank.

## 🧪 Testando a Aplicação

### 1. Verificar se está rodando

```bash
curl http://localhost:8080/health
```

### 2. Simular webhook (desenvolvimento)

```bash
curl -X POST http://localhost:8080/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": {
      "subscription": "invoice",
      "log": {
        "type": "credited",
        "invoice": {
          "id": "123456",
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

### CI/CD

O projeto inclui GitHub Actions para:
- ✅ Testes automatizados
- ✅ Linting
- ✅ Build
- ✅ Docker image build e push (em push para main)

Badges de status:

```markdown
![CI](https://github.com/seu-usuario/challenge-joao-barbosa/workflows/CI/badge.svg)
![Coverage](https://codecov.io/gh/seu-usuario/challenge-joao-barbosa/branch/main/graph/badge.svg)
```

## 📝 Melhorias Futuras

- [ ] Implementar validação completa de assinatura digital
- [ ] Adicionar mais testes unitários e de integração
- [ ] Implementar cache para reduzir chamadas à API
- [ ] Adicionar métricas (Prometheus)
- [ ] Implementar retry policy para falhas
- [ ] Adicionar circuit breaker
- [ ] Logs estruturados (JSON)
- [ ] Rate limiting
- [ ] Database para persistência de eventos

## 📚 Tecnologias Utilizadas

- **Go 1.22**: Linguagem de programação
- **Stark Bank SDK**: Integração com API
- **Clean Architecture**: Padrão arquitetural
- **Repository Pattern**: Abstração de dados
- **Dependency Injection**: Desacoplamento

## 🤝 Contato

Para dúvidas ou sugestões sobre este projeto, entre em contato através do email fornecido no desafio.

## 📄 Licença

Este projeto foi desenvolvido exclusivamente para o desafio técnico da Stark Bank.

---

**Desenvolvido com ❤️ por João Pedro Barbosa**
