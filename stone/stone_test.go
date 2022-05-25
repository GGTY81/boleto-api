//go:build integration || !unit
// +build integration !unit

package stone

import (
	"fmt"
	"testing"
	"time"

	"github.com/PMoneda/flow"
	"github.com/mundipagg/boleto-api/mock"
	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/test"
	"github.com/mundipagg/boleto-api/tmpl"
	"github.com/mundipagg/boleto-api/util"
	"github.com/stretchr/testify/assert"
)

var boletoTypeParameters = []test.Parameter{
	{Input: models.Title{BoletoType: ""}, Expected: "bill_of_exchange"},
	{Input: models.Title{BoletoType: "NSA"}, Expected: "bill_of_exchange"},
	{Input: models.Title{BoletoType: "BDP"}, Expected: "bill_of_exchange"},
}

var boletoResponseFailParameters = []test.Parameter{
	{Input: newStubBoletoRequestStone().WithAccessKey("").Build(), Expected: models.ErrorResponse{Code: `MP400`, Message: `o campo AccessKey não pode ser vazio`}},
	{Input: newStubBoletoRequestStone().WithAmountInCents(200).Build(), Expected: models.ErrorResponse{Code: "MPOurNumberFail", Message: "our number was not returned by the bank"}},
	{Input: newStubBoletoRequestStone().WithAmountInCents(401).Build(), Expected: models.ErrorResponse{Code: `srn:error:unauthenticated`, Message: `srn:error:unauthenticated`}},
	{Input: newStubBoletoRequestStone().WithAmountInCents(403).Build(), Expected: models.ErrorResponse{Code: `srn:error:unauthorized`, Message: `srn:error:unauthorized`}},
	{Input: newStubBoletoRequestStone().WithAmountInCents(409).Build(), Expected: models.ErrorResponse{Code: `srn:error:conflict`, Message: `srn:error:conflict`}},
	{Input: newStubBoletoRequestStone().WithAmountInCents(422).Build(), Expected: models.ErrorResponse{Code: `srn:error:product_not_enabled`, Message: `barcode_payment_invoice_bill_of_exchange is not ena bled on this account`}},
	{Input: newStubBoletoRequestStone().WithAmountInCents(4001).Build(), Expected: models.ErrorResponse{Code: `srn:error:validation`, Message: `[{error:is invalid,path:[customer,document]}]`}},
	{Input: newStubBoletoRequestStone().WithAmountInCents(4002).Build(), Expected: models.ErrorResponse{Code: `srn:error:validation`, Message: `[{error:can&#39;t be blank,path:[customer,legal_name]}]`}},
	{Input: newStubBoletoRequestStone().WithAmountInCents(4003).Build(), Expected: models.ErrorResponse{Code: `srn:error:validation`, Message: `[{error:not allowed,path:[amount]}]`}},
	{Input: newStubBoletoRequestStone().WithAmountInCents(4004).Build(), Expected: models.ErrorResponse{Code: `srn:error:validation`, Message: `[{error:is invalid,path:[receiver,document]}]`}},
	{Input: newStubBoletoRequestStone().WithAmountInCents(4005).Build(), Expected: models.ErrorResponse{Code: `srn:error:validation`, Message: `[{error:is invalid,path:[account_id]},{error:not allowed,path:[amount]}]`}},
	{Input: newStubBoletoRequestStone().WithAmountInCents(301).Build(), Expected: models.ErrorResponse{Code: `srn:error:validation`, Message: `[{error:Percentage must be equal or lower than 2.0,path:[fine,value]}]`}},
	{Input: newStubBoletoRequestStone().WithAmountInCents(302).Build(), Expected: models.ErrorResponse{Code: `srn:error:validation`, Message: `[{error:Percentage must be equal or lower than 1.0,path:[interest,value]}]`}},
	{Input: newStubBoletoRequestStone().WithAmountInCents(303).Build(), Expected: models.ErrorResponse{Code: `srn:error:validation`, Message: `[{error:fine date should be greater than expiration date,path:[fine]}]`}},
	{Input: newStubBoletoRequestStone().WithAmountInCents(304).Build(), Expected: models.ErrorResponse{Code: `srn:error:validation`, Message: `[{error:interest date should be greater than expiration date,path:[interest]}]`}},
	{Input: newStubBoletoRequestStone().WithAmountInCents(504).Build(), Expected: models.ErrorResponse{Code: `MPTimeout`, Message: `Post http://localhost:9099/stone/registrarBoleto: context deadline exceeded`}},
}

