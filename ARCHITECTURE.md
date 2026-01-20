# 📐 Arquitetura do Projeto

## Visão Geral

Este projeto segue os princípios de **Clean Architecture** e **Domain-Driven Design (DDD)**, proporcionando uma estrutura escalável, testável e manutenível.

## Camadas da Aplicação

```
┌─────────────────────────────────────────────────────┐
│                     cmd/api                         │
│                   (Entry Point)                     │
└───────────────────┬─────────────────────────────────┘
                    │
┌───────────────────▼─────────────────────────────────┐
│                   Handler Layer                     │
│         (HTTP Handlers + Middlewares)               │
│  - webhook_handler.go                               │
│  - health_handler.go                                │
│  - logger.go (middleware)                           │
│  - recovery.go (middleware)                         │
└───────────────────┬─────────────────────────────────┘
                    │
┌───────────────────▼─────────────────────────────────┐
│                  Service Layer                      │
│              (Business Logic)                       │
│  - invoice_service.go                               │
│  - transfer_service.go                              │
│  - webhook_service.go                               │
│  - scheduler_service.go                             │
└───────────────────┬─────────────────────────────────┘
                    │
┌───────────────────▼─────────────────────────────────┐
│                 Repository Layer                    │
│            (Data Access Abstraction)                │
│  - starkbank_invoice.go                             │
│  - starkbank_transfer.go                            │
└───────────────────┬─────────────────────────────────┘
                    │
┌───────────────────▼─────────────────────────────────┐
│                  Domain Layer                       │
│         (Entities + Interfaces)                     │
│  - invoice.go                                       │
│  - transfer.go                                      │
│  - webhook_event.go                                 │
└─────────────────────────────────────────────────────┘
```

## Descrição das Camadas

### 1. Domain Layer (`internal/domain/`)

**Responsabilidade**: Define as entidades de negócio e interfaces (contratos).

**Características**:
- Não possui dependências externas
- Define as regras de negócio centrais
- Contém apenas interfaces e structs de dados

**Arquivos**:
- `invoice.go`: Entidade Invoice + InvoiceRepository interface
- `transfer.go`: Entidade Transfer + TransferRepository interface
- `webhook_event.go`: Entidade WebhookEvent + WebhookService interface

### 2. Repository Layer (`internal/repository/`)

**Responsabilidade**: Implementa as interfaces definidas no domínio, abstraindo o acesso a dados externos (StarkBank API).

**Características**:
- Implementa interfaces do domínio
- Encapsula toda a lógica de comunicação com APIs externas
- Converte entre tipos do SDK e tipos do domínio

**Arquivos**:
- `starkbank_invoice.go`: Implementação de InvoiceRepository
- `starkbank_transfer.go`: Implementação de TransferRepository

**Benefícios**:
- Facilita troca de providers
- Permite mock para testes
- Desacopla lógica de negócio da API externa

### 3. Service Layer (`internal/service/`)

**Responsabilidade**: Contém toda a lógica de negócio da aplicação.

**Características**:
- Orquestra operações entre repositórios
- Implementa regras de negócio complexas
- Gerencia transações e validações

**Arquivos**:
- `invoice_service.go`: Lógica para geração e gerenciamento de invoices
- `transfer_service.go`: Lógica para criação de transferências
- `webhook_service.go`: Processamento de eventos de webhook
- `scheduler_service.go`: Gerenciamento de tarefas agendadas

**Exemplo de Fluxo**:
```go
WebhookService.ProcessEvent()
    ↓
TransferService.CreateFromInvoicePayment()
    ↓
TransferRepository.Create()
    ↓
StarkBank API
```

### 4. Handler Layer (`internal/handler/` + `internal/middleware/`)

**Responsabilidade**: Gerencia requisições HTTP e aplica middlewares.

**Características**:
- Converte HTTP requests em chamadas de serviço
- Aplica validações de entrada
- Gerencia serialização/deserialização JSON

