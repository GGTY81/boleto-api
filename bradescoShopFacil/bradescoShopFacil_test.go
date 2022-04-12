package bradescoShopFacil

import (
	"testing"
	"time"

	"github.com/mundipagg/boleto-api/mock"
	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/test"
	"github.com/stretchr/testify/assert"
)

var boletoTypeParameters = []test.Parameter{
	{Input: models.Title{BoletoType: ""}, Expected: "01"},
	{Input: models.Title{BoletoType: "NSA"}, Expected: "01"},
	{Input: models.Title{BoletoType: "BDP"}, Expected: "01"},
	{Input: models.Title{BoletoType: "DM"}, Expected: "01"},
	{Input: models.Title{BoletoType: "DS"}, Expected: "12"},
	{Input: models.Title{BoletoType: "NP"}, Expected: "02"},
	{Input: models.Title{BoletoType: "RC"}, Expected: "05"},
	{Input: models.Title{BoletoType: "OUT"}, Expected: "99"},
}

func TestProcessBoleto_WhenServiceRespondsSuccessfully_ShouldHasSuccessfulBoletoResponse(t *testing.T) {
	mock.StartMockService("9093")
	input := newStubBoletoRequestBradescoShopFacil().Build()
	bank := New()

	output, _ := bank.ProcessBoleto(input)

	test.AssertProcessBoletoWithSuccess(t, output)
}

func TestProcessBoleto_WhenServiceRespondsFailed_ShouldHasFailedBoletoResponse(t *testing.T) {
	mock.StartMockService("9093")
	input := newStubBoletoRequestBradescoShopFacil().WithAmountInCents(400).Build()
	bank := New()

	output, _ := bank.ProcessBoleto(input)

	test.AssertProcessBoletoFailed(t, output)
}

func TestBarcodeGenerationBradescoShopFacil(t *testing.T) {
	const expected = "23795796800000001990001250012446693212345670"
	expireDate, _ := time.Parse("02-01-2006", "01-08-2019")
	boleto := newStubBoletoRequestBradescoShopFacil().WithAgreementAgency("1").WithAgreementAccount("1234567").WithExpirationDate(expireDate).WithAmountInCents(199).WithOurNumber(124466932).Build()

	bc := getBarcode(*boleto)

	assert.Equal(t, expected, bc.toString(), "Deve-se montar o código de barras do BradescoShopFacil")
}

func TestRemoveDigitFromAccount(t *testing.T) {
	const expected = "23791796800000001992372250012446693300056000"

	bc := barcode{
		account:       "0005600",
		bankCode:      "237",
		currencyCode:  "9",
		agency:        "2372",
		dateDueFactor: "7968",
		ourNumber:     "00124466933",
		zero:          "0",
		wallet:        "25",
		value:         "0000000199",
	}

	assert.Equal(t, expected, bc.toString(), "Deve-se montar identificar e remover o digito da conta")
}

func TestGetBoletoType_WhenCalled_ShouldBeMapTypeSuccessful(t *testing.T) {
	BradescoShopFacilStub := newStubBoletoRequestBradescoShopFacil().Build()
	for _, fact := range boletoTypeParameters {
		BradescoShopFacilStub.Title = fact.Input.(models.Title)
		_, result := getBoletoType(BradescoShopFacilStub)
		assert.Equal(t, fact.Expected, result, "Deve mapear o boleto type corretamente")
	}
}
