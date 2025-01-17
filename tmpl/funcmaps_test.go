package tmpl

import (
	"fmt"
	"html/template"
	"testing"
	"time"

	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/test"
	"github.com/stretchr/testify/assert"
)

type TestFee struct {
	line          uint8
	AmountFee     uint64
	PercentageFee float64
	TitleAmount   uint64
}

var formatDigitableLineParameters = []test.Parameter{
	{Input: "34191123456789010111213141516171812345678901112", Expected: "34191.12345 67890.101112 13141.516171 8 12345678901112"},
}

var truncateParameters = []test.Parameter{
	{Input: "00000000000000000000", Length: 5, Expected: "00000"},
	{Input: "00000000000000000000", Length: 50, Expected: "00000000000000000000"},
	{Input: "Rua de teste para o truncate", Length: 20, Expected: "Rua de teste para o "},
	{Input: "", Length: 50, Expected: ""},
}

var clearStringParameters = []test.Parameter{
	{Input: "óláçñê", Expected: "olacne"},
	{Input: "ola", Expected: "ola"},
	{Input: "", Expected: ""},
	{Input: "Jardim Novo Cambuí ", Expected: "Jardim Novo Cambui"},
	{Input: "Jardim Novo Cambuí�", Expected: "Jardim Novo Cambui"},
}

var formatNumberParameters = []test.UInt64TestParameter{
	{Input: 50332, Expected: "503,32"},
	{Input: 55, Expected: "0,55"},
	{Input: 0, Expected: "0,00"},
}

var toFloatStrParameters = []test.UInt64TestParameter{
	{Input: 50332, Expected: "503.32"},
	{Input: 55, Expected: "0.55"},
	{Input: 0, Expected: "0.00"},
	{Input: 200, Expected: "2.00"},
}

var StrtoFloatParameters = []test.Parameter{
	{Input: "2.00", Expected: 2.00},
	{Input: "2.01", Expected: 2.01},
}

var formatDocParameters = []test.Parameter{
	{Input: models.Document{Type: "CPF", Number: "12312100100"}, Expected: "123.121.001-00"},
	{Input: models.Document{Type: "CNPJ", Number: "12123123000112"}, Expected: "12.123.123/0001-12"},
}

var docTypeParameters = []test.Parameter{
	{Input: models.Document{Type: "CPF", Number: "12312100100"}, Expected: 1},
	{Input: models.Document{Type: "CNPJ", Number: "12123123000112"}, Expected: 2},
}

var sanitizeCepParameters = []test.Parameter{
	{Input: "25368-100", Expected: "25368100"},
	{Input: "25368100", Expected: "25368100"},
}

var mod11BradescoShopFacilDvParameters = []test.Parameter{
	{Input: "00000000006", Expected: "0"},
	{Input: "00000000001", Expected: "P"},
	{Input: "00000000002", Expected: "8"},
}

var sanitizeCitibankSpecialCharacteresParameters = []test.Parameter{
	{Input: "", Length: 0, Expected: ""},       //Default string value
	{Input: "   ", Length: 3, Expected: "   "}, //Whitespaces
	{Input: "a b", Length: 3, Expected: "a b"},
	{Input: "/-;@", Length: 4, Expected: "/-;@"}, //Caracteres especiais aceitos pelo Citibank
	{Input: "???????????????????????????a-zA-Z0-9ÁÉÍÓÚÀÈÌÒÙÂÊÎÔÛÃÕáéíóúàèìòùâêîôûãõç.", Length: 45, Expected: "a-zA-Z0-9AEIOUAEIOUAEIOUAOaeiouaeiouaeiouaoc."},
	{Input: "Ol@ Mundo. você pode ver uma barra /, mas não uma exclamação!?; Nem Isso", Length: 60, Expected: "Ol@ Mundo. voce pode ver uma barra / mas nao uma exclamacao;"},
	{Input: "Avenida Andr? Rodrigues de Freitas", Length: 33, Expected: "Avenida Andr Rodrigues de Freitas"},
}

