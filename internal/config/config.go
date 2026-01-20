package config

import (
	"fmt"
	"os"
)

// Config armazena todas as configurações da aplicação
type Config struct {
	Server      ServerConfig
	StarkBank   StarkBankConfig
	Destination DestinationAccount
}

// ServerConfig configurações do servidor HTTP
type ServerConfig struct {
	Port string
	Host string
}

// StarkBankConfig configurações da StarkBank
type StarkBankConfig struct {
	ProjectID   string
	PrivateKey  string
	Environment string
}

// DestinationAccount conta de destino para transferências
type DestinationAccount struct {
	BankCode      string
	BranchCode    string
	AccountNumber string
	Name          string
	TaxID         string
	AccountType   string
}

// Load carrega as configurações da aplicação
func Load() (*Config, error) {
	privateKey, err := loadPrivateKey()
	if err != nil {
		return nil, err
	}

	// Validar variáveis obrigatórias
	projectID := os.Getenv("STARK_PROJECT_ID")
	if projectID == "" {
		return nil, fmt.Errorf("STARK_PROJECT_ID é obrigatório. Configure a variável de ambiente")
	}

	port := getEnv("PORT", "8080")
	environment := getEnv("STARK_ENVIRONMENT", "sandbox")

	return &Config{
		Server: ServerConfig{
			Port: port,
			Host: "0.0.0.0",
		},
		StarkBank: StarkBankConfig{
			ProjectID:   projectID,
			PrivateKey:  privateKey,
			Environment: environment,
		},
		Destination: DestinationAccount{
			BankCode:      "20018183",
			BranchCode:    "0001",
			AccountNumber: "6341320293482496",
			Name:          "Stark Bank S.A.",
			TaxID:         "20.018.183/0001-80",
			AccountType:   "payment",
		},
	}, nil
}

// loadPrivateKey carrega a chave privada do arquivo ou variável de ambiente
func loadPrivateKey() (string, error) {
	privateKey := os.Getenv("PRIVATE_KEY")
	if privateKey != "" {
		return privateKey, nil
	}

	content, err := os.ReadFile("privateKeyChallenge.pem")
	if err != nil {
		return "", fmt.Errorf("erro ao ler chave privada: %w", err)
	}

	return string(content), nil
}

// getEnv retorna o valor de uma variável de ambiente ou um valor padrão
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
