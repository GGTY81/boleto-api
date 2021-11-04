package jpmorgan

import (
	"time"

	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/test"
)

const day = time.Hour * 24

type stubBoletoRequestJPMorgan struct {
	test.StubBoletoRequest
}

func newStubBoletoRequestJPMorgan() *stubBoletoRequestJPMorgan {
	expirationDate := time.Now().Add(5 * day)

	base := test.NewStubBoletoRequest(models.JPMorgan)
	s := &stubBoletoRequestJPMorgan{
		StubBoletoRequest: *base,
	}

	s.Title = models.Title{
		ExpireDateTime: expirationDate,
		ExpireDate:     expirationDate.Format("2006-01-02"),
		AmountInCents:  200,
		Instructions:   "Sr. Caixa, favor não receber após vencimento",
		DocumentNumber: "999999999999999",
	}

	s.Recipient = models.Recipient{
		Document: models.Document{
			Type:   "CNPJ",
			Number: "33172537000198",
		},
	}

	s.Buyer = models.Buyer{
		Name:  "Nome do Comprador (Cliente)",
		Email: "",
		Document: models.Document{
			Type:   "CPF",
			Number: "37303489819",
		},
		Address: models.Address{
			Street:     "Logradouro do Comprador",
			Number:     "1000",
			Complement: "Sala 01",
			ZipCode:    "00000000",
			City:       "Cidade do Comprador",
			District:   "Bairro do Comprador",
			StateCode:  "SP",
		},
	}

	return s
}

func (s *stubBoletoRequestJPMorgan) WithBoletoType(bt string) *stubBoletoRequestJPMorgan {
	switch bt {
	case "DM":
		s.Title.BoletoType, s.Title.BoletoTypeCode = bt, "02"
	default:
		s.Title.BoletoType = bt
	}
	return s
}

func (s *stubBoletoRequestJPMorgan) WithDocument(number string, doctype string) *stubBoletoRequestJPMorgan {
	s.Buyer.Document.Type = doctype
	s.Buyer.Document.Number = number
	return s
}

func (s *stubBoletoRequestJPMorgan) WithAccessKey(key string) *stubBoletoRequestJPMorgan {
	s.Authentication.AccessKey = key
	return s
}
