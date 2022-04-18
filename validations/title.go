package validations

import (
	"strconv"
	"time"

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

//ValidateMaxExpirationDate O emissor Bradesco contém um bug na geração da linha digitável onde,
// quando a data de vencimento é maior do que 21-02-2025 a linha digitável se torna inválida(O própio Bradesco não consegue ler a linha gerada) e não conseguimos gerar a visualização do boleto
// Para evitarmos esse problema, adicionamos temporariamente essa trava que bloqueia a geração de boletos com data de vencimento após a data em questão.
func ValidateMaxExpirationDate(b interface{}) error {
	maxExpDate, _ := time.Parse("2006-01-02", "2025-02-21")

	switch t := b.(type) {
	case *models.BoletoRequest:
		if t.Title.ExpireDateTime.After(maxExpDate) {
			return models.NewErrorResponse("MPExpireDate", "Data de vencimento não pode ser maior que 21/02/2025")
		}
		return nil
	default:
		return InvalidType(t)
	}
}