var clearStringCaixaParameters = []test.Parameter{
	{Input: "CaixaAccepted:ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789,/()*=-+!:?.;_'", Expected: "CaixaAccepted:ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789,/()*=-+!:?.;_'"},
	{Input: "CaixaAccepted:,/()*=-+!:?.;_'", Expected: "CaixaAccepted:,/()*=-+!:?.;_'"},
	{Input: "XMLNotAccepted:&<>", Expected: "XMLNotAccepted:   "},
	{Input: "CaixaClearCharacter:ÁÉÍÓÚÀÈÌÒÙÂÊÎÔÛÃÕáéíóúàèìòùâêîôûãõç", Expected: "CaixaClearCharacter:AEIOUAEIOUAEIOUAOaeiouaeiouaeiouaoc"},
	{Input: "@#$%¨{}[]^~|ºª§°¹²³£¢¬\\�\"", Expected: "                         "},
}

var truncateOnlyParameters = []test.Parameter{
	{Input: "0000000000000000000", Length: 5, Expected: "00000"},
	{Input: "0000000000000000000", Length: 50, Expected: "0000000000000000000"},
	{Input: "Rua de teste para o truncate", Length: 20, Expected: "Rua de teste para o "},
	{Input: "", Length: 50, Expected: ""},
	{Input: "ÁÉÍÓÚÀÈÌÒÙÂÊÎÔÛÃÕáéíóúàèìòùâêîôûãõç,/()*&=-+!:?<>.;_'@#$%¨{}[]^~|ºª§°¹²³£¢¬\\\"", Length: 80, Expected: "ÁÉÍÓÚÀÈÌÒÙÂÊÎÔÛÃÕáéíóúàèìòùâêîôûãõç,/()*&=-+!:?<>.;_'@#$%¨{}[]^~|ºª§°¹²³£¢¬\\\""},
	{Input: "ÁÉÍÓÚÀÈÌÒÙÂÊÎÔÛÃÕáéíóúàèìòùâêîôûãõç,/()*&=-+!:?<>.;_'@#$%¨{}[]^~|ºª§°¹²³£¢¬\\\"", Length: 75, Expected: "ÁÉÍÓÚÀÈÌÒÙÂÊÎÔÛÃÕáéíóúàèìòùâêîôûãõç,/()*&=-+!:?<>.;_'@#$%¨{}[]^~|ºª§°¹²³£¢¬"},
}

var joinStringsParameters = []test.Parameter{
	{Input: []string{"a", "b", "c"}, Expected: "a b c"},
	{Input: []string{"abc", "d", "efgh"}, Expected: "abc d efgh"},
	{Input: []string{" ", " ", " "}, Expected: "     "},
	{Input: []string{"", "", "", "", ""}, Expected: "    "},
	{Input: []string{"ÁÉÍÓÚÀ", "()*&=-+!:?<", "^~|ºª§°¹²³£¢¬\\\"", "@#$%"}, Expected: "ÁÉÍÓÚÀ ()*&=-+!:?< ^~|ºª§°¹²³£¢¬\\\" @#$%"},
}

var float64ToStringParameters = []test.Parameter{
	{Input: -2.0, Expected: "-2.00"},
	{Input: -2.4, Expected: "-2.40"},
	{Input: 0.0, Expected: "0.00"},
	{Input: 0.02, Expected: "0.02"},
	{Input: 1.0, Expected: "1.00"},
	{Input: 1.23, Expected: "1.23"},
	{Input: 1.2379, Expected: "1.24"},
}

var calculateFeesParameters = []test.Parameter{
	{Input: TestFee{line: 1, AmountFee: 0, PercentageFee: 0, TitleAmount: 1}, Expected: 0.0},
	{Input: TestFee{line: 2, AmountFee: 200, PercentageFee: 0, TitleAmount: 2000}, Expected: 2.0},
	{Input: TestFee{line: 3, AmountFee: 346, PercentageFee: 0, TitleAmount: 2000}, Expected: 3.46},
	{Input: TestFee{line: 4, AmountFee: 122211, PercentageFee: 0, TitleAmount: 2000}, Expected: 1222.1100000000001},
	{Input: TestFee{line: 5, AmountFee: 0, PercentageFee: 1.00, TitleAmount: 2000}, Expected: 0.2},
	{Input: TestFee{line: 6, AmountFee: 0, PercentageFee: 1.26, TitleAmount: 2248}, Expected: 0.283248},
}

