package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mundipagg/boleto-api/bank"
	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/test"
	"github.com/stretchr/testify/assert"
)

var handleStoneErrorsParameters = []test.Parameter{
	{Input: models.GetBoletoResponseError("srn:error:validation", "erro de validação"), Expected: http.StatusBadRequest},
	{Input: models.GetBoletoResponseError("srn:error:unauthenticated", "erro de validação do token"), Expected: http.StatusInternalServerError},
	{Input: models.GetBoletoResponseError("srn:error:unauthorized", "erro de clientId"), Expected: http.StatusBadGateway},
	{Input: models.GetBoletoResponseError("srn:error:not_found", "recusro não encontrado"), Expected: http.StatusBadGateway},
	{Input: models.GetBoletoResponseError("srn:error:conflict", "conflito"), Expected: http.StatusBadGateway},
	{Input: models.GetBoletoResponseError("srn:error:product_not_enabled", "produto não disponível"), Expected: http.StatusBadRequest},
	{Input: models.GetBoletoResponseError("unexpect_error", "erro inesperado"), Expected: http.StatusInternalServerError},
}

var handleInternalErrorsParameters = []test.Parameter{
	{Input: models.GetBoletoResponseError("MP400", "falha de validação"), Expected: http.StatusBadRequest},
	{Input: models.GetBoletoResponseError("MPAmountInCents", "valor invalido"), Expected: http.StatusBadRequest},
	{Input: models.GetBoletoResponseError("MPExpireDate", "data de expiração invalida"), Expected: http.StatusBadRequest},
	{Input: models.GetBoletoResponseError("MPBuyerDocumentType", "tipo de documento invalido"), Expected: http.StatusBadRequest},
	{Input: models.GetBoletoResponseError("MPDocumentNumber", "documento invalido"), Expected: http.StatusBadRequest},
	{Input: models.GetBoletoResponseError("MPRecipientDocumentType", "tipo de documento invalido"), Expected: http.StatusBadRequest},
	{Input: models.GetBoletoResponseError("MPTimeout", "timeout"), Expected: http.StatusGatewayTimeout},
	{Input: models.GetBoletoResponseError("MPOurNumberFail", "resposta sem nosso numero"), Expected: http.StatusBadGateway},
}

func Test_GetErrorCodeToLog_WhenHasError_ReturnErrorCode(t *testing.T) {
	expectedErrorCode := "CODE"
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	err := models.GetBoletoResponseError(expectedErrorCode, "error")
	c.Set(responseKey, err)

	result := getErrorCodeToLog(c)

	assert.Equal(t, expectedErrorCode, result)
}

func Test_GetErrorCodeToLog_WhitoutError_ReturnEmptyString(t *testing.T) {
	expectedErrorCode := ""
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Set(responseKey, models.BoletoResponse{})

	result := getErrorCodeToLog(c)

	assert.Equal(t, expectedErrorCode, result)
}

func Test_HandleErrors_WhenNotQualifiedForNewHandleError(t *testing.T) {
	expectStatusCode := 200
	c := arrangeContextWithBankAndResponse(models.BancoDoBrasil, models.BoletoResponse{})

	handleErrors(c)

	assert.Equal(t, expectStatusCode, c.Writer.Status())
}

func Test_HandleErrors_WhenStoneBankResponse(t *testing.T) {
	for _, fact := range handleStoneErrorsParameters {
		c := arrangeContextWithBankAndResponse(models.Stone, fact.Input.(models.BoletoResponse))
		handleErrors(c)
		assert.Equal(t, fact.Expected.(int), c.Writer.Status())
	}
}

func Test_HandleErrors_WhenInternalErrorResponse(t *testing.T) {
	for _, fact := range handleInternalErrorsParameters {
		c := arrangeContextWithBankAndResponse(models.Stone, fact.Input.(models.BoletoResponse))
		handleErrors(c)
		assert.Equal(t, fact.Expected.(int), c.Writer.Status())
	}
}

func arrangeContextWithBankAndResponse(bankNumber int, response models.BoletoResponse) *gin.Context {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	b, _ := bank.Get(models.BoletoRequest{BankNumber: models.Stone})
	c.Set(bankKey, b)
	c.Set(responseKey, response)
	return c
}

func Test_RotaRegisterV1_WhenPanicOccurred_RunsPanicRecovery(t *testing.T) {
	router, w := arrangeMiddlewareRoute("/boleto/register", parseBoleto, registerBoletoLogger, errorResponseToClient, panicRecoveryHandler, mockPanicRegistration)
	req, _ := http.NewRequest("POST", "/boleto/register", bytes.NewBuffer([]byte(mockPanicRegistrationRequestJSON)))

	router.ServeHTTP(w, req)

	assert.Equal(t, 500, w.Code)
	assert.Equal(t, mockPanicRegistrationResponseJSON, w.Body.String())
}

func Test_RotaRegisterV2_WhenPanicOccurred_RunsPanicRecovery(t *testing.T) {
	router, w := arrangeMiddlewareRoute("/boleto/register", parseBoleto, registerBoletoLogger, handleErrors, panicRecoveryHandler, mockPanicRegistration)
	req, _ := http.NewRequest("POST", "/boleto/register", bytes.NewBuffer([]byte(mockPanicRegistrationRequestJSON)))

	router.ServeHTTP(w, req)

	assert.Equal(t, 500, w.Code)
	assert.Equal(t, mockPanicRegistrationResponseJSON, w.Body.String())
}

func Test_QualifiedForNewErrorHandling_WhenBankStoneWithError_ReturnTrue(t *testing.T) {
	request := models.BoletoRequest{BankNumber: models.Stone}
	response := models.GetBoletoResponseError("MP000", "error")
	bank, _ := bank.Get(request)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Set(boletoKey, request)
	c.Set(bankKey, bank)

	result := qualifiedForNewErrorHandling(c, response)

	assert.True(t, result)
}

func Test_QualifiedForNewErrorHandling_WhenBankStoneWithoutError_ReturnFalse(t *testing.T) {
	request := models.BoletoRequest{BankNumber: models.Stone}
	response := models.BoletoResponse{}
	bank, _ := bank.Get(request)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Set(boletoKey, request)
	c.Set(bankKey, bank)

	result := qualifiedForNewErrorHandling(c, response)

	assert.False(t, result)
}

func Test_QualifiedForNewErrorHandling_WhenAnotherBankWithError_ReturnFalse(t *testing.T) {
	request := models.BoletoRequest{BankNumber: models.Caixa}
	response := models.GetBoletoResponseError("MP000", "error")
	bank, _ := bank.Get(request)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Set(boletoKey, request)
	c.Set(bankKey, bank)

	result := qualifiedForNewErrorHandling(c, response)

	assert.False(t, result)
}

func Test_QualifiedForNewErrorHandling_WhenAnotherBankWithoutError_ReturnFalse(t *testing.T) {
	request := models.BoletoRequest{BankNumber: models.Caixa}
	response := models.BoletoResponse{}
	bank, _ := bank.Get(request)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Set(boletoKey, request)
	c.Set(bankKey, bank)

	result := qualifiedForNewErrorHandling(c, response)

	assert.False(t, result)
}
