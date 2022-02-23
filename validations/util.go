package validations

import (
	"strconv"

	"github.com/mundipagg/boleto-api/models"
)

func SumAccountDigits(a string, m []int) int {
	sum := 0
	for idx, c := range a {
		i, _ := strconv.Atoi(string(c))
		sum += i * m[idx]
	}
	return sum
}

func InvalidType(t interface{}) error {
	return models.NewErrorResponse("MP500", "Tipo inválido")
}

func ModElevenCalculator(a string, m []int) string {
	sum := SumAccountDigits(a, m)

	digit := 11 - sum%11

	if digit == 10 {
		return "X"
	}

	if digit == 11 {
		return "0"
	}
	return strconv.Itoa(digit)
}

//ValidatePayeeGuarantorName Verifica se o nome do lojista é existe
func ValidatePayeeGuarantorName(b interface{}) error {
	switch t := b.(type) {
	case *models.BoletoRequest:
		if t.HasPayeeGuarantor() {
			if !t.PayeeGuarantor.HasName() {
				return models.NewErrorResponse("MPPayeeGuarantorNameType", "Nome do sacador avalista está vazio")
			}
		}
		return nil
	default:
		return InvalidType(t)
	}
}
