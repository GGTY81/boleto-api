package issuer

import (
	"fmt"
	"testing"

	"github.com/mundipagg/boleto-api/test"
	"github.com/stretchr/testify/assert"
)

var barCodeParameter = []test.Parameter{
	{Input: Issuer{barCode: "03391877100003841119"}, Expected: true},
	{Input: Issuer{barCode: "03391877100003841119794611200000957223350101"}, Expected: true},
	{Input: Issuer{barCode: "0339187710000384111979461120000095722335010103391877100003841119794611200000957223350101"}, Expected: true},
	{Input: Issuer{barCode: ""}, Expected: false},
	{Input: Issuer{barCode: "     "}, Expected: false},
	{Input: Issuer{barCode: "abcdefghijklmnopqrstuvxwyz"}, Expected: false},
	{Input: Issuer{barCode: "a03391877100003841119"}, Expected: false},
	{Input: Issuer{barCode: "03391877100003841119b"}, Expected: false},
	{Input: Issuer{barCode: "033918771C00003841119"}, Expected: false},
	{Input: Issuer{barCode: "033 9187 71000 03841 119"}, Expected: false},
	{Input: Issuer{barCode: "*03391877100003841???1197946112´00000957223350101"}, Expected: false},
}

var digitableLineParameter = []test.Parameter{
	{Input: Issuer{digitableLine: "23792.69307 40004.617383 11000.180908 9 877"}, Expected: true},
	{Input: Issuer{digitableLine: "23792.69307 40004.617383 11000.180908 9 87720000007290"}, Expected: true},
	{Input: Issuer{digitableLine: "23792.69307 40004.617383 11000.180908 9 87720000007290452564156412"}, Expected: true},
	{Input: Issuer{digitableLine: ""}, Expected: false},
	{Input: Issuer{digitableLine: "     "}, Expected: false},
	{Input: Issuer{digitableLine: "a23792.69307 40004.617383 11000.180908 9 87720000007290"}, Expected: false},
	{Input: Issuer{digitableLine: "23792.b69307 40004.617383 11000.180908 9 87720000007290"}, Expected: false},
	{Input: Issuer{digitableLine: "23792.69307 400c04.617383 11000.180908 9 87720000007290"}, Expected: false},
	{Input: Issuer{digitableLine: "23792.69307 40004.617d383 11000.180908 9 87720000007290"}, Expected: false},
	{Input: Issuer{digitableLine: "23792.69307 40004.617383 110e00.180908 9 87720000007290"}, Expected: false},
	{Input: Issuer{digitableLine: "23792.69307 40004.617383 11000.1809FFF0f8 9 87720000007290"}, Expected: false},
	{Input: Issuer{digitableLine: "23792.69307 40004.617383 11000.180908 9g 87720000007290"}, Expected: false},
	{Input: Issuer{digitableLine: "23792.69307 40004.617383 11000.180908 9 877200000h07290"}, Expected: false},
	{Input: Issuer{digitableLine: "23792.69307 40004.617383 11000 9 87720000007290"}, Expected: false},
	{Input: Issuer{digitableLine: "2372.69307 40004.617383 11000.180908 9 87720000007290"}, Expected: false},
	{Input: Issuer{digitableLine: "23792.6937 40004.617383 11000.180908 9 87720000007290"}, Expected: false},
	{Input: Issuer{digitableLine: "23792.69307 4004.617383 11000.180908 9 87720000007290"}, Expected: false},
	{Input: Issuer{digitableLine: "23792.69307 40004.61383 11000.180908 9 87720000007290"}, Expected: false},
	{Input: Issuer{digitableLine: "23792.69307 40004.617383 1100.180908 9 87720000007290"}, Expected: false},
	{Input: Issuer{digitableLine: "23792.69307 40004.617383 11000.18008 9 87720000007290"}, Expected: false},
	{Input: Issuer{digitableLine: "23792.69307 40004.617383 11000.180908  87720000007290"}, Expected: false},
	{Input: Issuer{digitableLine: "23792.69307 40004.617383 11000.180908 9"}, Expected: false},
	{Input: Issuer{digitableLine: "^23792.69307 40004.617383 11000.180908 9 87720000007290"}, Expected: false},
}

var barCodeAndDigitableLineParameters = []test.Parameter{
	{Input: Issuer{barCode: "03391877100003841119794611200000957223350101", digitableLine: "23792.69307 40004.617383 11000.180908 9 87720000007290"}, Expected: true},
	{Input: Issuer{barCode: "0339187abc7100003841119794611200000957223350101", digitableLine: "23792.69307 40004.617383 11000.180908 9 87720000007290"}, Expected: false},
	{Input: Issuer{barCode: "03391877100003841119794611200000957223350101", digitableLine: "23a792.693b07 400c04.61d7383 1100e0.180f908 9g 8772000h0007290"}, Expected: false},
	{Input: Issuer{barCode: "a0339187710000v,^^´[´p3841119794611", digitableLine: "11000.180908 9 87720000007290"}, Expected: false},
}

func TestNewIssuer(t *testing.T) {
	expectedBarCode := "0339187710000384111794611200000957223350101"
	expectedDigitableLine := "2379.69307 40004.617383 11000.180908 9 87720000007290"

	issuer := NewIssuer(expectedBarCode, expectedDigitableLine)

	assert.Equal(t, expectedBarCode, issuer.barCode, "O barcode não foi atribuído corretamente")
	assert.Equal(t, expectedDigitableLine, issuer.digitableLine, "A digitableline não foi atribuído corretamnte")
}

func TestIsValidBarCode(t *testing.T) {
	for _, fact := range barCodeParameter {
		issuer := fact.Input.(Issuer)
		result := issuer.IsValidBarCode()

		assert.Equal(t, fact.Expected, result, fmt.Sprintf("O barCode não é válido %v ", fact.Input))
	}
}

func TestIsValidDigitableLine(t *testing.T) {
	for _, fact := range digitableLineParameter {
		issuer := fact.Input.(Issuer)
		result := issuer.IsValidDigitableLine()

		assert.Equal(t, fact.Expected, result, fmt.Sprintf("A linha digitável não é válida %v ", fact.Input))
	}
}

func TestIsValidDigitableLineAndBarCode(t *testing.T) {
	for _, fact := range barCodeAndDigitableLineParameters {
		issuer := fact.Input.(Issuer)
		result := issuer.IsValidBarCode() && issuer.IsValidDigitableLine()

		assert.Equal(t, fact.Expected, result, fmt.Sprintf("A linha digitávl ou o código de barras não são válidos %v", fact.Input))
	}
}