var getFineInstructionParameters = []test.Parameter{
	{Input: models.Title{AmountInCents: 2000, Fees: &models.Fees{Fine: &models.Fine{DaysAfterExpirationDate: 1, AmountInCents: 200}}}, Expected: "A PARTIR DE 10/03/2022: MULTA..........R$ 2.00"},
	{Input: models.Title{AmountInCents: 2000, Fees: &models.Fees{Fine: &models.Fine{DaysAfterExpirationDate: 2, AmountInCents: 346}}}, Expected: "A PARTIR DE 11/03/2022: MULTA..........R$ 3.46"},
	{Input: models.Title{AmountInCents: 2000, Fees: &models.Fees{Fine: &models.Fine{DaysAfterExpirationDate: 2, AmountInCents: 122211}}}, Expected: "A PARTIR DE 11/03/2022: MULTA..........R$ 1222.11"},
	{Input: models.Title{AmountInCents: 2000, Fees: &models.Fees{Fine: &models.Fine{DaysAfterExpirationDate: 2, PercentageOnTotal: 1.00}}}, Expected: "A PARTIR DE 11/03/2022: MULTA..........R$ 0.20"},
	{Input: models.Title{AmountInCents: 2248, Fees: &models.Fees{Fine: &models.Fine{DaysAfterExpirationDate: 2, PercentageOnTotal: 1.26}}}, Expected: "A PARTIR DE 11/03/2022: MULTA..........R$ 0.28"},
	{Input: models.Title{AmountInCents: 10, Fees: &models.Fees{Fine: &models.Fine{DaysAfterExpirationDate: 1, PercentageOnTotal: 0.5}}}, Expected: "A PARTIR DE 10/03/2022: MULTA..........R$ 0.00"},
}

var calculateInterestByDayParameters = []test.Parameter{
	{Input: TestFee{line: 1, AmountFee: 0, PercentageFee: 0, TitleAmount: 1}, Expected: 0.0},
	{Input: TestFee{line: 2, AmountFee: 200, PercentageFee: 0, TitleAmount: 2000}, Expected: 2.0},
	{Input: TestFee{line: 3, AmountFee: 346, PercentageFee: 0, TitleAmount: 2000}, Expected: 3.46},
	{Input: TestFee{line: 4, AmountFee: 122211, PercentageFee: 0, TitleAmount: 2000}, Expected: 1222.1100000000001},
	{Input: TestFee{line: 5, AmountFee: 0, PercentageFee: 1.00, TitleAmount: 6000}, Expected: 0.02},
	{Input: TestFee{line: 6, AmountFee: 0, PercentageFee: 1.26, TitleAmount: 83448}, Expected: 0.3504816},
}

var getInterestInstructionParameters = []test.Parameter{
	{Input: models.Title{AmountInCents: 2000, Fees: &models.Fees{Interest: &models.Interest{DaysAfterExpirationDate: 1, AmountPerDayInCents: 20}}}, Expected: "A PARTIR DE 10/03/2022: JUROS POR DIA DE ATRASO.........R$ 0.200"},
	{Input: models.Title{AmountInCents: 2000, Fees: &models.Fees{Interest: &models.Interest{DaysAfterExpirationDate: 1, AmountPerDayInCents: 346}}}, Expected: "A PARTIR DE 10/03/2022: JUROS POR DIA DE ATRASO.........R$ 3.460"},
	{Input: models.Title{AmountInCents: 2000, Fees: &models.Fees{Interest: &models.Interest{DaysAfterExpirationDate: 2, AmountPerDayInCents: 122211}}}, Expected: "A PARTIR DE 11/03/2022: JUROS POR DIA DE ATRASO.........R$ 1222.110"},
	{Input: models.Title{AmountInCents: 6000, Fees: &models.Fees{Interest: &models.Interest{DaysAfterExpirationDate: 2, PercentagePerMonth: 1.00}}}, Expected: "A PARTIR DE 11/03/2022: JUROS POR DIA DE ATRASO.........R$ 0.020"},
	{Input: models.Title{AmountInCents: 83448, Fees: &models.Fees{Interest: &models.Interest{DaysAfterExpirationDate: 2, PercentagePerMonth: 1.26}}}, Expected: "A PARTIR DE 11/03/2022: JUROS POR DIA DE ATRASO.........R$ 0.350"},
	{Input: models.Title{AmountInCents: 10, Fees: &models.Fees{Interest: &models.Interest{DaysAfterExpirationDate: 1, PercentagePerMonth: 0.5}}}, Expected: "A PARTIR DE 10/03/2022: JUROS POR DIA DE ATRASO.........R$ 0.000"},
	{Input: models.Title{AmountInCents: 2000, Fees: &models.Fees{Interest: &models.Interest{DaysAfterExpirationDate: 1, PercentagePerMonth: 1.0}}}, Expected: "A PARTIR DE 10/03/2022: JUROS POR DIA DE ATRASO.........R$ 0.006"},
}

