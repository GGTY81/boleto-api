package caixa

import (
	"fmt"
	"testing"
	"time"

	"github.com/PMoneda/flow"
	"github.com/mundipagg/boleto-api/mock"
	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/test"
	"github.com/mundipagg/boleto-api/tmpl"
	"github.com/stretchr/testify/assert"
)

var boletoTypeParameters = []test.Parameter{
	{Input: models.Title{BoletoType: ""}, Expected: "99"},
	{Input: models.Title{BoletoType: "NSA"}, Expected: "99"},
	{Input: models.Title{BoletoType: "BDP"}, Expected: "99"},
}

var boletoBuyerNameParameters = []test.Parameter{
	{Input: "Leonardo Jasmim", Expected: "<NOME>Leonardo Jasmim</NOME>"},
	{Input: "Ântôníõ Tùpìnâmbáú", Expected: "<NOME>Antonio Tupinambau</NOME>"},
	{Input: "Accepted , / ( ) * = - + ! : ? . ; _ ' ", Expected: "<NOME>Accepted , / ( ) * = - &#43; ! : ? . ; _ &#39; </NOME>"},
	{Input: "NotAccepted @#$%¨{}[]^~\"&<>\\", Expected: "<NOME>NotAccepted                 </NOME>"},
}

var boletoInstructionsParameters = []test.Parameter{
	{Input: ", / ( ) * = - + ! : ? . ; _ ' ", Expected: "<MENSAGEM>, / ( ) * = - &#43; ! : ? . ; _ &#39; </MENSAGEM>"},
	{Input: "@ # $ % ¨ { } [ ] ^ ~ \" & < > \\", Expected: "                              "},
}

var boletoCepParameters = []test.Parameter{
	{Input: "12345-123", Expected: "<CEP>12345123</CEP>"},
	{Input: "01000-000", Expected: "<CEP>01000000</CEP>"},
	{Input: "08000-123", Expected: "<CEP>08000123</CEP>"},
	{Input: "5000-098", Expected: "<CEP>5000098</CEP>"},
	{Input: "98765-456", Expected: "<CEP>98765456</CEP>"},
}

var interestPercentageParameters = []test.Parameter{
	{Input: models.Interest{DaysAfterExpirationDate: 1, AmountPerDayInCents: 0, PercentagePerMonth: 3.27}},
	{Input: models.Interest{DaysAfterExpirationDate: 1, AmountPerDayInCents: 0, PercentagePerMonth: 4.238}},
}

func TestProcessBoleto_WhenServiceRespondsSuccessfully_ShouldHasSuccessfulBoletoResponse(t *testing.T) {
	mock.StartMockService("9094")

	input := newStubBoletoRequestCaixa().Build()
	bank := New()

	output, _ := bank.ProcessBoleto(input)

	test.AssertProcessBoletoWithSuccess(t, output)
}

func TestProcessBoleto_WhenServiceRespondsFailed_ShouldHasFailedBoletoResponse(t *testing.T) {
	mock.StartMockService("9094")

	input := newStubBoletoRequestCaixa().WithAmountInCents(400).Build()
	bank := New()

	output, _ := bank.ProcessBoleto(input)

	test.AssertProcessBoletoFailed(t, output)
}

func TestProcessBoleto_WhenRequestContainsInvalidOurNumberParameter_ShouldHasFailedBoletoResponse(t *testing.T) {
	largeOurNumber := uint(9999999999999999)
	mock.StartMockService("9094")

	input := newStubBoletoRequestCaixa().WithOurNumber(largeOurNumber).Build()

	bank := New()

	output, _ := bank.ProcessBoleto(input)

	test.AssertProcessBoletoFailed(t, output)
}