func Test_GetBoletoType_WhenCalled_ShouldBeMapTypeSuccessful(t *testing.T) {
	request := new(models.BoletoRequest)
	for _, fact := range boletoTypeParameters {
		request.Title = fact.Input.(models.Title)
		_, result := getBoletoType(request)
		assert.Equal(t, fact.Expected, result, "Deve mapear o boleto type corretamente")
	}
}
func Test_TemplateRequestStone_WhenBuyerIsPerson_ParseSuccessful(t *testing.T) {
	var result map[string]interface{}
	f := flow.NewFlow()
	input := newStubBoletoRequestStone().WithBoletoType("DM").Build()

	body := fmt.Sprintf("%v", f.From("message://?source=inline", input, templateRequest, tmpl.GetFuncMaps()).GetBody())
	util.FromJSON(body, &result)

	assert.Equal(t, result["account_id"], input.Authentication.AccessKey)
	assert.Equal(t, uint64(result["amount"].(float64)), input.Title.AmountInCents)
	assert.Equal(t, result["expiration_date"], input.Title.ExpireDate)
	assert.Equal(t, result["invoice_type"], input.Title.BoletoTypeCode)
	assert.Equal(t, result["customer"].(map[string]interface{})["document"], input.Buyer.Document.Number)
	assert.Equal(t, result["customer"].(map[string]interface{})["legal_name"], input.Buyer.Name)
	assert.Equal(t, result["customer"].(map[string]interface{})["trade_name"], nil)
}
func Test_TemplateRequestStone_WhenBuyerIsCompany_ParseSuccessful(t *testing.T) {
	var result map[string]interface{}
	f := flow.NewFlow()
	input := newStubBoletoRequestStone().WithDocument("12123123000112", "CNPJ").WithBoletoType("DM").Build()

	body := fmt.Sprintf("%v", f.From("message://?source=inline", input, templateRequest, tmpl.GetFuncMaps()).GetBody())
	util.FromJSON(body, &result)

	assert.Equal(t, result["account_id"], input.Authentication.AccessKey)
	assert.Equal(t, uint64(result["amount"].(float64)), input.Title.AmountInCents)
	assert.Equal(t, result["expiration_date"], input.Title.ExpireDate)
	assert.Equal(t, result["invoice_type"], input.Title.BoletoTypeCode)
	assert.Equal(t, result["customer"].(map[string]interface{})["document"], input.Buyer.Document.Number)
	assert.Equal(t, result["customer"].(map[string]interface{})["legal_name"], input.Buyer.Name)
	assert.Equal(t, result["customer"].(map[string]interface{})["trade_name"], input.Buyer.Name)
}

func Test_ProcessBoleto_WhenServiceRespondsSuccessfully_ShouldHasSuccessfulBoletoResponse(t *testing.T) {
	mock.StartMockService("9093")

	input := newStubBoletoRequestStone().WithAmountInCents(201).Build()
	bank := New()

	output, _ := bank.ProcessBoleto(input)

	test.AssertProcessBoletoWithSuccess(t, output)
}

func Test_ProcessBoleto_WhenServiceRespondsUnsuccessful_ShouldHasErrorResponse(t *testing.T) {
	bank := New()
	mock.StartMockService("9093")

	for _, fact := range boletoResponseFailParameters {
		request := fact.Input.(*models.BoletoRequest)
		response, _ := bank.ProcessBoleto(request)

		test.AssertProcessBoletoFailed(t, response)
		assert.Equal(t, fact.Expected.(models.ErrorResponse).Code, response.Errors[0].Code)
		assert.Equal(t, fact.Expected.(models.ErrorResponse).Message, response.Errors[0].Message)
	}
}

func Test_GetBankNumber(t *testing.T) {
	bank := New()

	result := bank.GetBankNumber()

	assert.Equal(t, models.Stone, int(result))
}

func Test_GetBankNameIntegration(t *testing.T) {
	bank := New()

	result := bank.GetBankNameIntegration()

	assert.Equal(t, "Stone", result)
}

func Test_GetBankLog(t *testing.T) {
	bank := New()

	result := bank.Log()

	assert.NotNil(t, result)
}

