package itau

import (
	"time"

	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/test"
)

const day = time.Hour * 24

type stubBoletoRequestItau struct {
	test.StubBoletoRequest
}

//newStubBoletoRequestItau Cria um novo StubBoletoRequest com valores default validáveis para o Itau
func newStubBoletoRequestItau() *stubBoletoRequestItau {
	expirationDate := time.Now().Add(5 * day)

	base := test.NewStubBoletoRequest(models.Itau)
	s := &stubBoletoRequestItau{
		StubBoletoRequest: *base,
	}

	s.Authentication = models.Authentication{
		Username:  "a",
		Password:  "b",
		AccessKey: "c",
	}

	s.Agreement = models.Agreement{
		AgreementNumber: 267,
		Wallet:          109,
		Agency:          "0407",
		Account:         "55292",
		AccountDigit:    "6",
	}

	s.Title = models.Title{
		ExpireDateTime: expirationDate,
		ExpireDate:     "2050-12-30",
		AmountInCents:  200,
	}

	s.Recipient = models.Recipient{
		Document: models.Document{
			Type:   "CNPJ",
			Number: "00123456789067",
		},
	}

	s.Buyer = models.Buyer{
		Name:  "Usuario Teste",
		Email: "p@p.com",
		Document: models.Document{
			Type:   "CNPJ",
			Number: "00001234567890",
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

func (s *stubBoletoRequestItau) WithAuthenticationUserName(userName string) *stubBoletoRequestItau {
	s.Authentication = models.Authentication{
		Username: userName,
	}
	return s
}