func TestGetCaixaCheckSumInfo(t *testing.T) {
	const expectedSumCode = "0200656000000000000000003008201700000000000100000732159000109"
	const expectedToken = "LvWr1op5Ayibn6jsCQ3/2bW4KwThVAlLK5ftxABlq20="

	bank := New()
	agreement := uint(200656)
	expiredAt := time.Date(2017, 8, 30, 12, 12, 12, 12, time.Local)
	doc := "00732159000109"

	s := newStubBoletoRequestCaixa()
	s.WithAgreementNumber(agreement).WithOurNumber(0).WithAmountInCents(1000)
	s.WithExpirationDate(expiredAt).WithRecipientDocumentNumber(doc)

	input := s.Build()

	assert.Equal(t, expectedSumCode, bank.getCheckSumCode(*input), "Deve-se formar uma string seguindo o padrão da documentação")
	assert.Equal(t, expectedToken, bank.getAuthToken(bank.getCheckSumCode(*input)), "Deve-se fazer um hash sha256 e encodar com base64")
}

func TestShouldCalculateAccountDigitCaixa(t *testing.T) {
	input := newStubBoletoRequestCaixa().WithAgreementAccount("100000448").WithAgreementAgency("2004").Build()

	assert.Nil(t, caixaValidateAccountAndDigit(input))
	assert.Nil(t, caixaValidateAgency(input))
}

func TestGetBoletoType_WhenCalled_ShouldBeMapTypeSuccessful(t *testing.T) {
	request := new(models.BoletoRequest)
	for _, fact := range boletoTypeParameters {
		request.Title = fact.Input.(models.Title)
		_, result := getBoletoType(request)
		assert.Equal(t, fact.Expected, result, "Deve mapear o boleto type corretamente")
	}
}

func TestTemplateRequestCaixa_CEP_ParseSuccessful(t *testing.T) {
	f := flow.NewFlow()
	s := newStubBoletoRequestCaixa()

	for _, fact := range boletoCepParameters {
		request := s.WithBuyerZipCode(fact.Input.(string)).Build()
		result := fmt.Sprintf("%v", f.From("message://?source=inline", request, getRequestCaixa(), tmpl.GetFuncMaps()).GetBody())
		assert.Contains(t, result, fact.Expected, "Conversão não realizada como esperado")
	}
}

func TestTemplateRequestCaixa_BuyerName_ParseSuccessful(t *testing.T) {
	f := flow.NewFlow()
	s := newStubBoletoRequestCaixa()

	for _, fact := range boletoBuyerNameParameters {
		request := s.WithBuyerName(fact.Input.(string)).Build()
		result := fmt.Sprintf("%v", f.From("message://?source=inline", request, getRequestCaixa(), tmpl.GetFuncMaps()).GetBody())
		assert.Contains(t, result, fact.Expected, "Conversão não realizada como esperado")
	}
}

func TestTemplateRequestCaixa_Instructions_ParseSuccessful(t *testing.T) {
	f := flow.NewFlow()
	s := newStubBoletoRequestCaixa()

	for _, fact := range boletoInstructionsParameters {
		request := s.WithInstructions(fact.Input.(string)).Build()
		result := fmt.Sprintf("%v", f.From("message://?source=inline", request, getRequestCaixa(), tmpl.GetFuncMaps()).GetBody())
		assert.Contains(t, result, fact.Expected, "Conversão não realizada como esperado")
	}
}

func TestTemplateRequestCaixa_WhenRequestV1_ParseSuccessful(t *testing.T) {
	f := flow.NewFlow()
	input := newStubBoletoRequestCaixa().WithBuyerAddress().Build()

	b := fmt.Sprintf("%v", f.From("message://?source=inline", input, getRequestCaixa(), tmpl.GetFuncMaps()).GetBody())

	for _, expected := range expectedBasicTitleRequestFields {
		assert.Contains(t, b, expected, "Erro no mapeamento dos campos básicos do Título")
	}

	for _, expected := range expectedBuyerRequestFields {
		assert.Contains(t, b, expected, "Erro no mapeamento dos campos básicos do Comprador")
	}

	for _, notExpected := range expectedStrictRulesFieldsV2 {
		assert.NotContains(t, b, notExpected, "Não devem haver campos de regras de pagamento na V1")
	}

	for _, notExpected := range expectedFlexRulesFieldsV2 {
		assert.NotContains(t, b, notExpected, "Não devem haver campos de regras de pagamento na V1")
	}
}

