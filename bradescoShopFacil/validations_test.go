package bradescoShopFacil

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
	{Input: "123456789", Expected: true},
	{Input: "  ", Expected: true},
	{Input: "", Expected: false},
}

var walletParameters = []test.Parameter{
	{Input: uint16(25), Expected: true},
	{Input: uint16(26), Expected: true},
	{Input: uint16(10), Expected: false},
	{Input: uint16(20), Expected: false},
	{Input: uint16(40), Expected: false},
	{Input: uint16(50), Expected: false},
}

var authenticationParameters = []test.Parameter{
	{Input: models.Authentication{Username: "Usuario", Password: "Senha"}, Expected: true},
	{Input: models.Authentication{Username: "   ", Password: "Senha"}, Expected: false},
	{Input: models.Authentication{Username: "", Password: "Senha"}, Expected: false},
	{Input: models.Authentication{Username: "Usuario", Password: "   "}, Expected: false},
	{Input: models.Authentication{Username: "Usuario", Password: ""}, Expected: false},
	{Input: models.Authentication{Username: "", Password: ""}, Expected: false},
}

var agreementNumberParameters = []test.Parameter{
	{Input: uint(1), Expected: true},
	{Input: uint(4), Expected: true},
	{Input: uint(0), Expected: false},
}

var boletoTypeKeyParameters = []test.Parameter{
	{Input: models.Title{BoletoType: ""}, Expected: true},
	{Input: models.Title{BoletoType: "DM"}, Expected: true},
	{Input: models.Title{BoletoType: "NP"}, Expected: true},
	{Input: models.Title{BoletoType: "RC"}, Expected: true},
	{Input: models.Title{BoletoType: "DS"}, Expected: true},
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
		request := newStubBoletoRequestBradescoShopFacil().WithAgreementAgency(fact.Input.(string)).Build()
		result := bradescoShopFacilValidateAgency(request)
		assert.Equal(t, fact.Expected, result == nil, fmt.Sprintf("O resultado: %d não condiz com o esperado: %d, utilizando o input: %d", result, fact.Expected, fact.Input))
	}
}

func Test_AgencyValidation_WhenTypeIsInvalid(t *testing.T) {

	request := "Não é um boleto request"
	result := bradescoShopFacilValidateAgency(request)
	assert.IsType(t, models.ErrorResponse{}, result, fmt.Sprintf("O tipo do resultado: %T não condiz com o tipo esperado: %T", result, models.ErrorResponse{}))
}

func Test_AccountValidation_WhenTypeIsBoletoRequest(t *testing.T) {

	for _, fact := range accountNumberParameters {
		request := newStubBoletoRequestBradescoShopFacil().WithAgreementAccount(fact.Input.(string)).Build()
		result := bradescoShopFacilValidateAccount(request)
		assert.Equal(t, fact.Expected, result == nil, fmt.Sprintf("O resultado: %d não condiz com o esperado: %d, utilizando o input: %d", result, fact.Expected, fact.Input))
	}
}

func Test_AccountValidation_WhenTypeIsInvalid(t *testing.T) {

	request := "Não é um boleto request"
	result := bradescoShopFacilValidateAccount(request)
	assert.IsType(t, models.ErrorResponse{}, result, fmt.Sprintf("O tipo do resultado: %T não condiz com o tipo esperado: %T", result, models.ErrorResponse{}))
}

func Test_WalletValidation_WhenTypeIsBoletoRequest(t *testing.T) {

	for _, fact := range walletParameters {
		request := newStubBoletoRequestBradescoShopFacil().WithWallet(fact.Input.(uint16)).Build()
		result := bradescoShopFacilValidateWallet(request)
		assert.Equal(t, fact.Expected, result == nil, fmt.Sprintf("O resultado: %d não condiz com o esperado: %d, utilizando o input: %d", result, fact.Expected, fact.Input))
	}
}

