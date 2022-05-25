package pefisa

import (
	"testing"

	"github.com/mundipagg/boleto-api/mock"
	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/test"
	"github.com/stretchr/testify/assert"
)

var boletoTypeParameters = []test.Parameter{
	{Input: models.Title{BoletoType: ""}, Expected: "1"},
	{Input: models.Title{BoletoType: "NSA"}, Expected: "1"},
	{Input: models.Title{BoletoType: "DM"}, Expected: "1"},
	{Input: models.Title{BoletoType: "DS"}, Expected: "2"},
	{Input: models.Title{BoletoType: "NP"}, Expected: "3"},
	{Input: models.Title{BoletoType: "SE"}, Expected: "4"},
	{Input: models.Title{BoletoType: "CH"}, Expected: "10"},
	{Input: models.Title{BoletoType: "OUT"}, Expected: "99"},
}

func TestProcessBoleto_WhenServiceRespondsSuccessfully_ShouldHasSuccessfulBoletoResponse(t *testing.T) {
	mock.StartMockService("9093")

	input := newStubBoletoRequestPefisa().Build()
	bank := New()

	output, _ := bank.ProcessBoleto(input)

	test.AssertProcessBoletoWithSuccess(t, output)
}

func TestProcessBoleto_WhenServiceRespondsFailed_ShouldHasFailedBoletoResponse(t *testing.T) {
	mock.StartMockService("9093")

	input := newStubBoletoRequestPefisa().WithAmountInCents(201).Build()
	bank := New()

	output, _ := bank.ProcessBoleto(input)

	test.AssertProcessBoletoFailed(t, output)
}

func TestGetBoletoType_WhenCalled_ShouldBeMapTypeSuccessful(t *testing.T) {
	request := newStubBoletoRequestPefisa().Build()
	for _, fact := range boletoTypeParameters {
		request.Title = fact.Input.(models.Title)
		_, result := getBoletoType(request)
		assert.Equal(t, fact.Expected, result, "Deve mapear o boleto type corretamente")
	}
}

func TestTemplateResponse_WhenRequestHasTabCharacter_ShouldBeParsedSuccessful(t *testing.T) {
	mock.StartMockService("9093")
	input := newStubBoletoRequestPefisa().WithBuyerName("Usuario \tTeste").Build()
	bank := New()

	output, _ := bank.ProcessBoleto(input)

	test.AssertProcessBoletoWithSuccess(t, output)
}