func TestTemplateRequestCaixa_WhenRequestWithStrictRulesV2_ParseSuccessful(t *testing.T) {
	f := flow.NewFlow()
	input := newStubBoletoRequestCaixa().WithBuyerAddress().WithStrictRules().Build()

	b := fmt.Sprintf("%v", f.From("message://?source=inline", input, getRequestCaixa(), tmpl.GetFuncMaps()).GetBody())

	for _, expected := range expectedBasicTitleRequestFields {
		assert.Contains(t, b, expected, "Erro no mapeamento dos campos básicos do Título")
	}

	for _, expected := range expectedBuyerRequestFields {
		assert.Contains(t, b, expected, "Erro no mapeamento dos campos básicos do Comprador")
	}

	for _, expected := range expectedStrictRulesFieldsV2 {
		assert.Contains(t, b, expected, "Erro no mapeamento das regras de pagamento")
	}
}

func TestTemplateRequestCaixa_WhenRequestWithFlexRulesV2_ParseSuccessful(t *testing.T) {
	f := flow.NewFlow()
	input := newStubBoletoRequestCaixa().WithBuyerAddress().WithFlexRules().Build()

	b := fmt.Sprintf("%v", f.From("message://?source=inline", input, getRequestCaixa(), tmpl.GetFuncMaps()).GetBody())

	for _, expected := range expectedBasicTitleRequestFields {
		assert.Contains(t, b, expected, "Erro no mapeamento dos campos básicos do Título")
	}

	for _, expected := range expectedBuyerRequestFields {
		assert.Contains(t, b, expected, "Erro no mapeamento dos campos básicos do Comprador")
	}

	for _, expected := range expectedFlexRulesFieldsV2 {
		assert.Contains(t, b, expected, "Erro no mapeamento das regras de pagamento")
	}
}

func TestTemplateRequestCaixa_NumberOfDaysAfterExpirationEqualsOne_NumberOfDaysAfterExpirationIsRight(t *testing.T) {
	flow := flow.NewFlow()
	caixaRequestStub := newStubBoletoRequestCaixa()

	request := caixaRequestStub.WithStrictRules().Build()
	result := fmt.Sprintf("%v", flow.From("message://?source=inline", request, getRequestCaixa(), tmpl.GetFuncMaps()).GetBody())
	assert.Contains(t, result, "<NUMERO_DIAS>1</NUMERO_DIAS>", "Falha ao encontrar o campo <NUMERO_DIAS>1<NUMERO_DIAS> no request")
}

func TestTemplateRequestCaixa_WhenRequestWithFineAndAmountInCents_ParseSuccessful(t *testing.T) {
	f := flow.NewFlow()
	input := newStubBoletoRequestCaixa().WithFine(1, 200, 0).Build()
	dateExpected := input.Title.ExpireDateTime.Format("2006-01-02")

	b := fmt.Sprintf("%v", f.From("message://?source=inline", input, getRequestCaixa(), tmpl.GetFuncMaps()).GetBody())

	assert.Contains(t, b, "<MULTA>", "Erro no mapeamento dos campos básicos do multa - <MULTA>")
	assert.Contains(t, b, fmt.Sprintf("<DATA>%s</DATA>", dateExpected), "Erro no mapeamento dos campos básicos do multa - <DATA>")
	assert.Contains(t, b, "<VALOR>2.00</VALOR>", "Erro no mapeamento dos campos básicos do multa - <VALOR>")
	assert.Contains(t, b, "</MULTA>", "Erro no mapeamento dos campos básicos do multa - </MULTA>")
}

