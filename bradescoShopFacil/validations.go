package bradescoShopFacil

import (
	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/validations"
	"strings"
	"time"
)

func bradescoShopFacilValidateAgency(b interface{}) error {
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

func bradescoShopFacilValidateAccount(b interface{}) error {
	switch t := b.(type) {
	case *models.BoletoRequest:
		if t.Agreement.Account == "" {
			return models.NewErrorResponse("MP400", "a conta deve ser preenchida")
		}
		return nil
	default:
		return validations.InvalidType(t)
	}
}

func bradescoShopFacilValidateWallet(b interface{}) error {
	switch t := b.(type) {
	case *models.BoletoRequest:
		if t.Agreement.Wallet != 25 && t.Agreement.Wallet != 26 {
			return models.NewErrorResponse("MP400", "a carteira deve ser 25 ou 26 para o BradescoShopFacil")
		}
		return nil
	default:
		return validations.InvalidType(t)
	}
}

func bradescoShopFacilValidateAuth(b interface{}) error {
	switch t := b.(type) {
	case *models.BoletoRequest:
		usr := strings.TrimSpace(t.Authentication.Username)
		pwd := strings.TrimSpace(t.Authentication.Password)
		if usr == "" || pwd == "" {
			return models.NewErrorResponse("MP400", "o nome de usuário e senha devem ser preenchidos")
		}
		return nil
	default:
		return validations.InvalidType(t)
	}
}

func bradescoShopFacilValidateAgreement(b interface{}) error {
	switch t := b.(type) {
	case *models.BoletoRequest:
		if t.Agreement.AgreementNumber == 0 {
			return models.NewErrorResponse("MP400", "o código do contrato deve ser preenchido")
		}
		return nil
	default:
		return validations.InvalidType(t)
	}
}

func bradescoShopFacilBoletoTypeValidate(b interface{}) error {
	bt := bradescoShopFacilBoletoTypes()

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

// bradescoShopFacilMaxExpirationDateValidate O emissor Bradesco contém um bug na geração da linha digitável onde,
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
