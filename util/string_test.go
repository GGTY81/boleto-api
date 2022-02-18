package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type UtilTestParameter struct {
	Input    interface{}
	Expected interface{}
}

var padLeftParameters = []UtilTestParameter{
	{Input: "123", Expected: "0000000123"},
	{Input: "1234567890", Expected: "1234567890"},
}

var digitParameters = []UtilTestParameter{
	{Input: "0123456789", Expected: true},
	{Input: " ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz,/()*&=-+!:?<>.;_\"", Expected: false},
}

var basicCharacter = []UtilTestParameter{
	{Input: " 0123456789,/()*&=-+!:?<>.;_\"", Expected: false},
	{Input: "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz", Expected: true},
}

var caixaSpecialCharacter = []UtilTestParameter{
	{Input: " ,/()*=-+!:?.;_'", Expected: true},
	{Input: "01223456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz@#$%¨{}[]^~Çç\"&<>\\", Expected: false},
}

func TestPadLeft(t *testing.T) {
	length := 10
	paddingCaracter := "0"

	for _, fact := range padLeftParameters {
		result := PadLeft(fact.Input.(string), paddingCaracter, uint(length))
		assert.Equal(t, fact.Expected, result, "O numero deve ser ajustado corretamente")
	}
}

func TestIsDigit(t *testing.T) {
	for _, fact := range digitParameters {
		s := fact.Input.(string)
		for _, c := range s {
			result := IsDigit(c)
			assert.Equal(t, fact.Expected, result, "A verificação de dígito deve ocorrer corretamente")
		}
	}
}

func TestIsBasicCharacter(t *testing.T) {
	for _, fact := range basicCharacter {
		s := fact.Input.(string)
		for _, c := range s {
			result := IsBasicCharacter(c)
			assert.Equal(t, fact.Expected, result, "A verificação de caracter deve ocorrer corretamente")
		}
	}
}

func TestIsSpecialCharacterCaixa(t *testing.T) {
	for _, fact := range caixaSpecialCharacter {
		s := fact.Input.(string)
		for _, c := range s {
			result := IsCaixaSpecialCharacter(c)
			assert.Equal(t, fact.Expected, result, "A verificação de caracter deve ocorrer corretamente")
		}
	}
}

func TestStringfy(t *testing.T) {
	expected := `{"Input":"Texto","Expected":1234}`

	input := UtilTestParameter{
		Input:    "Texto",
		Expected: 1234,
	}

	result := Stringify(input)

	assert.Equal(t, expected, result)
}

func TestParseJson(t *testing.T) {
	input := `{"Input":"Texto","Expected":1234.0}`

	result := ParseJSON(input, new(UtilTestParameter)).(*UtilTestParameter)

	assert.Equal(t, "Texto", result.Input)
	assert.Equal(t, 1234.0, result.Expected)
}

func TestMinifyString(t *testing.T) {
	input := `<html>
			 	<body>
					<p><b>Get My PDF</b></p>
				</body>
			</html>`

	expected := `<html><body><p><b>Get My PDF</b></p></body></html>`

	result := MinifyString(input, "text/html")

	assert.Equal(t, expected, result)

	input = `{
				"Input":"Texto",
				"Expected":1234.0
			 }`
	expected = `{"Input":"Texto","Expected":1234.0}`

	result = MinifyString(input, "application/json")

	assert.Equal(t, expected, result)

}

func TestSanitizeBody(t *testing.T) {
	input := `{
    "bankNumber": 174,
    "authentication": {
            "Username": "altsa",
            "Password": "altsa"
	},
	"agreement": {
		"agreementNumber": 267,
		"wallet": 36,
		"agency": "00000"
	},
	"title": {           
		"expireDate": "2050-12-30",
		"amountInCents": 200,
		"ourNumber": 1,
		"instructions": "Não receber após a data de vencimento.",
		"documentNumber": "1234567890"
	},
	"recipient": {
		"name": "Empresa - Boletos",
		"document": {
			"type": "CNPJ",
			"number": "29799428000128"
		},
		"address": {
			"street": "Avenida Miguel Estefno, 2394",
			"complement": "Água Funda",
			"zipCode": "04301-002",
			"city": "São Paulo",
			"stateCode": "SP"
		}
	},
	"buyer": {
		"name": "Usuario Teste",
		"email": "p@p.com",
		"document": {
			"type": "CNPJ",
			"number": "29.799.428/0001-28"
		},
		"address": {
			"street": "Rua Teste",
			"number": "2",
			"complement": "SALA 1",
			"zipCode": "20931-001",
			"district": "Centro",
			"city": "Rio de Janeiro",
			"stateCode": "RJ"
		}
	}
}`

	result := SanitizeBody(input)
	assert.NotContains(t, result, "\t")
}