func TestTemplateRequestCaixa_WhenRequestWithFineAndPercentageOnTotal_ParseSuccessful(t *testing.T) {
	f := flow.NewFlow()
	input := newStubBoletoRequestCaixa().WithFine(1, 0, 1.3).Build()
	dateExpected := input.Title.ExpireDateTime.Format("2006-01-02")

	b := fmt.Sprintf("%v", f.From("message://?source=inline", input, getRequestCaixa(), tmpl.GetFuncMaps()).GetBody())

	assert.Contains(t, b, "<MULTA>", "Erro no mapeamento dos campos básicos do multa - <MULTA>")
	assert.Contains(t, b, fmt.Sprintf("<DATA>%s</DATA>", dateExpected), "Erro no mapeamento dos campos básicos do multa - <DATA>")
	assert.Contains(t, b, "<PERCENTUAL>1.30</PERCENTUAL>", "Erro no mapeamento dos campos básicos do multa - <VALOR>")
	assert.Contains(t, b, "</MULTA>", "Erro no mapeamento dos campos básicos do multa - </MULTA>")
}

func TestTemplateRequestCaixa_WhenRequestWithInterestAndTypeIsento_ParseSuccessful(t *testing.T) {
	f := flow.NewFlow()
	input := newStubBoletoRequestCaixa().Build()

	b := fmt.Sprintf("%v", f.From("message://?source=inline", input, getRequestCaixa(), tmpl.GetFuncMaps()).GetBody())

	assert.Contains(t, b, "<TIPO>ISENTO</TIPO>", "Erro no mapeamento dos campos básicos de juros - <TIPO>")
	assert.Contains(t, b, "<VALOR>0</VALOR>", "Erro no mapeamento dos campos básicos de juros - <VALOR>")
}

func TestTemplateRequestCaixa_WhenRequestWithInterestAndAmountPerDayInCents_ParseSuccessful(t *testing.T) {
	f := flow.NewFlow()
	input := newStubBoletoRequestCaixa().WithInterest(1, 300, 0).Build()
	daysToAdd := input.Title.Fees.Interest.DaysAfterExpirationDate
	dateExpected := input.Title.ExpireDateTime.UTC().Add(day * time.Duration(daysToAdd)).Format("2006-01-02")

	b := fmt.Sprintf("%v", f.From("message://?source=inline", input, getRequestCaixa(), tmpl.GetFuncMaps()).GetBody())

	assert.Contains(t, b, "<TIPO>VALOR_POR_DIA</TIPO>", "Erro no mapeamento dos campos básicos de juros - <TIPO>")
	assert.Contains(t, b, fmt.Sprintf("<DATA>%s</DATA>", dateExpected), "Erro no mapeamento dos campos básicos de juros - <DATA>")
	assert.Contains(t, b, "<VALOR>3.00</VALOR>", "Erro no mapeamento dos campos básicos de juros - <VALOR>")
}

func TestTemplateRequestCaixa_WhenRequestWithInterestAndPercentagePerMonth_ParseSuccessful(t *testing.T) {
	f := flow.NewFlow()
	stub := newStubBoletoRequestCaixa()

	for _, fact := range interestPercentageParameters {
		interest := fact.Input.(models.Interest)
		request := stub.WithInterest(interest.DaysAfterExpirationDate, interest.AmountPerDayInCents, interest.PercentagePerMonth).Build()
		daysToAdd := request.Title.Fees.Interest.DaysAfterExpirationDate
		dateExpected := request.Title.ExpireDateTime.UTC().Add(day * time.Duration(daysToAdd)).Format("2006-01-02")
		percentualExpected := fmt.Sprintf("%.2f", request.Title.Fees.Interest.PercentagePerMonth)

		b := fmt.Sprintf("%v", f.From("message://?source=inline", request, getRequestCaixa(), tmpl.GetFuncMaps()).GetBody())

		assert.Contains(t, b, "<TIPO>TAXA_MENSAL</TIPO>", "Erro no mapeamento dos campos básicos de juros - <TIPO>")
		assert.Contains(t, b, fmt.Sprintf("<DATA>%s</DATA>", dateExpected), "Erro no mapeamento dos campos básicos de juros - <DATA>")
		assert.Contains(t, b, fmt.Sprintf("<PERCENTUAL>%s</PERCENTUAL>", percentualExpected), "Erro no mapeamento dos campos básicos de juros - <VALOR>")
	}
}

