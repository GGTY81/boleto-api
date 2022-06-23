package bb

import (
	"fmt"
	"testing"

	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/test"
	"github.com/stretchr/testify/assert"
)

var boletoRequestParameters = []test.Parameter{
	{
		Line:     0,
		Input:    newStubBoletoRequestBB().Build(),
		Expected: nil,
	},
	{
		Line:     1,
		Input:    newStubBoletoRequestBB().WithTitle(models.Title{DocumentNumber: "loja"}).Build(),
		Expected: nil,
	},
	{
		Line:     2,
		Input:    newStubBoletoRequestBB().WithTitle(models.Title{DocumentNumber: "123456"}).Build(),
		Expected: nil,
	},
	{
		Line:     3,
		Input:    newStubBoletoRequestBB().WithTitle(models.Title{DocumentNumber: "lojas-15-digito"}).Build(),
		Expected: nil,
	},
	{
		Line:     4,
		Input:    newStubBoletoRequestBB().WithTitle(models.Title{DocumentNumber: "lojas-16-digitos"}).Build(),
		Expected: models.NewErrorResponse("MP400", "O campo documentNumber do título ultrapassou o limite permitido de 15 caracteres"),
	},
}

func Test_GivenTheValidateTitleDocumentNumberMethodWasCalled_ThenItShouldCorrectlyValidateTheField(t *testing.T) {
	for _, fact := range boletoRequestParameters {
		boletoRequest := fact.Input.(*models.BoletoRequest)

		result := bbValidateTitleDocumentNumber(boletoRequest)

		assert.Equal(t, fact.Expected, result, fmt.Sprintf("bbValidateTitleDocumentNumber - Linha %d: Deve validar o campo documentNumber corretamente", fact.Line))
	}
}