func Test_bankStone_ProcessBoleto(t *testing.T) {
	mock.StartMockService("9093")

	bankInst := New()

	type args struct {
		request *models.BoletoRequest
	}
	tests := []struct {
		name    string
		b       bankStone
		args    args
		want    models.BoletoResponse
		wantErr bool
	}{
		{
			name: "StoneEmptyAccessKeyRequest",
			b:    bankInst,
			args: args{
				request: successRequest,
			},
			want: models.BoletoResponse{
				StatusCode: 0,
				Errors: []models.ErrorResponse{
					{
						Code:    "MP400",
						Message: "o campo AccessKey não pode ser vazio",
					},
				},
				ID:            "",
				DigitableLine: "",
				BarCodeNumber: "",
				OurNumber:     "",
				Links:         []models.Link{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := tt.b.ProcessBoleto(tt.args.request)
			assert.Greater(t, len(got.Errors), 0)
			err := got.Errors[0]
			assert.Equal(t, err.Code, "MP400")
			assert.Equal(t, err.Message, "o campo AccessKey não pode ser vazio")
		})
	}
}

func BenchmarkBankStoneProcessBoleto(b *testing.B) {
	mock.StartMockService("9093")

	input := newStubBoletoRequestStone().WithAmountInCents(201).Build()
	bank := New()

	bank.ProcessBoleto(input)
}

func TestTemplateResponse_WhenRequestHasSpecialCharacter_ShouldBeParsedSuccessful(t *testing.T) {
	mock.StartMockService("9093")
	input := newStubBoletoRequestStone().WithAmountInCents(201).WithBuyerName("Nome do \tComprador (Cliente)").Build()
	bank := New()

	output, _ := bank.ProcessBoleto(input)

	test.AssertProcessBoletoWithSuccess(t, output)
}

func Test_TemplateRequestStone_WhenHasFine_AndWithPercentageOnTotal_ParseSuccessful(t *testing.T) {
	mock.StartMockService("9099")
	var result map[string]interface{}
	f := flow.NewFlow()

	var boletoDate uint = 1
	var amount uint64 = 0
	var percentage float64 = 2.0

	input := newStubBoletoRequestStone().WithFine(boletoDate, amount, percentage).Build()

	valueFineExpected := "2.00"
	daysToAdd := input.Title.Fees.Fine.DaysAfterExpirationDate
	dateFineExpected := input.Title.ExpireDateTime.UTC().Add(day * time.Duration(daysToAdd)).Format("2006-01-02")

	body := fmt.Sprintf("%v", f.From("message://?source=inline", input, templateRequest, tmpl.GetFuncMaps()).GetBody())
	_ = util.FromJSON(body, &result)

	assert.Equal(t, result["fine"].(map[string]interface{})["date"], dateFineExpected)
	assert.Equal(t, result["fine"].(map[string]interface{})["value"], valueFineExpected)
}

func Test_TemplateRequestStone_WhenHasFine_AndWithAmountInCents_ParseSuccessful(t *testing.T) {
	mock.StartMockService("9099")
	var result map[string]interface{}
	f := flow.NewFlow()

	var boletoDate uint = 1
	var amount uint64 = 25
	var percentage float64 = 0
	input := newStubBoletoRequestStone().WithAmountInCents(20000).WithFine(boletoDate, amount, percentage).Build()

	valueFineExpected := "0.12"
	daysToAdd := input.Title.Fees.Fine.DaysAfterExpirationDate
	dateFineExpected := input.Title.ExpireDateTime.UTC().Add(day * time.Duration(daysToAdd)).Format("2006-01-02")

	body := fmt.Sprintf("%v", f.From("message://?source=inline", input, templateRequest, tmpl.GetFuncMaps()).GetBody())
	_ = util.FromJSON(body, &result)

	assert.Equal(t, result["fine"].(map[string]interface{})["date"], dateFineExpected)
	assert.Equal(t, result["fine"].(map[string]interface{})["value"], valueFineExpected)
}

func Test_TemplateRequestStone_WhenHasInterest_AndWithPercentageOnTotal_ParseSuccessful(t *testing.T) {
	mock.StartMockService("9099")
	var result map[string]interface{}
	f := flow.NewFlow()

	var boletoDate uint = 1
	var amount uint64 = 0
	var percentage float64 = 1.0
	input := newStubBoletoRequestStone().WithInterest(boletoDate, amount, percentage).Build()

	valueFineExpected := "1.00"
	daysToAdd := input.Title.Fees.Interest.DaysAfterExpirationDate
	dateFineExpected := input.Title.ExpireDateTime.UTC().Add(day * time.Duration(daysToAdd)).Format("2006-01-02")

	body := fmt.Sprintf("%v", f.From("message://?source=inline", input, templateRequest, tmpl.GetFuncMaps()).GetBody())
	_ = util.FromJSON(body, &result)

	assert.Equal(t, result["interest"].(map[string]interface{})["date"], dateFineExpected)
	assert.Equal(t, result["interest"].(map[string]interface{})["value"], valueFineExpected)
}

func Test_TemplateRequestStone_WhenHasInterest_AndWithAmountInCents_ParseSuccessful(t *testing.T) {
	mock.StartMockService("9099")
	var result map[string]interface{}
	f := flow.NewFlow()

	var boletoDate uint = 1
	var amount uint64 = 5
	var percentage float64 = 0
	input := newStubBoletoRequestStone().WithAmountInCents(20000).WithInterest(boletoDate, amount, percentage).Build()

	valueInterestxpected := "0.75"
	daysToAdd := input.Title.Fees.Interest.DaysAfterExpirationDate
	dateInterestExpected := input.Title.ExpireDateTime.UTC().Add(day * time.Duration(daysToAdd)).Format("2006-01-02")

	body := fmt.Sprintf("%v", f.From("message://?source=inline", input, templateRequest, tmpl.GetFuncMaps()).GetBody())
	_ = util.FromJSON(body, &result)

	assert.Equal(t, result["interest"].(map[string]interface{})["date"], dateInterestExpected)
	assert.Equal(t, result["interest"].(map[string]interface{})["value"], valueInterestxpected)
}