var alphanumericsStringsParameters = []test.Parameter{
	{Input: "zÁÉÍÓÚÀÈÌÒÙÂÊÎÔÛÃÕáéíóúàèìòùâêîôûãõç", Expected: "zÁÉÍÓÚÀÈÌÒÙÂÊÎÔÛÃÕáéíóúàèìòùâêîôûãõç"},
	{Input: "1234567890zÁÉÍÓÚÀÈÌÒÙÂÊÎÔÛÃÕáéíóúàèìòùâêîôûãõç", Expected: "1234567890zÁÉÍÓÚÀÈÌÒÙÂÊÎÔÛÃÕáéíóúàèìòùâêîôûãõç"},
	{Input: "1234567890zÁÉÍÓÚÀÈÌÒÙÂÊÎÔÛÃÕáéíóúàèìòùâêîôûãõçabcdefABCDEFzZ", Expected: "1234567890zÁÉÍÓÚÀÈÌÒÙÂÊÎÔÛÃÕáéíóúàèìòùâêîôûãõçabcdefABCDEFzZ"},
	{Input: "1@234#567890zÁÉÍÓÚÀÈ%ÌÒÙÂÊÎÔÛÃÕáéíóúàèìòùâêîôûãõçabcdefABCDEFzZ", Expected: "1234567890zÁÉÍÓÚÀÈÌÒÙÂÊÎÔÛÃÕáéíóúàèìòùâêîôûãõçabcdefABCDEFzZ"},
}

var alphabeticsStringsParameters = []test.Parameter{
	{Input: "zÁÉÍÓÚÀÈÌÒÙÂÊÎÔÛÃÕáéíóúàèìòùâêîôûãõç", Expected: "zÁÉÍÓÚÀÈÌÒÙÂÊÎÔÛÃÕáéíóúàèìòùâêîôûãõç"},
	{Input: "zÁÉÍ1234ÓÚÀÈÌÒÙÂÊÎÔÛÃÕáéíóú1234àèìòùâêîôûãõç", Expected: "zÁÉÍÓÚÀÈÌÒÙÂÊÎÔÛÃÕáéíóúàèìòùâêîôûãõç"},
	{Input: "abcdefgzÁÉÍÓÚÀÈÌÒÙÂÊÎÔÛÃÕáéíóúàèìòùâêîôûãõçABCDEFGZ", Expected: "abcdefgzÁÉÍÓÚÀÈÌÒÙÂÊÎÔÛÃÕáéíóúàèìòùâêîôûãõçABCDEFGZ"},
	{Input: "abcdef$gzÁÉÍÓÚÀ%ÈÌÒÙÂÊÎÔÛÃ1()ÕáéíóúàèìòùâêîôûãõçABCDEFGZ", Expected: "abcdefgzÁÉÍÓÚÀÈÌÒÙÂÊÎÔÛÃÕáéíóúàèìòùâêîôûãõçABCDEFGZ"},
}

var onlyOneSpaceStringsParameters = []test.Parameter{
	{Input: "campo1 campo2", Expected: "campo1 campo2"},
	{Input: "campo1 \tcampo2", Expected: "campo1 campo2"},
	{Input: "campo1 \t	campo2", Expected: "campo1 campo2"},
	{Input: "campo1 \t	campo2 campo3", Expected: "campo1 campo2 campo3"},
}

var removeAllSpacesStringsParameters = []test.Parameter{
	{Input: "campo1 campo2", Expected: "campo1campo2"},
	{Input: "S P", Expected: "SP"},
	{Input: "S P ", Expected: "SP"},
	{Input: "S 	\tP \t", Expected: "SP"},
}

var roundDownParameters = []test.Parameter{
	{Input: 0.12346, Expected: 0.1234},
	{Input: 0.12349, Expected: 0.1234},
	{Input: 0.12340, Expected: 0.1234},
	{Input: 1.0, Expected: 1.0},
	{Input: 0.0, Expected: 0.0},
}

