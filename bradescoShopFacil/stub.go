package bradescoShopFacil

import (
	"time"

	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/test"
)

const day = time.Hour * 24

type stubBoletoRequestBradescoShopFacil struct {
	test.StubBoletoRequest
}

//newStubBoletoRequestBradescoShopFacil Cria um novo StubBoletoRequest com valores default validáveis para o BradescoShopFacil
func newStubBoletoRequestBradescoShopFacil() *stubBoletoRequestBradescoShopFacil {
	expirationDate := time.Now().Add(5 * day)

	base := test.NewStubBoletoRequest(models.Bradesco)
	s := &stubBoletoRequestBradescoShopFacil{
		StubBoletoRequest: *base,
	}

	s.Authentication = models.Authentication{
		Username: "55555555555",
		Password: "55555555555555555",
	}

	s.Agreement = models.Agreement{
		AgreementNumber: 55555555,
		Wallet:          25,
		Agency:          "5555",
		Account:         "55555",
	}

	s.Title = models.Title{
		ExpireDateTime: expirationDate,
		ExpireDate:     expirationDate.Format("2006-01-02"),
		OurNumber:      12446688,
		AmountInCents:  200,
		DocumentNumber: "1234566",
		Instructions:   "Senhor caixa, não receber após o vencimento",
	}

	s.Recipient = models.Recipient{
		Name: "TESTE",
		Document: models.Document{
			Type:   "CNPJ",
			Number: "00555555000109",
		},
		Address: models.Address{
			Street:     "TESTE",
			Number:     "111",
			Complement: "TESTE",
			ZipCode:    "11111111",
			City:       "Teste",
			District:   "",
			StateCode:  "SP",
		},
	}

	s.Buyer = models.Buyer{
		Name: "Luke Skywalker",
		Document: models.Document{
			Type:   "CPF",
			Number: "01234567890",
		},
		Address: models.Address{
			Street:     "Mos Eisley Cantina",
			Number:     "123",
			Complement: "Apto",
			ZipCode:    "20001-000",
			City:       "Tatooine",
			District:   "Tijuca",
			StateCode:  "RJ",
		},
	}
	return s
}
