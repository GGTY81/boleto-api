package bradescoNetEmpresa

import (
	"fmt"
	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/test"
	"github.com/mundipagg/boleto-api/validations"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var agencyNumberParameters = []test.Parameter{
	{Input: "1", Expected: true},
	{Input: "12", Expected: true},
	{Input: "123", Expected: true},
	{Input: "1234", Expected: true},
	{Input: "123A", Expected: true},
	{Input: "", Expected: false},
	{Input: "  ", Expected: false},
	{Input: "123456789", Expected: false},
}

var accountNumberParameters = []test.Parameter{
	{Input: "1", Expected: true},
	{Input: "12", Expected: true},
	{Input: "123", Expected: true},
	{Input: "1234", Expected: true},
	{Input: "123A", Expected: true},
	{Input: "123456789", Expected: false},
	{Input: "  ", Expected: false},
	{Input: "", Expected: false},
}

var walletParameters = []test.Parameter{
	{Input: uint16(4), Expected: true},
	{Input: uint16(9), Expected: true},
	{Input: uint16(19), Expected: true},
	{Input: uint16(1), Expected: false},
	{Input: uint16(6), Expected: false},
	{Input: uint16(20), Expected: false},
}

var boletoTypeKeyParameters = []test.Parameter{
	{Input: models.Title{BoletoType: ""}, Expected: true},
	{Input: models.Title{BoletoType: "CH"}, Expected: true},
	{Input: models.Title{BoletoType: "DM"}, Expected: true},
	{Input: models.Title{BoletoType: "DS"}, Expected: true},
	{Input: models.Title{BoletoType: "NP"}, Expected: true},
	{Input: models.Title{BoletoType: "RC"}, Expected: true},
	{Input: models.Title{BoletoType: "OUT"}, Expected: true},
	{Input: models.Title{BoletoType: "A"}, Expected: false},
	{Input: models.Title{BoletoType: "ABC"}, Expected: false},
}

var expirationDateParameters = []test.Parameter{
	{Input: "2025-01-01", Expected: true},
	{Input: "2025-02-20", Expected: true},
	{Input: "2025-02-21", Expected: true},
	{Input: "2025-02-22", Expected: false},
	{Input: "2025-12-31", Expected: false},
}

func Test_AgencyValidation_WhenTypeIsBoletoRequest(t *testing.T) {

	for _, fact := range agencyNumberParameters {
		request := newStubBoletoRequestBradescoNetEmpresa().WithAgreementAgency(fact.Input.(string)).Build()
		result := bradescoNetEmpresaValidateAgency(request)
		assert.Equal(t, fact.Expected, result == nil, fmt.Sprintf("O resultado: %d não condiz com o esperado: %d, utilizando o input: %d", result, fact.Expected, fact.Input))
	}
}

func Test_AgencyValidation_WhenTypeIsInvalid(t *testing.T) {

	request := "Não é um boleto request"
	result := bradescoNetEmpresaValidateAgency(request)
	assert.IsType(t, models.ErrorResponse{}, result, fmt.Sprintf("O tipo do resultado: %T não condiz com o tipo esperado: %T", result, models.ErrorResponse{}))
}

func Test_AccountValidation_WhenTypeIsBoletoRequest(t *testing.T) {

	for _, fact := range accountNumberParameters {
		request := newStubBoletoRequestBradescoNetEmpresa().WithAgreementAccount(fact.Input.(string)).Build()
		result := bradescoNetEmpresaValidateAccount(request)
		assert.Equal(t, fact.Expected, result == nil, fmt.Sprintf("O resultado: %d não condiz com o esperado: %d, utilizando o input: %d", result, fact.Expected, fact.Input))
	}
}

func Test_AccountValidation_WhenTypeIsInvalid(t *testing.T) {

	request := "Não é um boleto request"
	result := bradescoNetEmpresaValidateAccount(request)
	assert.IsType(t, models.ErrorResponse{}, result, fmt.Sprintf("O tipo do resultado: %T não condiz com o tipo esperado: %T", result, models.ErrorResponse{}))
}

func Test_WalletValidation_WhenTypeIsBoletoRequest(t *testing.T) {

	for _, fact := range walletParameters {
		request := newStubBoletoRequestBradescoNetEmpresa().WithWallet(fact.Input.(uint16)).Build()
		result := bradescoNetEmpresaValidateWallet(request)
		assert.Equal(t, fact.Expected, result == nil, fmt.Sprintf("O resultado: %d não condiz com o esperado: %d, utilizando o input: %d", result, fact.Expected, fact.Input))
	}
}

func Test_WalletValidation_WhenTypeIsInvalid(t *testing.T) {

	request := "Não é um boleto request"
	result := bradescoNetEmpresaValidateWallet(request)
	assert.IsType(t, models.ErrorResponse{}, result, fmt.Sprintf("O tipo do resultado: %T não condiz com o tipo esperado: %T", result, models.ErrorResponse{}))
}

func Test_BoletoTypeValidation_WhenTypeIsBoletoRequest(t *testing.T) {

	for _, fact := range boletoTypeKeyParameters {
		request := newStubBoletoRequestBradescoNetEmpresa().WithBoletoType(fact.Input.(models.Title)).Build()
		result := bradescoNetEmpresaBoletoTypeValidate(request)
		assert.Equal(t, fact.Expected, result == nil, fmt.Sprintf("O resultado: %d não condiz com o esperado: %d, utilizando o input: %d", result, fact.Expected, fact.Input))
	}
}

func Test_BoletoTypeValidation_WhenTypeIsInvalid(t *testing.T) {

	request := "Não é um boleto request"
	result := bradescoNetEmpresaBoletoTypeValidate(request)
	assert.IsType(t, models.ErrorResponse{}, result, fmt.Sprintf("O tipo do resultado: %T não condiz com o tipo esperado: %T", result, models.ErrorResponse{}))
}

func Test_MaxDateExpirationBlockValidation_WhenTypeIsBoletoRequest(t *testing.T) {

	for _, fact := range expirationDateParameters {
		expDate, _ := time.Parse("2006-01-02", fact.Input.(string))
		request := newStubBoletoRequestBradescoNetEmpresa().WithExpirationDate(expDate).Build()
		result := validations.ValidateMaxExpirationDate(request)
		assert.Equal(t, fact.Expected, result == nil, fmt.Sprintf("O resultado: %d não condiz com o esperado: %d, utilizando o input: %d", result, fact.Expected, fact.Input))
	}
}

func Test_MaxDateExpirationBlockValidation_WhenTypeIsInvalid(t *testing.T) {

	request := "Não é um boleto request"
	result := validations.ValidateMaxExpirationDate(request)
	assert.IsType(t, models.ErrorResponse{}, result, fmt.Sprintf("O tipo do resultado: %T não condiz com o tipo esperado: %T", result, models.ErrorResponse{}))
}