var float64ToStringTruncateParameters = []test.Parameter{
	{Input: 0.12346, Expected: "0.1234"},
	{Input: 0.12349, Expected: "0.1234"},
	{Input: 0.12340, Expected: "0.1234"},
	{Input: 1.0, Expected: "1.0000"},
	{Input: 0.0, Expected: "0.0000"},
}

func TestShouldPadLeft(t *testing.T) {
	expected := "00005"

	result := padLeft("5", "0", 5)

	assert.Equal(t, expected, result, "O texto deve ter zeros a esqueda e até 5 caracteres")
}

func TestShouldReturnString(t *testing.T) {
	expected := "5"

	result := toString(5)

	assert.Equal(t, expected, result, "O número deve ser uma string")
}

func TestFormatDigitableLine(t *testing.T) {
	for _, fact := range formatDigitableLineParameters {
		result := fmtDigitableLine(fact.Input.(string))
		assert.Equal(t, fact.Expected, result, "A linha digitável deve ser formatada corretamente")
	}
}

func TestTruncate(t *testing.T) {
	for _, fact := range truncateParameters {
		result := truncateString(fact.Input.(string), fact.Length)
		assert.Equal(t, fact.Expected, result, "Deve-se truncar uma string corretamente")
	}
}

func TestClearString(t *testing.T) {
	for _, fact := range clearStringParameters {
		result := clearString(fact.Input.(string))
		assert.Equal(t, fact.Expected, result, "Deve-se limpar uma string corretamente")
	}
}

func TestJoinSpace(t *testing.T) {
	for _, fact := range joinStringsParameters {
		params := fact.Input.([]string)
		result := joinSpace(params...)
		assert.Equal(t, fact.Expected, result, "Deve-se limpar uma string corretamente")
	}
}

func TestFormatCNPJ(t *testing.T) {
	expected := "01.000.000/0001-00"

	result := fmtCNPJ("01000000000100")

	assert.Equal(t, expected, result, "O CNPJ deve ser formatado corretamente")
}

func TestFormatCPF(t *testing.T) {
	expected := "123.121.001-00"

	result := fmtCPF("12312100100")

	assert.Equal(t, expected, result, "O CPF deve ser formatado corretamente")
}

func TestFormatNumber(t *testing.T) {
	for _, fact := range formatNumberParameters {
		result := fmtNumber(fact.Input)
		assert.Equal(t, fact.Expected, result, "O valor em inteiro deve ser convertido para uma string com duas casas decimais separado por vírgula (0,00)")
	}
}

func TestMod11OurNumber(t *testing.T) {
	var expected, onlyDigitExpected uint
	expected = 120000001148
	onlyDigitExpected = 8

	result := calculateOurNumberMod11(12000000114, false)
	onlyDigitResult := calculateOurNumberMod11(12000000114, true)

	assert.Equal(t, expected, result, "Deve-se calcular o mod11 do nosso número e retornar o digito à esquerda")
	assert.Equal(t, onlyDigitExpected, onlyDigitResult, "Deve-se calcular o mod11 do nosso número e retornar o digito à esquerda")
}

func TestToFloatStr(t *testing.T) {
	for _, fact := range toFloatStrParameters {
		result := toFloatStr(fact.Input)
		assert.Equal(t, fact.Expected, result, "O valor em inteiro deve ser convertido para uma string com duas casas decimais separado por ponto (0.00)")
	}
}

func TestStrToFloat(t *testing.T) {
	for _, fact := range StrtoFloatParameters {
		result := strToFloat(fact.Input.(string))
		assert.Equal(t, fact.Expected, result, "O valor em inteiro deve ser convertido para uma string com duas casas decimais separado por ponto (0.00)")
	}
}

func TestFormatDoc(t *testing.T) {
	for _, fact := range formatDocParameters {
		result := fmtDoc(fact.Input.(models.Document))
		assert.Equal(t, fact.Expected, result, "O documento deve ser formatado corretamente")
	}
}

func TestDocType(t *testing.T) {
	for _, fact := range docTypeParameters {
		result := docType(fact.Input.(models.Document))
		assert.Equal(t, fact.Expected, result, "O documento deve ser do tipo correto")
	}
}

func TestTrim(t *testing.T) {
	expected := "hue br festa"

	result := trim(" hue br festa ")

	assert.Equal(t, expected, result, "O texto não deve ter espaços no início e no final")
}

