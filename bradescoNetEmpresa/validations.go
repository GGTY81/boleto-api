package bradescoNetEmpresa

import (
	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/validations"
	"time"
)

func bradescoNetEmpresaValidateAgency(b interface{}) error {
	switch t := b.(type) {
	case *models.BoletoRequest:
		err := t.Agreement.IsAgencyValid()
		if err != nil {
			return models.NewErrorResponse("MP400", err.Error())
		}
		return nil
	default:
		return validations.InvalidType(t)
	}
}

func bradescoNetEmpresaValidateAccount(b interface{}) error {
	switch t := b.(type) {
	case *models.BoletoRequest:
		err := t.Agreement.IsAccountValid(7)
		if err != nil {
			return models.NewErrorResponse("MP400", err.Error())
		}
		return nil
	default:
		return validations.InvalidType(t)
	}
}

func bradescoNetEmpresaValidateWallet(b interface{}) error {
	switch t := b.(type) {
	case *models.BoletoRequest:
		if t.Agreement.Wallet != 4 && t.Agreement.Wallet != 9 && t.Agreement.Wallet != 19 {
			return models.NewErrorResponse("MP400", "a carteira deve ser 4, 9 ou 19 para o bradescoNetEmpresa")
		}
		return nil
	default:
		return validations.InvalidType(t)
	}
}

func bradescoNetEmpresaBoletoTypeValidate(b interface{}) error {
	bt := bradescoNetEmpresaBoletoTypes()

	switch t := b.(type) {

	case *models.BoletoRequest:
		if len(t.Title.BoletoType) > 0 && bt[t.Title.BoletoType] == "" {
			return models.NewErrorResponse("MP400", "espécie de boleto informada não existente")
		}
		return nil
	default:
		return validations.InvalidType(t)
	}
}

// bradescoNetEmpresaMaxExpirationDateValidate O emissor Bradesco contém um bug na geração da linha digitável onde,
// quando a data de vencimento é maior do que 21-02-2025 a linha digitável se torna inválida(O própio Bradesco não consegue ler a linha gerada) e não conseguimos gerar a visualização do boleto
// Para evitarmos esse problema, adicionamos temporariamente essa trava que bloqueia a geração de boletos com data de vencimento após a data em questão.
func ValidateMaxExpirationDate(b interface{}) error {
	maxExpDate, _ := time.Parse("2006-01-02", "2025-02-21")

	switch t := b.(type) {
	case *models.BoletoRequest:
		if t.Title.ExpireDateTime.After(maxExpDate) {
			return models.NewErrorResponse("MPExpireDate", "Data de vencimento não pode ser maior que 21-02-2025")
		}
		return nil
	default:
		return validations.InvalidType(t)
	}
}
