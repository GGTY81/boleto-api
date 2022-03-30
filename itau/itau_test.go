package itau

import (
	"testing"

	"github.com/mundipagg/boleto-api/mock"
	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/test"
	"github.com/stretchr/testify/assert"
)

var boletoTypeParameters = []test.Parameter{
	{Input: models.Title{BoletoType: ""}, Expected: "01"},
	{Input: models.Title{BoletoType: "NSA"}, Expected: "01"},
	{Input: models.Title{BoletoType: "BDP"}, Expected: "18"},
	{Input: models.Title{BoletoType: "DM"}, Expected: "01"},
	{Input: models.Title{BoletoType: "DS"}, Expected: "08"},
	{Input: models.Title{BoletoType: "NP"}, Expected: "02"},
	{Input: models.Title{BoletoType: "RC"}, Expected: "05"},
	{Input: models.Title{BoletoType: "OUT"}, Expected: "99"},
}

func TestProcessBoleto_WhenServiceRespondsSuccessfully_ShouldHasSuccessfulBoletoResponse(t *testing.T) {
	mock.StartMockService("9096")
	input := newStubBoletoRequestItau().Build()
	bank := New()

	output, _ := bank.ProcessBoleto(input)

	test.AssertProcessBoletoWithSuccess(t, output)
}

func TestProcessBoleto_WhenServiceRespondsFailed_ShouldHasFailedBoletoResponse(t *testing.T) {
	mock.StartMockService("9096")
	input := newStubBoletoRequestItau().WithAmountInCents(400).Build()

	bank := New()

	output, err := bank.ProcessBoleto(input)

	assert.Nil(t, err, "Não deve haver um erro")
	test.AssertProcessBoletoFailed(t, output)
}

func TestProcessBoleto_WhenServiceRespondsFailedWithWrongContentAndStatusCodeIs500_ShouldHasFailedBoletoResponseWithWrongContentAndStatusCodeIs500(t *testing.T) {
	mock.StartMockService("9096")
	input := newStubBoletoRequestItau().WithAmountInCents(500).Build()

	bank := New()

	_, errProcessBoleto := bank.ProcessBoleto(input)

	test.AssertError(t, errProcessBoleto, models.BadGatewayError{})
}

func TestProcessBoleto_WhenRequestHasInvalidAccountParameters_ShouldHasFailedBoletoResponse(t *testing.T) {
	mock.StartMockService("9096")
	input := newStubBoletoRequestItau().WithAmountInCents(200).WithAgreementAccount("").Build()

	bank := New()

	output, _ := bank.ProcessBoleto(input)

	test.AssertProcessBoletoFailed(t, output)
}

func TestProcessBoleto_WhenRequestHasInvalidUserNameParameter_ShouldHasFailedBoletoResponse(t *testing.T) {
	mock.StartMockService("9096")
	input := newStubBoletoRequestItau().WithAuthenticationUserName("").WithAmountInCents(200).Build()

	bank := New()

	_, err := bank.ProcessBoleto(input)

	assert.NotNil(t, err, "Deve ocorrer um erro")
}

func TestGetBoletoType_WhenCalled_ShouldBeMapTypeSuccessful(t *testing.T) {
	request := newStubBoletoRequestItau().Build()
	for _, fact := range boletoTypeParameters {
		request.Title = fact.Input.(models.Title)
		_, result := getBoletoType(request)
		assert.Equal(t, fact.Expected, result, "Deve mapear o boleto type corretamente")
	}
}

func TestTemplateResponse_WhenRequestHasSpecialCharacter_ShouldBeParsedSuccessful(t *testing.T) {
	mock.StartMockService("9096")
	input := newStubBoletoRequestItau().WithBuyerName("Usuario \tTeste").Build()
	bank := New()

	output, _ := bank.ProcessBoleto(input)

	test.AssertProcessBoletoWithSuccess(t, output)
}