func TestCaixaValidateFine(t *testing.T) {
	input := newStubBoletoRequestCaixa().WithFine(1, 200, 0).Build()

	assert.Nil(t, caixaValidateFine(input))
}

func TestCaixaValidateFineWithNilFees(t *testing.T) {
	input := newStubBoletoRequestCaixa().Build()

	assert.Nil(t, caixaValidateFine(input))
}

func TestCaixaValidateInterest(t *testing.T) {
	input := newStubBoletoRequestCaixa().WithInterest(1, 0, 10.2).Build()

	assert.Nil(t, caixaValidateInterest(input))
}

func TestCaixaValidateInterestWithNilFees(t *testing.T) {
	input := newStubBoletoRequestCaixa().Build()

	assert.Nil(t, caixaValidateInterest(input))
}

func TestTemplateRequestCaixa_WhenRequestWithPayeeGuarantorIsCNPJ_ParseSuccessful(t *testing.T) {
	f := flow.NewFlow()
	input := newStubBoletoRequestCaixa().WithPayeeGuarantorName("PayeeGuarantor Test Name").WithPayeeGuarantorDocumentType("CNPJ").Build()

	b := fmt.Sprintf("%v", f.From("message://?source=inline", input, getRequestCaixa(), tmpl.GetFuncMaps()).GetBody())
	nodeContent := test.GetNode(b, "SACADOR_AVALISTA")

	assert.Contains(t, nodeContent, "<RAZAO_SOCIAL>PayeeGuarantor Test Name</RAZAO_SOCIAL>", "Erro no mapeamento do campo name do PayeeGuarantor")
	assert.Contains(t, nodeContent, fmt.Sprintf("<CNPJ>%s</CNPJ>", input.PayeeGuarantor.Document.Number), "Erro no mapeamento do document type do PayeeGuarantor")
}

func TestTemplateRequestCaixa_WhenRequestWithPayeeGuarantorIsCPF_ParseSuccessful(t *testing.T) {
	f := flow.NewFlow()
	input := newStubBoletoRequestCaixa().WithPayeeGuarantorName("PayeeGuarantor Test Name").WithPayeeGuarantorDocumentType("CPF").Build()

	b := fmt.Sprintf("%v", f.From("message://?source=inline", input, getRequestCaixa(), tmpl.GetFuncMaps()).GetBody())
	nodeContent := test.GetNode(b, "SACADOR_AVALISTA")

	assert.Contains(t, nodeContent, "<NOME>PayeeGuarantor Test Name</NOME>", "Erro no mapeamento do campo name do PayeeGuarantor")
	assert.Contains(t, nodeContent, fmt.Sprintf("<CPF>%s</CPF>", input.PayeeGuarantor.Document.Number), "Erro no mapeamento do document type do PayeeGuarantor")
}

func TestTemplateRequestCaixa_WhenRequestWithPayeeGuarantorIsCPF_ParseFailed(t *testing.T) {
	f := flow.NewFlow()
	input := newStubBoletoRequestCaixa().WithPayeeGuarantorName("PayeeGuarantor Test Name").WithPayeeGuarantorDocumentType("CNPJ").Build()

	b := fmt.Sprintf("%v", f.From("message://?source=inline", input, getRequestCaixa(), tmpl.GetFuncMaps()).GetBody())
	nodeContent := test.GetNode(b, "SACADOR_AVALISTA")

	assert.NotContains(t, nodeContent, "<NOME>PayeeGuarantor Test Name</NOME>", "Erro no mapeamento do campo PayeeGuarantor")
	assert.NotContains(t, nodeContent, fmt.Sprintf("<CPF>%s</CPF>", input.PayeeGuarantor.Document.Number), "Erro no mapeamento do document type do PayeeGuarantor")
}

