package pefisa

import (
	"time"

	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/test"
)

const day = time.Hour * 24

type stubBoletoRequestPefisa struct {
	test.StubBoletoRequest
}

//newStubBoletoRequestPefisa Cria um novo StubBoletoRequest com valores default validáveis para a Pefisa
func newStubBoletoRequestPefisa() *stubBoletoRequestPefisa {
	expirationDate := time.Now().Add(5 * day)

	base := test.NewStubBoletoRequest(models.Pefisa)
	s := &stubBoletoRequestPefisa{
		StubBoletoRequest: *base,
	}

	s.Agreement = models.Agreement{
		AgreementNumber: 267,
		Wallet:          36,
		Agency:          "00000",
		Account:         "0062145",
	}

	s.Title = models.Title{
		ExpireDateTime: expirationDate,
		ExpireDate:     "2050-12-30",
		OurNumber:      1,
		AmountInCents:  200,
		DocumentNumber: "1234567890",
		Instructions:   "Não receber após a data de vencimento.",
		BoletoType:     "OUT",
		BoletoTypeCode: "99",
	}

	s.Recipient = models.Recipient{
		Document: models.Document{
			Type:   "CNPJ",
			Number: "29799428000128",
		},
	}

	s.Buyer = models.Buyer{
		Name:  "Usuario Teste",
		Email: "p@p.com",
		Document: models.Document{
			Type:   "CNPJ",
			Number: "29.799.428/0001-28",
		},
		Address: models.Address{
			Street:     "Rua Teste",
			Number:     "2",
			Complement: "SALA 1",
			ZipCode:    "20931-001",
			City:       "Rio de Janeiro",
			District:   "Centro",
			StateCode:  "RJ",
		},
	}
	return s
}
