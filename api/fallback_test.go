// +build integration !unit

package api

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mundipagg/boleto-api/bank"
	"github.com/mundipagg/boleto-api/mock"
	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/test"
	"github.com/stretchr/testify/assert"
)

func Test_fallback(t *testing.T) {
	mock.StartMockService("9101")

	var fallback = new(Fallback)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	boleto := *test.NewStubBoletoRequest(models.BancoDoBrasil).Build()
	bank, _ := bank.Get(boleto)
	c.Set(boletoKey, boleto)
	c.Set(bankKey, bank)
	test.CreateClientIP(c)
	payload := `{"ID":"6127b37d36b0e8770b1668ae","uid":"7ce410cd-0682-11ec-852e-00059a3c7a00","secretkey":"7ce410cd-0682-11ec-852e-00059a3c7a00","publickey":"dad58ecd903ceda1ce6e479ff1e6fab399c8207dd000af127cbcdbb5cd3dfe8d","boleto":{"authentication":{},"agreement":{"agreementNumber":1103388,"agency":"3337"},"title":{"createDate":"2021-08-26T00:00:00Z","expireDateTime":"2021-08-31T00:00:00Z","expireDate":"2021-08-31","amountInCents":200,"ourNumber":14000000019047441,"instructions":"NÃO RECEBER APÓS O VENCIMENTO. O prazo de compensação de boleto é de até 3 dias úteis após o pagamento, o valor do limite poderá ficar bloqueado até o processamento.","documentNumber":"12345678901","boletoType":"OUT","BoletoTypeCode":"99"},"recipient":{"name":"Nome do Recebedor (Loja)","document":{"type":"CNPJ","number":"18727053000174"},"address":{"street":"Logradouro do Recebedor","number":"1000","complement":"Sala 01","zipCode":"00000000","city":"Cidade do Recebedor","district":"Bairro do Recebdor","stateCode":"RJ"}},"buyer":{"name":"Nome do Comprador (Cliente)","email":"comprador@gmail.com","document":{"type":"CPF","number":"11282705792"},"address":{"street":"Logradouro do Comprador","number":"1000","complement":"Casa 01","zipCode":"01001000","city":"Cidade do Comprador","district":"Bairro do Comprador","stateCode":"SC"}},"bankNumber":104,"requestKey":"5239ad4a-2a97-4d39-905a-2cc304971d11"},"bankId":104,"createDate":"2021-08-26T12:30:05.2479181-03:00","bankNumber":"104-0","digitableLine":"10492.00650 61000.100042 09922.269841 3 72670000001000","ourNumber":"14000000099222698","barcode":"10493726700000010002006561000100040992226984","links":[{"href":"http://localhost:3000/boleto?fmt=html\u0026id=6127b37d36b0e8770b1668ae\u0026pk=dad58ecd903ceda1ce6e479ff1e6fab399c8207dd000af127cbcdbb5cd3dfe8d","rel":"html","method":"GET"},{"href":"http://localhost:3000/boleto?fmt=pdf\u0026id=6127b37d36b0e8770b1668ae\u0026pk=dad58ecd903ceda1ce6e479ff1e6fab399c8207dd000af127cbcdbb5cd3dfe8d","rel":"pdf","method":"GET"}]}`

	fallback.Save(c, "testeXXX", payload)
}

func Test_getLogUploadProperties(t *testing.T) {
	var time int64 = 100
	var registerId = "RegisterXXX"
	var content = "Content"

	response := getLogUploadProperties(time, registerId, content)

	assert.Equal(t, len(response), 3)
	assert.Equal(t, response["TotalElapsedTimeInMilliseconds"], time)
	assert.Equal(t, response["RegisterId"], registerId)
	assert.Equal(t, response["Content"], content)

}
