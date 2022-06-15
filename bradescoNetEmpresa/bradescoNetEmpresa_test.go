package bradescoNetEmpresa

import (
	"testing"

	"github.com/mundipagg/boleto-api/mock"
	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/test"
	"github.com/stretchr/testify/assert"
)

var boletoTypeParameters = []test.Parameter{
	{Input: models.Title{BoletoType: ""}, Expected: "02"},
	{Input: models.Title{BoletoType: "NSA"}, Expected: "02"},
	{Input: models.Title{BoletoType: "BDP"}, Expected: "02"},
	{Input: models.Title{BoletoType: "CH"}, Expected: "01"},
	{Input: models.Title{BoletoType: "DM"}, Expected: "02"},
	{Input: models.Title{BoletoType: "DS"}, Expected: "04"},
	{Input: models.Title{BoletoType: "NP"}, Expected: "12"},
	{Input: models.Title{BoletoType: "RC"}, Expected: "17"},
	{Input: models.Title{BoletoType: "OUT"}, Expected: "99"},
}

func TestProcessBoleto_WhenServiceRespondsSuccessfully_ShouldHasSuccessfulBoletoResponse(t *testing.T) {
	mock.StartMockService("9012")
	input := newStubBoletoRequestBradescoNetEmpresa().WithAmountInCents(200).Build()
	bank := New()

	output, _ := bank.ProcessBoleto(input)

	test.AssertProcessBoletoWithSuccess(t, output)
}

func TestProcessBoleto_WhenServiceRespondsFailed_ShouldHasFailedBoletoResponse(t *testing.T) {
	mock.StartMockService("9011")
	input := newStubBoletoRequestBradescoNetEmpresa().WithAmountInCents(201).Build()
	bank := New()

	output, _ := bank.ProcessBoleto(input)

	test.AssertProcessBoletoFailed(t, output)
}

func TestProcessBoleto_WhenServiceRespondsCertificateFailed_ShouldHasFailedBoletoResponse(t *testing.T) {
	mock.StartMockService("9010")
	input := newStubBoletoRequestBradescoNetEmpresa().WithAmountInCents(202).Build()
	bank := New()

	output, _ := bank.ProcessBoleto(input)

	test.AssertProcessBoletoFailed(t, output)
}

func TestGetBoletoType_WhenCalled_ShouldBeMapTypeSuccessful(t *testing.T) {
	BradescoNetEmpresaRequestStub := newStubBoletoRequestBradescoNetEmpresa().Build()
	for _, fact := range boletoTypeParameters {
		BradescoNetEmpresaRequestStub.Title = fact.Input.(models.Title)
		_, result := getBoletoType(BradescoNetEmpresaRequestStub)
		assert.Equal(t, fact.Expected, result, "Deve mapear o boleto type corretamente")
	}
}

func TestTemplateResponse_WhenRequestHasSpecialCharacter_ShouldBeParsedSuccessful(t *testing.T) {
	mock.StartMockService("9013")
	input := newStubBoletoRequestBradescoNetEmpresa().WithAmountInCents(204).WithBuyerName("Usuario 	Teste").Build()
	bank := New()

	output, _ := bank.ProcessBoleto(input)

	test.AssertProcessBoletoWithSuccess(t, output)
}
