package api

import (
	"net/http"
	"testing"

	"github.com/mundipagg/boleto-api/models"
	"github.com/stretchr/testify/assert"
)

func Test_GetResultFromContext_WhenResultExists_ReturnResult(t *testing.T) {
	c, _, _ := arrangeGetBoleto()
	url := "http://localhost:3000/boleto?fmt=html&id=1234567890&pk=1234567890"
	c.Request, _ = http.NewRequest(http.MethodGet, url, nil)

	expected := models.NewGetBoletoResult(c)
	c.Set(resultGetBoletoKey, expected)

	result := getResultFromContext(c)

	assert.NotNil(t, result)
	assert.EqualValues(t, expected, result, "O resultado não é igual ao esperado")
}

func Test_GetResultFromContext_WhenResultNotExists_ReturnNil(t *testing.T) {
	c, _, _ := arrangeGetBoleto()
	url := "http://localhost:3000/boleto?fmt=html&id=1234567890&pk=1234567890"
	c.Request, _ = http.NewRequest(http.MethodGet, url, nil)

	result := getResultFromContext(c)

	assert.Nil(t, result)
}

func Test_GetRequestKeyFromContext_WhenRequestKeyExists_ReturnHeaderRequestKey(t *testing.T) {
	expectedRequestKey := "00000000-0000-0000-0000-000000000000"

	c, _, _ := arrangeGetBoleto()
	url := "http://localhost:3000/boleto?fmt=html&id=1234567890&pk=1234567890"
	c.Request, _ = http.NewRequest(http.MethodGet, url, nil)
	c.Request.Header.Add("RequestKey", expectedRequestKey)

	result := getRequestKeyFromContext(c)

	assert.Equal(t, expectedRequestKey, result, "A requestKey não é a mesma do contexto")
}

func Test_GetRequestKeyFromContext_WhenRequestKeyNotExists_ReturnNewRequestKey(t *testing.T) {
	c, _, _ := arrangeGetBoleto()
	url := "http://localhost:3000/boleto?fmt=html&id=1234567890&pk=1234567890"
	c.Request, _ = http.NewRequest(http.MethodGet, url, nil)

	result := getRequestKeyFromContext(c)

	assert.NotEqual(t, "", result, "A requestKey está vazia")
}
