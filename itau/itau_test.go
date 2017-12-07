package itau

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
	"BankNumber": 341,
	"Authentication": {
		"Username": "a",
		"Password": "b",
		"AccessKey":"c"
	},
	"Agreement": {
		"Wallet":109,
		"Agency":"0407",
		"Account":"55292",
		"AccountDigit":"6"
	},
	"Title": {
		"ExpireDate": "2017-12-31",
		"AmountInCents": 200			
	},
	"Buyer": {
		"Name": "TESTE",
		"Document": {
			"Type": "CNPJ",
			"Number": "00001234567890"
		}
	},
	"Recipient": {
		"Name": "TESTE",
		"Document": {
			"Type": "CNPJ",
			"Number": "00123456789067"
		}
	}
}
`

func TestRegiterBoletoItau(t *testing.T) {
	env.Config(true, true, true)
	input := new(models.BoletoRequest)
	if err := util.FromJSON(baseMockJSON, input); err != nil {
		t.Fail()
	}
	bank := New()
	go mock.Run("9096")
	time.Sleep(2 * time.Second)
	Convey("deve-se processar um boleto itau com sucesso", t, func() {
		output, err := bank.ProcessBoleto(input)
		So(err, ShouldBeNil)
		So(output.BarCodeNumber, ShouldNotBeEmpty)
		So(output.DigitableLine, ShouldNotBeEmpty)
		So(output.Errors, ShouldBeEmpty)
	})
	input.Title.AmountInCents = 400
	Convey("deve-se processar uma falha no registro de boleto no itau", t, func() {
		output, err := bank.ProcessBoleto(input)
		So(err, ShouldBeNil)
		So(output.BarCodeNumber, ShouldBeEmpty)
		So(output.DigitableLine, ShouldBeEmpty)
		So(output.Errors, ShouldNotBeEmpty)
	})
	input.Title.AmountInCents = 200
	ac := input.Agreement.Account
	input.Agreement.Account = ""
	Convey("deve-se tratar uma validacao de conta no itau", t, func() {
		output, err := bank.ProcessBoleto(input)
		So(err, ShouldBeNil)
		So(output.BarCodeNumber, ShouldBeEmpty)
		So(output.DigitableLine, ShouldBeEmpty)
		So(output.Errors, ShouldNotBeEmpty)
	})
	input.Agreement.Account = ac
	input.Authentication.Username = ""
	Convey("deve-se tratar uma falha de login no itau", t, func() {
		output, err := bank.ProcessBoleto(input)
		So(err, ShouldNotBeNil)
		So(output.BarCodeNumber, ShouldBeEmpty)
		So(output.DigitableLine, ShouldBeEmpty)
	})

}