func TestTemplateRequestCaixa_WhenRequestWithPayeeGuarantorIsCNPJ_ParseFailed(t *testing.T) {
	f := flow.NewFlow()
	input := newStubBoletoRequestCaixa().WithPayeeGuarantorName("PayeeGuarantor Test Name").WithPayeeGuarantorDocumentType("CPF").Build()

	b := fmt.Sprintf("%v", f.From("message://?source=inline", input, getRequestCaixa(), tmpl.GetFuncMaps()).GetBody())
	nodeContent := test.GetNode(b, "SACADOR_AVALISTA")

	assert.NotContains(t, nodeContent, "<RAZAO_SOCIAL>PayeeGuarantor Test Name</RAZAO_SOCIAL>", "Erro no mapeamento do campo PayeeGuarantor")
	assert.NotContains(t, nodeContent, fmt.Sprintf("<CNPJ>%s</CNPJ>", input.PayeeGuarantor.Document.Number), "Erro no mapeamento do document type do PayeeGuarantor")
}

func TestTemplateRequestCaixa_WhenRequestWithPayeeGuarantorDocumentInNotValidCPF_ParseFailed(t *testing.T) {
	mock.StartMockService("9094")
	cnpjDocument := "00732159000109"

	input := newStubBoletoRequestCaixa().WithPayeeGuarantorName("PayeeGuarantor Test Name").WithPayeeGuarantorDocumentNumber(cnpjDocument).WithPayeeGuarantorDocumentType("CPF").Build()

	bank := New()

	output, _ := bank.ProcessBoleto(input)

	test.AssertProcessBoletoFailed(t, output)
}

func TestTemplateRequestCaixa_WhenRequestWithPayeeGuarantorDocumentInNotValidCNPJ_ParseFailed(t *testing.T) {
	mock.StartMockService("9094")
	cpfDocument := "08013156036"

	input := newStubBoletoRequestCaixa().WithPayeeGuarantorName("PayeeGuarantor Test Name").WithPayeeGuarantorDocumentNumber(cpfDocument).WithPayeeGuarantorDocumentType("CNPJ").Build()

	bank := New()

	output, _ := bank.ProcessBoleto(input)

	test.AssertProcessBoletoFailed(t, output)
}

func TestTemplateRequestCaixa_WhenRequestWithoutPayeeGuarantor_HasNotPayeeGuarantorNode(t *testing.T) {
	f := flow.NewFlow()
	input := newStubBoletoRequestCaixa().Build()

	b := fmt.Sprintf("%v", f.From("message://?source=inline", input, getRequestCaixa(), tmpl.GetFuncMaps()).GetBody())
	nodeContent := test.GetNode(b, "SACADOR_AVALISTA")

	assert.Contains(t, nodeContent, "", "Não deve haver o nó PayeeGuarantor")
}

func TestTemplateRequestCaixa_WhenRequestWithoutBuyerAddress_HasNotBuyerAddressNode(t *testing.T) {
	f := flow.NewFlow()
	input := newStubBoletoRequestCaixa().Build()

	b := fmt.Sprintf("%v", f.From("message://?source=inline", input, getRequestCaixa(), tmpl.GetFuncMaps()).GetBody())
	nodeContent := test.GetNode(b, "ENDERECO")

	assert.Empty(t, nodeContent, "Não deve haver o nó Address no Buyer")
}

func TestTemplateRequestCaixa_WhenRequestWithBuyerAddress_HasBuyerAddressNode(t *testing.T) {
	f := flow.NewFlow()
	input := newStubBoletoRequestCaixa().WithBuyerAddress().Build()

	b := fmt.Sprintf("%v", f.From("message://?source=inline", input, getRequestCaixa(), tmpl.GetFuncMaps()).GetBody())
	nodeContent := test.GetNode(b, "ENDERECO")

	assert.NotEmpty(t, nodeContent, "Deve haver o nó Address no Buyer")
}
