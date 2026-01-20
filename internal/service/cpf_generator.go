package service

import (
	"fmt"
	"math/rand"
	"time"
)

// GenerateCPF gera um CPF válido aleatório
func GenerateCPF() string {
	rand.Seed(time.Now().UnixNano())

	// Gerar 9 primeiros dígitos aleatórios
	cpf := make([]int, 11)
	for i := 0; i < 9; i++ {
		cpf[i] = rand.Intn(10)
	}

	// Calcular primeiro dígito verificador
	sum := 0
	for i := 0; i < 9; i++ {
		sum += cpf[i] * (10 - i)
	}
	remainder := sum % 11
	if remainder < 2 {
		cpf[9] = 0
	} else {
		cpf[9] = 11 - remainder
	}

	// Calcular segundo dígito verificador
	sum = 0
	for i := 0; i < 10; i++ {
		sum += cpf[i] * (11 - i)
	}
	remainder = sum % 11
	if remainder < 2 {
		cpf[10] = 0
	} else {
		cpf[10] = 11 - remainder
	}

	// Converter para string (sem formatação)
	return fmt.Sprintf("%d%d%d%d%d%d%d%d%d%d%d",
		cpf[0], cpf[1], cpf[2], cpf[3], cpf[4],
		cpf[5], cpf[6], cpf[7], cpf[8], cpf[9], cpf[10])
}
