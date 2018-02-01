package bradescoShopFacil

import (
	"testing"
	"time"

	"github.com/mundipagg/boleto-api/env"
	"github.com/mundipagg/boleto-api/mock"
	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/util"
	. "github.com/smartystreets/goconvey/convey"
)

const baseMockJSON = `
{
    "BankNumber": 237,
     "Authentication": {
        "Username": "55555555555",
        "Password": "55555555555555555"
    },
    "Agreement": {
        "AgreementNumber": 55555555,
        "Wallet": 25,
        "Agency":"5555",
        "Account":"55555"
    },
    "Title": {
        "ExpireDate": "2029-08-01",
        "AmountInCents": 200,
        "OurNumber": 12446688,
        "Instructions": "Senhor caixa, não receber após o vencimento",
        "DocumentNumber": "1234566"
    },
    "Buyer": {
        "Name": "Luke Skywalker",
        "Document": {
            "Type": "CPF",
            "Number": "01234567890"
        },
        "Address": {
            "Street": "Mos Eisley Cantina",
            "Number": "123",
            "Complement": "Apto",
            "ZipCode": "20001-000",
            "City": "Tatooine",
            "District": "Tijuca",
            "StateCode": "RJ"
        }
    },
    "Recipient": {
      "Name": "TESTE",
        "Document": {
            "Type": "CNPJ",
            "Number": "00555555000109"
        },

        "Address": {
            "Street": "TESTE",
            "Number": "111",
            "Complement": "TESTE",
            "ZipCode": "11111111",
            "City": "Teste",
            "District": "",
            "StateCode": "SP"
        }

    }
}
`

func TestRegiterBoleto(t *testing.T) {
	env.Config(true, true, true)
	input := new(models.BoletoRequest)
	if err := util.FromJSON(baseMockJSON, input); err != nil {
		t.Fail()
	}
	bank := New()
	go mock.Run("9093")
	time.Sleep(2 * time.Second)

	Convey("deve-se processar um boleto BradescoShopFacil com sucesso", t, func() {
		output, err := bank.ProcessBoleto(input)
		So(err, ShouldBeNil)
		So(output.BarCodeNumber, ShouldNotBeEmpty)
		So(output.DigitableLine, ShouldNotBeEmpty)
		So(output.Errors, ShouldBeEmpty)
	})
	input.Title.AmountInCents = 400
	Convey("deve-se processar um boleto BradescoShopFacil com sucesso", t, func() {
		output, err := bank.ProcessBoleto(input)
		So(err, ShouldBeNil)
		So(output.BarCodeNumber, ShouldBeEmpty)
		So(output.DigitableLine, ShouldBeEmpty)
		So(output.Errors, ShouldNotBeEmpty)
	})
}

func TestBarcodeGenerationBradescoShopFacil(t *testing.T) {
	//example := "23795796800000001990001250012446693212345670"
	boleto := models.BoletoRequest{}
	boleto.BankNumber = models.Bradesco
	boleto.Agreement = models.Agreement{
		Account: "1234567",
		Agency:  "1",
		Wallet:  25,
	}
	expireDate, _ := time.Parse("02-01-2006", "01-08-2019")
	boleto.Title = models.Title{
		AmountInCents:  199,
		OurNumber:      124466932,
		ExpireDateTime: expireDate,
	}
	bc := getBarcode(boleto)
	Convey("deve-se montar o código de barras do BradescoShopFacil", t, func() {
		So(bc.toString(), ShouldEqual, "23795796800000001990001250012446693212345670")
	})
}

func TestRemoveDigitFromAccount(t *testing.T) {
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
	Convey("deve-se montar identificar e remover o digito da conta", t, func() {
		So(bc.toString(), ShouldEqual, "23791796800000001992372250012446693300056000")
	})
}