func Test_WalletValidation_WhenTypeIsInvalid(t *testing.T) {

	request := "Não é um boleto request"
	result := bradescoShopFacilValidateWallet(request)
	assert.IsType(t, models.ErrorResponse{}, result, fmt.Sprintf("O tipo do resultado: %T não condiz com o tipo esperado: %T", result, models.ErrorResponse{}))
}

func Test_AuthenticationValidation_WhenTypeIsBoletoRequest(t *testing.T) {
	for _, fact := range authenticationParameters {
		request := newStubBoletoRequestBradescoShopFacil().WithAuthentication(fact.Input.(models.Authentication)).Build()
		result := bradescoShopFacilValidateAuth(request)
		assert.Equal(t, fact.Expected, result == nil, fmt.Sprintf("O resultado: %d não condiz com o esperado: %d, utilizando o input: %d", result, fact.Expected, fact.Input))
	}
}

func Test_AuthenticationValidation_WhenTypeIsInvalid(t *testing.T) {

	request := "Não é um boleto request"
	result := bradescoShopFacilValidateAuth(request)
	assert.IsType(t, models.ErrorResponse{}, result, fmt.Sprintf("O tipo do resultado: %T não condiz com o tipo esperado: %T", result, models.ErrorResponse{}))
}

func Test_AgreementNumberValidation_WhenTypeIsBoletoRequest(t *testing.T) {
	for _, fact := range agreementNumberParameters {
		request := newStubBoletoRequestBradescoShopFacil().WithAgreementNumber(fact.Input.(uint)).Build()
		result := bradescoShopFacilValidateAgreement(request)
		assert.Equal(t, fact.Expected, result == nil, fmt.Sprintf("O resultado: %d não condiz com o esperado: %d, utilizando o input: %d", result, fact.Expected, fact.Input))
	}
}

func Test_AgreementNumberValidation_WhenTypeIsInvalid(t *testing.T) {

	request := "Não é um boleto request"
	result := bradescoShopFacilValidateAgreement(request)
	assert.IsType(t, models.ErrorResponse{}, result, fmt.Sprintf("O tipo do resultado: %T não condiz com o tipo esperado: %T", result, models.ErrorResponse{}))
}

func Test_BoletoTypeValidation_WhenTypeIsBoletoRequest(t *testing.T) {

	for _, fact := range boletoTypeKeyParameters {
		request := newStubBoletoRequestBradescoShopFacil().WithBoletoType(fact.Input.(models.Title)).Build()
		result := bradescoShopFacilBoletoTypeValidate(request)
		assert.Equal(t, fact.Expected, result == nil, fmt.Sprintf("O resultado: %d não condiz com o esperado: %d, utilizando o input: %d", result, fact.Expected, fact.Input))
	}
}

func Test_BoletoTypeValidation_WhenTypeIsInvalid(t *testing.T) {

	request := "Não é um boleto request"
	result := bradescoShopFacilBoletoTypeValidate(request)
	assert.IsType(t, models.ErrorResponse{}, result, fmt.Sprintf("O tipo do resultado: %T não condiz com o tipo esperado: %T", result, models.ErrorResponse{}))
}

func Test_MaxDateExpirationBlockValidation_WhenTypeIsBoletoRequest(t *testing.T) {

	for _, fact := range expirationDateParameters {
		expDate, _ := time.Parse("2006-01-02", fact.Input.(string))
		request := newStubBoletoRequestBradescoShopFacil().WithExpirationDate(expDate).Build()
		result := validations.ValidateMaxExpirationDate(request)
		assert.Equal(t, fact.Expected, result == nil, fmt.Sprintf("O resultado: %d não condiz com o esperado: %d, utilizando o input: %d", result, fact.Expected, fact.Input))
	}
}

func Test_MaxDateExpirationBlockValidation_WhenTypeIsInvalid(t *testing.T) {

	request := "Não é um boleto request"
	result := validations.ValidateMaxExpirationDate(request)
	assert.IsType(t, models.ErrorResponse{}, result, fmt.Sprintf("O tipo do resultado: %T não condiz com o tipo esperado: %T", result, models.ErrorResponse{}))
}
