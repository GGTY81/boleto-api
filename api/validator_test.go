package api

import (
	"bytes"
	"net/http"
	"testing"
	"time"

	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/test"
	"github.com/mundipagg/boleto-api/util"
	"github.com/stretchr/testify/assert"
)

type bankNumberParameter struct {
	input    models.BankNumber
	expected interface{}
}

type bankNumberExpectedParameter struct {
	code    int
	message string
}

var validationBanksAcceptedRulesAndFees = []bankNumberParameter{
	{input: models.Caixa, expected: true},
	{input: models.Stone, expected: true},
	{input: models.BancoDoBrasil, expected: false},
	{input: models.Santander, expected: false},
	{input: models.Citibank, expected: false},
	{input: models.Itau, expected: false},
	{input: models.JPMorgan, expected: false},
	{input: models.Pefisa, expected: false},
	{input: models.Bradesco, expected: false},
}

var validationRulesAndFeesSuccessParametersV2 = []bankNumberParameter{
	{input: models.Caixa, expected: bankNumberExpectedParameter{code: 200}},
	{input: models.Stone, expected: bankNumberExpectedParameter{code: 200}},
}

var validationFeesFailedParametersV2 = []bankNumberParameter{
	{input: models.BancoDoBrasil, expected: bankNumberExpectedParameter{code: 400, message: `{"errors":[{"code":"MP400","message":"title.fees not available for this bank"}]}`}},
	{input: models.Santander, expected: bankNumberExpectedParameter{code: 400, message: `{"errors":[{"code":"MP400","message":"title.fees not available for this bank"}]}`}},
	{input: models.Citibank, expected: bankNumberExpectedParameter{code: 400, message: `{"errors":[{"code":"MP400","message":"title.fees not available for this bank"}]}`}},
	{input: models.Itau, expected: bankNumberExpectedParameter{code: 400, message: `{"errors":[{"code":"MP400","message":"title.fees not available for this bank"}]}`}},
	{input: models.JPMorgan, expected: bankNumberExpectedParameter{code: 400, message: `{"errors":[{"code":"MP400","message":"title.fees not available for this bank"}]}`}},
	{input: models.Pefisa, expected: bankNumberExpectedParameter{code: 400, message: `{"errors":[{"code":"MP400","message":"title.fees not available for this bank"}]}`}},
	{input: models.Bradesco, expected: bankNumberExpectedParameter{code: 400, message: `{"errors":[{"code":"MP400","message":"title.fees not available for this bank"}]}`}},
}

var validationRulesFailedParametersV2 = []bankNumberParameter{
	{input: models.BancoDoBrasil, expected: bankNumberExpectedParameter{code: 400, message: `{"errors":[{"code":"MP400","message":"title.rules not available for this bank"}]}`}},
	{input: models.Santander, expected: bankNumberExpectedParameter{code: 400, message: `{"errors":[{"code":"MP400","message":"title.rules not available for this bank"}]}`}},
	{input: models.Citibank, expected: bankNumberExpectedParameter{code: 400, message: `{"errors":[{"code":"MP400","message":"title.rules not available for this bank"}]}`}},
	{input: models.Itau, expected: bankNumberExpectedParameter{code: 400, message: `{"errors":[{"code":"MP400","message":"title.rules not available for this bank"}]}`}},
	{input: models.JPMorgan, expected: bankNumberExpectedParameter{code: 400, message: `{"errors":[{"code":"MP400","message":"title.rules not available for this bank"}]}`}},
	{input: models.Pefisa, expected: bankNumberExpectedParameter{code: 400, message: `{"errors":[{"code":"MP400","message":"title.rules not available for this bank"}]}`}},
}

var validationSuccessParametersV1 = []bankNumberParameter{
	{input: models.BancoDoBrasil, expected: bankNumberExpectedParameter{code: 200}},
}

var validationRulesFailedParametersV1 = []bankNumberParameter{
	{input: models.Caixa, expected: bankNumberExpectedParameter{code: 400, message: `{"errors":[{"code":"MP400","message":"title.rules not available in this version"}]}`}},
	{input: models.BancoDoBrasil, expected: bankNumberExpectedParameter{code: 400, message: `{"errors":[{"code":"MP400","message":"title.rules not available in this version"}]}`}},
	{input: models.Bradesco, expected: bankNumberExpectedParameter{code: 400, message: `{"errors":[{"code":"MP400","message":"title.rules not available in this version"}]}`}},
	{input: models.Citibank, expected: bankNumberExpectedParameter{code: 400, message: `{"errors":[{"code":"MP400","message":"title.rules not available in this version"}]}`}},
	{input: models.Itau, expected: bankNumberExpectedParameter{code: 400, message: `{"errors":[{"code":"MP400","message":"title.rules not available in this version"}]}`}},
	{input: models.JPMorgan, expected: bankNumberExpectedParameter{code: 400, message: `{"errors":[{"code":"MP400","message":"title.rules not available in this version"}]}`}},
	{input: models.Pefisa, expected: bankNumberExpectedParameter{code: 400, message: `{"errors":[{"code":"MP400","message":"title.rules not available in this version"}]}`}},
	{input: models.Santander, expected: bankNumberExpectedParameter{code: 400, message: `{"errors":[{"code":"MP400","message":"title.rules not available in this version"}]}`}},
}

func Test_ValidateRegisterV1_WhenWithoutRules_PassSuccessful(t *testing.T) {
	for _, fact := range validationSuccessParametersV1 {
		router, w := arrangeMiddlewareRoute("/validateV1", parseBoleto, validateRegisterV1)
		body := test.NewStubBoletoRequest(fact.input).WithExpirationDate(time.Now()).WithWallet(25).Build()
		req, _ := http.NewRequest("POST", "/validateV1", bytes.NewBuffer([]byte(util.ToJSON(body))))

		router.ServeHTTP(w, req)

		assert.Equal(t, fact.expected.(bankNumberExpectedParameter).code, w.Code)
	}
}

