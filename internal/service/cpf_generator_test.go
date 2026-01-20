package service

import (
	"strconv"
	"testing"
)

func TestGenerateCPF(t *testing.T) {
	// Gerar 100 CPFs e validar todos
	for i := 0; i < 100; i++ {
		cpf := GenerateCPF()

		// Verificar tamanho
		if len(cpf) != 11 {
			t.Errorf("CPF gerado tem tamanho inválido: %d (esperado 11)", len(cpf))
		}

		// Verificar se é apenas números
		for _, char := range cpf {
			if char < '0' || char > '9' {
				t.Errorf("CPF contém caracteres não numéricos: %s", cpf)
			}
		}

		// Validar CPF usando o algoritmo
		if !isValidCPF(cpf) {
			t.Errorf("CPF gerado é inválido: %s", cpf)
		}
	}
}

func TestGenerateCPFUniqueness(t *testing.T) {
	// Gerar múltiplos CPFs e verificar que não são todos iguais
	cpfs := make(map[string]bool)
	for i := 0; i < 50; i++ {
		cpf := GenerateCPF()
		cpfs[cpf] = true
	}

	// Deve haver pelo menos alguns CPFs diferentes (não todos iguais)
	if len(cpfs) < 10 {
		t.Errorf("CPFs gerados não são suficientemente únicos: apenas %d únicos em 50", len(cpfs))
	}
}

// isValidCPF valida um CPF usando o algoritmo oficial
func isValidCPF(cpf string) bool {
	if len(cpf) != 11 {
		return false
	}

	// Converter string para slice de ints
	digits := make([]int, 11)
	for i, char := range cpf {
		digit, err := strconv.Atoi(string(char))
		if err != nil {
			return false
		}
		digits[i] = digit
	}

	// Verificar se todos os dígitos são iguais (inválido)
	allEqual := true
	for i := 1; i < 11; i++ {
		if digits[i] != digits[0] {
			allEqual = false
			break
		}
	}
	if allEqual {
		return false
	}

	// Validar primeiro dígito verificador
	sum := 0
	for i := 0; i < 9; i++ {
		sum += digits[i] * (10 - i)
	}
	remainder := sum % 11
	expectedDigit1 := 0
	if remainder >= 2 {
		expectedDigit1 = 11 - remainder
	}
	if digits[9] != expectedDigit1 {
		return false
	}

	// Validar segundo dígito verificador
	sum = 0
	for i := 0; i < 10; i++ {
		sum += digits[i] * (11 - i)
	}
	remainder = sum % 11
	expectedDigit2 := 0
	if remainder >= 2 {
		expectedDigit2 = 11 - remainder
	}
	if digits[10] != expectedDigit2 {
		return false
	}

	return true
}

func BenchmarkGenerateCPF(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateCPF()
	}
}