**Arquivos**:
- `webhook_handler.go`: Processa webhooks da StarkBank
- `health_handler.go`: Health check endpoint
- `logger.go`: Middleware de logging
- `recovery.go`: Middleware de recuperação de panics

### 5. Entry Point (`cmd/api/`)

**Responsabilidade**: Inicializa a aplicação e configura dependências.

**Características**:
- Dependency Injection manual
- Configuração do servidor HTTP
- Gerenciamento de lifecycle da aplicação

## Dependency Injection

O projeto usa **Dependency Injection manual** para manter as camadas desacopladas:

```go
// Criar repositórios (camada mais baixa)
invoiceRepo := repository.NewStarkBankInvoiceRepository()
transferRepo := repository.NewStarkBankTransferRepository()

// Criar serviços (injetar repositórios)
invoiceService := service.NewInvoiceService(invoiceRepo)
transferService := service.NewTransferService(transferRepo, cfg.Destination)
webhookService := service.NewWebhookService(transferService)

// Criar handlers (injetar serviços)
webhookHandler := handler.NewWebhookHandler(webhookService)
```

**Benefícios**:
- Testabilidade: Fácil mockar dependências
- Flexibilidade: Trocar implementações sem alterar código
- Clareza: Dependências explícitas

## Fluxo de Dados

### Geração de Invoices

```
Scheduler (3h timer)
    ↓
InvoiceService.GenerateRandomInvoices()
    ↓
InvoiceRepository.Create()
    ↓
StarkBank API
```

### Processamento de Webhook

```
POST /webhook
    ↓
WebhookHandler.Handle()
    ↓
WebhookService.ProcessEvent()
    ↓
TransferService.CreateFromInvoicePayment()
    ↓
TransferRepository.Create()
    ↓
StarkBank API
```

## Padrões de Design Utilizados

### 1. Repository Pattern
Abstrai o acesso a dados, permitindo trocar facilmente a fonte de dados.

### 2. Dependency Injection
Injeta dependências via construtores, facilitando testes e desacoplamento.

### 3. Service Layer Pattern
Centraliza a lógica de negócio em services, mantendo handlers e repositories simples.

### 4. Middleware Chain
Aplica funcionalidades transversais (logging, recovery) de forma modular.

### 5. Interface Segregation
Interfaces pequenas e focadas, seguindo o princípio SOLID ISP.

## Benefícios da Arquitetura

### ✅ Testabilidade
- Cada camada pode ser testada isoladamente
- Interfaces facilitam criação de mocks
- Lógica de negócio separada de I/O

### ✅ Manutenibilidade
- Responsabilidades bem definidas
- Código organizado e fácil de navegar
- Mudanças localizadas em camadas específicas

### ✅ Escalabilidade
- Fácil adicionar novos features
- Estrutura suporta crescimento
- Camadas independentes

### ✅ Flexibilidade
- Trocar StarkBank por outro provider: apenas alterar repository
- Adicionar cache: injetar no service layer
- Mudar protocolo (HTTP → gRPC): apenas alterar handler layer

## Possíveis Melhorias

### 1. Adicionar Camada de Use Cases
Para aplicações maiores, separar use cases específicos dos services genéricos.

### 2. Implementar Unit of Work
Para gerenciar transações entre múltiplos repositórios.

### 3. Adicionar Event Sourcing
Para auditoria completa de todas as operações.

### 4. Implementar CQRS
Separar comandos (writes) de queries (reads) para melhor performance.

## Conclusão

Esta arquitetura fornece uma base sólida para crescimento futuro, mantendo o código limpo, testável e manutenível. Cada decisão arquitetural foi tomada pensando em:

- **Separação de Concerns**
- **Testabilidade**
- **Manutenibilidade**
- **Escalabilidade**

---

**Autor**: João Pedro Barbosa  
**Challenge**: Stark Bank Backend Developer
