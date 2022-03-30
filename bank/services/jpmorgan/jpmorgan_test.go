package jpmorgan

import (
	"testing"
	"time"

	"github.com/mundipagg/boleto-api/certificate"
	"github.com/mundipagg/boleto-api/mock"
	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/test"
	"github.com/stretchr/testify/assert"
)

var boletoResponseFailParameters = []test.Parameter{
	{Input: newStubBoletoRequestJPMorgan().WithAmountInCents(0).Build(), Expected: models.ErrorResponse{Code: `MPAmountInCents`, Message: `Valor não pode ser menor do que 1 centavo`}},
	{Input: newStubBoletoRequestJPMorgan().WithExpirationDate(time.Date(2021, 10, 17, 12, 12, 12, 12, time.Local)).Build(), Expected: models.ErrorResponse{Code: `MPExpireDate`, Message: `Data de expiração não pode ser menor que a data de hoje`}},
	{Input: newStubBoletoRequestJPMorgan().WithAmountInCents(201).Build(), Expected: models.ErrorResponse{Code: `Internal Error - Contact Service Provider`, Message: `Internal Error - Contact Service Provider`}},
	{Input: newStubBoletoRequestJPMorgan().WithAmountInCents(202).Build(), Expected: models.ErrorResponse{Code: `GCA-010`, Message: `The account was not found.`}},
	{Input: newStubBoletoRequestJPMorgan().WithAmountInCents(205).Build(), Expected: models.ErrorResponse{Code: `MPOurNumberFail`, Message: `our number was not returned by the bank`}},
	{Input: newStubBoletoRequestJPMorgan().WithAmountInCents(203).Build(), Expected: models.ErrorResponse{Code: `BOL-3`, Message: `Não foi possível processar a requisição por inconsistencia nos campos abaixo`}},
	{Input: newStubBoletoRequestJPMorgan().WithAmountInCents(204).Build(), Expected: models.ErrorResponse{Code: `BOL-144`, Message: `Boleto já Existente. Detalhes abaixo:`}},
	{Input: newStubBoletoRequestJPMorgan().WithAmountInCents(206).Build(), Expected: models.ErrorResponse{Code: `GCA-111`, Message: `codContaCorrente/codAgencia is required`}},
	{Input: newStubBoletoRequestJPMorgan().WithAmountInCents(207).Build(), Expected: models.ErrorResponse{Code: `Signature Verification Failure`, Message: `Signature Verification Failure`}},
	{Input: newStubBoletoRequestJPMorgan().WithAmountInCents(209).Build(), Expected: models.ErrorResponse{Code: `BOL-144`, Message: `Boleto já Existente. Detalhes abaixo:`}},
	{Input: newStubBoletoRequestJPMorgan().WithAmountInCents(300).Build(), Expected: models.ErrorResponse{Code: `MPTimeout`, Message: `GatewayTimeout`}},
}

func Test_ProcessBoleto_WhenServiceRespondsSuccessfully_ShouldHasSuccessfulBoletoResponse(t *testing.T) {
	mock.StartMockService("9102")
	certificate.LoadMockCertificates()
	input := newStubBoletoRequestJPMorgan().Build()
	bank, _ := New()

	output, err := bank.ProcessBoleto(input)

	assert.Nil(t, err, "Não deve haver um erro")
	assert.Equal(t, 12, len(output.OurNumber))
	assert.Equal(t, "123456789012", output.OurNumber)
	test.AssertProcessBoletoWithSuccess(t, output)
}

func Test_ProcessBoleto_WhenServiceRespondsWithShortOurNUmber_ShouldHasPadZerosAndSuccessfulBoletoResponse(t *testing.T) {
	mock.StartMockService("9102")
	certificate.LoadMockCertificates()
	input := newStubBoletoRequestJPMorgan().WithAmountInCents(211).Build()
	bank, _ := New()

	output, err := bank.ProcessBoleto(input)

	assert.Nil(t, err, "Não deve haver um erro")
	assert.Equal(t, 12, len(output.OurNumber))
	assert.Equal(t, "000000123456", output.OurNumber)
	test.AssertProcessBoletoWithSuccess(t, output)
}

func Test_ProcessBoleto_WhenServiceRespondsUnsuccessful_ShouldHasErrorResponse(t *testing.T) {
	mock.StartMockService("9102")
	certificate.LoadMockCertificates()
	bank, _ := New()

	for _, fact := range boletoResponseFailParameters {
		request := fact.Input.(*models.BoletoRequest)
		response, err := bank.ProcessBoleto(request)
		assert.Nil(t, err, "Não deve haver um erro fora do objeto de response")

		test.AssertProcessBoletoFailed(t, response)
		assert.Equal(t, fact.Expected.(models.ErrorResponse).Code, response.Errors[0].Code)
		assert.Contains(t, response.Errors[0].Message, fact.Expected.(models.ErrorResponse).Message)
	}
}

func TestTemplateResponse_WhenRequestHasSpecialCharacter_ShouldBeParsedSuccessful(t *testing.T) {
	mock.StartMockService("9092")
	certificate.LoadMockCertificates()
	input := newStubBoletoRequestJPMorgan().WithBuyerName("Nome do \tComprador (Cliente)").Build()
	bank, _ := New()

	output, _ := bank.ProcessBoleto(input)

	test.AssertProcessBoletoWithSuccess(t, output)
}