func TestSanitizeHtml(t *testing.T) {
	expected := "hu3 br festa"

	result := sanitizeHtmlString("<b>hu3 br festa</b>")

	assert.Equal(t, expected, result, "O texto não deve conter HTML tags")
}

func TestUnscapeHtml(t *testing.T) {
	var expected template.HTML
	expected = "ó"

	result := unescapeHtmlString("&#243;")

	assert.Equal(t, expected, result, "A string não deve ter caracteres Unicode")
}

func TestSanitizeCep(t *testing.T) {
	for _, fact := range sanitizeCepParameters {
		result := extractNumbers(fact.Input.(string))
		assert.Equal(t, fact.Expected, result, "o zipcode deve conter apenas números")
	}
}

func TestDVOurNumberMod11BradescoShopFacil(t *testing.T) {
	wallet := "19"
	for _, fact := range mod11BradescoShopFacilDvParameters {
		result := mod11BradescoShopFacilDv(fact.Input.(string), wallet)
		assert.Equal(t, fact.Expected, result, "o dígito verificador deve ser equivalente ao OurNumber")
	}
}

func TestEscape(t *testing.T) {
	expected := "KM 5,00    "

	result := escapeStringOnJson("KM 5,00 \t \f \r \b")

	assert.Equal(t, expected, result, "O texto deve ser escapado")
}

func TestRemoveCharacterSpecial(t *testing.T) {
	expected := "Texto com carácter especial   -"

	result := removeSpecialCharacter("Texto? com \"carácter\" especial * ' -")

	assert.Equal(t, expected, result, "Os caracteres especiais devem ser removidos")
}

func TestCitiBankSanitizeString(t *testing.T) {
	for _, fact := range sanitizeCitibankSpecialCharacteresParameters {
		input := fact.Input.(string)
		result := sanitizeCitibankSpecialCharacteres(input, fact.Length)
		assert.Equal(t, fact.Expected, result, "Caracteres especiais e acentos devem ser removidos")
		assert.Equal(t, fact.Length, len(result), "O texto deve ser devidamente truncado")
	}
}

func TestClearStringCaixa(t *testing.T) {
	for _, fact := range clearStringCaixaParameters {
		result := clearStringCaixa(fact.Input.(string))
		assert.Equal(t, fact.Expected, result, "Deve-se limpar uma string corretamente")
	}
}

func TestTruncateOnly(t *testing.T) {
	for _, fact := range truncateOnlyParameters {
		result := truncateOnly(fact.Input.(string), fact.Length)
		assert.Equal(t, fact.Expected, result, "Deve-se truncar uma string corretamente")
	}
}

func TestDatePlusDays(t *testing.T) {
	var daysToAdd uint = 2
	dateNow := time.Now()
	dateExpected := dateNow.UTC().Add(time.Hour * 24 * time.Duration(daysToAdd))

	result := datePlusDays(dateNow, daysToAdd)

	assert.Equal(t, dateExpected, result, "Deve incrementar na data, os dias passados no método")
}

func TestDatePlusDaysConsideringZeroAsStart(t *testing.T) {
	var daysPassedForTheMethod uint = 3
	var daysCountingFromZero uint = 2
	dateNow := time.Now()
	dateExpected := dateNow.UTC().Add(time.Hour * 24 * time.Duration(daysCountingFromZero))

	result := datePlusDaysConsideringZeroAsStart(dateNow, daysPassedForTheMethod)

	assert.Equal(t, dateExpected, result, "Deve incrementar na data, os dias passados no método, considerando o zero como start")
}

func TestFloat64ToStringWith2f(t *testing.T) {
	format := "%.2f"
	for _, fact := range float64ToStringParameters {
		result := float64ToString(format, fact.Input.(float64))
		assert.Equal(t, fact.Expected, result, "Deve formatar o float com duas casas decimais e retornar como string")
	}
}

func TestDatePlusDaysLocalTime(t *testing.T) {
	var daysToAdd uint = 2
	dateNow := time.Now()
	dateExpected := dateNow.Add(time.Hour * 24 * time.Duration(daysToAdd))

	result := datePlusDaysLocalTime(dateNow, daysToAdd)

	assert.Equal(t, dateExpected, result, "Deve incrementar na data, os dias passados no método, considerando o local time")
}