func Test_ValidateRegisterV1_WhenHasRules_ReturnBadRequest(t *testing.T) {
	for _, fact := range validationRulesFailedParametersV1 {
		router, w := arrangeMiddlewareRoute("/validateV1", parseBoleto, validateRegisterV1)
		body := test.NewStubBoletoRequest(fact.input).WithExpirationDate(time.Now()).WithWallet(25).WithAcceptDivergentAmount(true).Build()
		req, _ := http.NewRequest("POST", "/validateV1", bytes.NewBuffer([]byte(util.ToJSON(body))))

		router.ServeHTTP(w, req)

		assert.Equal(t, fact.expected.(bankNumberExpectedParameter).code, w.Code)
		assert.Equal(t, fact.expected.(bankNumberExpectedParameter).message, w.Body.String())
	}
}

func Test_ValidateRegisterV1_WhenHasBankStone_ReturnBadRequest(t *testing.T) {
	router, w := arrangeMiddlewareRoute("/validateV1", parseBoleto, validateRegisterV1)
	body := test.NewStubBoletoRequest(models.Stone).WithExpirationDate(time.Now()).Build()
	req, _ := http.NewRequest("POST", "/validateV1", bytes.NewBuffer([]byte(util.ToJSON(body))))

	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	assert.Equal(t, `{"errors":[{"code":"MP400","message":"bank Stone not available in this version"}]}`, w.Body.String())
}

func Test_ValidateRegisterV2_WhenWithoutRules_PassSuccessful(t *testing.T) {
	router, w := arrangeMiddlewareRoute("/validateV2", parseBoleto, validateRegisterV2)
	body := test.NewStubBoletoRequest(models.BancoDoBrasil).WithExpirationDate(time.Now()).Build()
	req, _ := http.NewRequest("POST", "/validateV2", bytes.NewBuffer([]byte(util.ToJSON(body))))

	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func Test_ValidateRegisterV2_WhenHasRulesAndNotAcceptedBanks_ReturnBadRequest(t *testing.T) {
	for _, fact := range validationRulesFailedParametersV2 {
		router, w := arrangeMiddlewareRoute("/validateV2", parseBoleto, validateRegisterV2)
		body := test.NewStubBoletoRequest(fact.input).WithExpirationDate(time.Now()).WithAcceptDivergentAmount(true).Build()
		req, _ := http.NewRequest("POST", "/validateV2", bytes.NewBuffer([]byte(util.ToJSON(body))))

		router.ServeHTTP(w, req)

		assert.Equal(t, fact.expected.(bankNumberExpectedParameter).code, w.Code)
		assert.Equal(t, fact.expected.(bankNumberExpectedParameter).message, w.Body.String())
	}
}

func Test_ValidateRegisterV2_WhenHasRules_PassSuccessful(t *testing.T) {
	for _, fact := range validationRulesAndFeesSuccessParametersV2 {
		router, w := arrangeMiddlewareRoute("/validateV2", parseBoleto, validateRegisterV2)
		body := test.NewStubBoletoRequest(fact.input).WithExpirationDate(time.Now()).WithMaxDaysToPayPastDue(30).Build()
		req, _ := http.NewRequest("POST", "/validateV2", bytes.NewBuffer([]byte(util.ToJSON(body))))

		router.ServeHTTP(w, req)

		assert.Equal(t, fact.expected.(bankNumberExpectedParameter).code, w.Code)
	}
}

func Test_ValidateRegisterV2_WhenHasFees_PassSuccessful(t *testing.T) {
	var boletoDate uint = 1
	var amount uint64 = 25
	var percentage float64 = 0
	for _, fact := range validationRulesAndFeesSuccessParametersV2 {
		router, w := arrangeMiddlewareRoute("/validateV2", parseBoleto, validateRegisterV2)
		body := test.NewStubBoletoRequest(fact.input).WithExpirationDate(time.Now()).WithFine(boletoDate, amount, percentage).WithMaxDaysToPayPastDue(30).Build()
		req, _ := http.NewRequest("POST", "/validateV2", bytes.NewBuffer([]byte(util.ToJSON(body))))

		router.ServeHTTP(w, req)

		assert.Equal(t, fact.expected.(bankNumberExpectedParameter).code, w.Code)
	}
}

func Test_ValidateRegisterV2_WhenHasFeesAndNotAcceptedBanks_ReturnBadRequest(t *testing.T) {
	var boletoDate uint = 1
	var amount uint64 = 25
	var percentage float64 = 0
	for _, fact := range validationFeesFailedParametersV2 {
		router, w := arrangeMiddlewareRoute("/validateV2", parseBoleto, validateRegisterV2)
		body := test.NewStubBoletoRequest(fact.input).WithExpirationDate(time.Now()).WithFine(boletoDate, amount, percentage).Build()
		req, _ := http.NewRequest("POST", "/validateV2", bytes.NewBuffer([]byte(util.ToJSON(body))))

		router.ServeHTTP(w, req)

		assert.Equal(t, fact.expected.(bankNumberExpectedParameter).code, w.Code)
	}
}

func Test_Banks_Accepted_Rules(t *testing.T) {
	for _, fact := range validationBanksAcceptedRulesAndFees {
		result := isBankNumberAcceptRules(fact.input)
		assert.Equal(t, fact.expected, result)
	}
}

func Test_Banks_Accepted_Fees(t *testing.T) {
	for _, fact := range validationBanksAcceptedRulesAndFees {
		result := isBankNumberAcceptFees(fact.input)
		assert.Equal(t, fact.expected, result)
	}
}