func TestCalculateFees(t *testing.T) {
	for _, fact := range calculateFeesParameters {
		testFee := fact.Input.(TestFee)
		result := calculateFees(testFee.AmountFee, testFee.PercentageFee, testFee.TitleAmount)
		assert.Equal(t, fact.Expected, result, fmt.Sprintf("CalculateFees - Linha %d: Deve calcular corretamente o Fees", testFee.line))
	}
}

func TestCalculateInterestByDay(t *testing.T) {
	for _, fact := range calculateInterestByDayParameters {
		testFee := fact.Input.(TestFee)
		result := calculateInterestByDay(testFee.AmountFee, testFee.PercentageFee, testFee.TitleAmount)
		assert.Equal(t, fact.Expected, result, fmt.Sprintf("CalculateInterestByDay - Linha %d: Deve calcular corretamente o juros", testFee.line))
	}
}

func TestGetFineInstruction(t *testing.T) {
	expireDateTime, _ := time.Parse("2006-01-02", "2022-03-09")

	for _, fact := range getFineInstructionParameters {
		title := fact.Input.(models.Title)
		title.ExpireDateTime = expireDateTime
		result := getFineInstruction(title)
		assert.Equal(t, fact.Expected, result, "Deve trazer a instrução de multa corretamente")
	}
}

func TestGetInterestInstruction(t *testing.T) {
	expireDateTime, _ := time.Parse("2006-01-02", "2022-03-09")

	for _, fact := range getInterestInstructionParameters {
		title := fact.Input.(models.Title)
		title.ExpireDateTime = expireDateTime
		result := getInterestInstruction(title)
		assert.Equal(t, fact.Expected, result, "Deve trazer a instrução de juros corretamente")
	}
}

func TestOnlyAlphanumerics(t *testing.T) {
	for _, fact := range alphanumericsStringsParameters {
		result := onlyAlphanumerics(fact.Input.(string))
		assert.Equal(t, fact.Expected, result, "Deve manter todos os caracteres alfanuméricos")
	}
}

func TestOnlyAlphabetics(t *testing.T) {
	for _, fact := range alphabeticsStringsParameters {
		result := onlyAlphabetics(fact.Input.(string))
		assert.Equal(t, fact.Expected, result, "Deve manter todos os caracteres alfabéticos")
	}
}

func TestOnlyOneSpace(t *testing.T) {
	for _, fact := range onlyOneSpaceStringsParameters {
		result := onlyOneSpace(fact.Input.(string))
		assert.Equal(t, fact.Expected, result, "Deve manter apenas um espaço entre strings, quando existir")
	}
}

func TestRemoveAllSpaces(t *testing.T) {
	for _, fact := range removeAllSpacesStringsParameters {
		result := removeAllSpaces(fact.Input.(string))
		assert.Equal(t, fact.Expected, result, "Deve remover qualquer espaço entre strings")
	}
}

func TestRoundDown(t *testing.T) {
	decimalPlaces := 4
	for _, fact := range roundDownParameters {
		result := roundDown(fact.Input.(float64), decimalPlaces)
		assert.Equal(t, fact.Expected, result, "Deve truncar o número float na quarta casa decimal")
	}
}

func TestFloat64ToStringTruncate(t *testing.T) {
	decimalPlaces := 4
	numberFormat := "%.4f"
	for _, fact := range float64ToStringTruncateParameters {
		result := float64ToStringTruncate(numberFormat, decimalPlaces, fact.Input.(float64))
		assert.Equal(t, fact.Expected, result, "Converte um número float com 4 decimais")
	}
}

func TestConvertAmountInCentsToPercent(t *testing.T) {
	var totalAmount uint64 = 2000
	var amount uint64 = 1
	percentageExpected := 0.05
	result := convertAmountInCentsToPercent(totalAmount, amount)
	assert.Equal(t, percentageExpected, result, "Deve retornar a quantidade em porcento de amount dado um totalAmount")
}

func TestConvertAmountInCentsToPercentPerDay(t *testing.T) {
	var totalAmount uint64 = 3000
	var amount uint64 = 1
	percentageExpected := 1.0
	result := convertAmountInCentsToPercentPerDay(totalAmount, amount)
	assert.Equal(t, percentageExpected, result, "Deve retornar a quantidade em porcento de amount dado um totalAmount por dia")
}
