//go:build !unit
// +build !unit

package test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/mundipagg/boleto-api/models"
	"github.com/pact-foundation/pact-go/dsl"
	"github.com/pact-foundation/pact-go/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var boleto *models.BoletoView

// The actual Provider test itself
func TestMessageProvider_Success(t *testing.T) {
	pact := createPact()

	// Map test descriptions to message producer (handlers)
	functionMappings := dsl.MessageHandlers{
		"a boleto": func(m dsl.Message) (interface{}, error) {
			if boleto != nil {
				return boleto, nil
			} else {
				return map[string]string{
					"message": "not found",
				}, nil
			}
		},
	}

	stateMappings := dsl.StateHandlers{
		"boleto with id exists": func(s dsl.State) error {
			id := primitive.NewObjectID()
			boleto = &models.BoletoView{
				ID:        id,
				UID:       "d4862472-0f5e-11ec-844b-00059a3c7a00",
				SecretKey: "d4862472-0f5e-11ec-844b-00059a3c7a00",
				PublicKey: "00c9acbc0ad001fd0316c83aeb8ef8d6dccefe2a76698963b305f478f7b0e263",
				Format:    "json",
				Boleto: models.BoletoRequest{
					Authentication: models.Authentication{
						Username:           "eyJpZCI6IjgwNDNiNTMtZjQ5Mi00YyIsImNvZGlnb1B1YmxpY2Fkb3IiOjEwOSwiY29kaWdvU29mdHdhcmUiOjEsInNlcXVlbmNpYWxJbnN0YWxhY2FvIjoxfQ",
						Password:           "eyJpZCI6IjBjZDFlMGQtN2UyNC00MGQyLWI0YSIsImNvZGlnb1B1YmxpY2Fkb3IiOjEwOSwiY29kaWdvU29mdHdhcmUiOjEsInNlcXVlbmNpYWxJbnN0YWxhY2FvIjoxLCJzZXF1ZW5jaWFsQ3JlZGVuY2lhbCI6MX0",
						AuthorizationToken: "NWVjNWI5MGI3M2Y0OGM4YTVjNDRlYzBkOmFiNjE3YzM4LTliNzAtNGE1OS1hMzhmLTMzMTU0ZmFiMDEwYw==",
						AccessKey:          "3aaa712f-8b31-457d-abf3-56cdf7f71fe4",
					},
					Agreement: models.Agreement{
						AgreementNumber: 1660590,
						Wallet:          1,
						WalletVariation: 35,
						Agency:          "1234",
						AgencyDigit:     "3",
						Account:         "01234567",
						AccountDigit:    "9",
					},
					Title: models.Title{
						CreateDate:     time.Time{},
						ExpireDateTime: time.Time{},
						ExpireDate:     "2021-09-11",
						AmountInCents:  200,
						OurNumber:      94614469,
						Instructions:   "NÃO RECEBER APÓS O VENCIMENTO. O prazo de compensação de boleto é de até 3 dias úteis após o pagamento, o valor do limite poderá ficar bloqueado até o processamento.",
						DocumentNumber: "1234567890",
						NSU:            "123",
						BoletoType:     "ND",
						Fees: &models.Fees{
							Fine: &models.Fine{
								DaysAfterExpirationDate: 1,
								AmountInCents:           200,
								PercentageOnTotal:       1,
							},
							Interest: &models.Interest{
								DaysAfterExpirationDate: 1,
								AmountPerDayInCents:     200,
								PercentagePerMonth:      1,
							},
						},
						BoletoTypeCode: "19",
					},
					Recipient: models.Recipient{
						Name: "Nome do Recebedor (Loja)",
						Document: models.Document{
							Type:   "CNPJ",
							Number: "11123123000199",
						},
						Address: models.Address{
							Street:     "Logradouro do Recebedor",
							Number:     "1000",
							Complement: "Sala 01",
							ZipCode:    "00000000",
							City:       "Cidade do Recebedor",
							District:   "Bairro do Recebdor",
							StateCode:  "RJ",
						},
					},
					PayeeGuarantor: &models.PayeeGuarantor{
						Name: "Nome do PayeeGuarantor (Loja)",
						Document: models.Document{
							Type:   "CPF",
							Number: "11282705792",
						},
					},
					Buyer: models.Buyer{
						Name:  "Nome do Comprador (Cliente)",
						Email: "teste@pagar.me",
						Document: models.Document{
							Type:   "CPF",
							Number: "12332279717",
						},
						Address: models.Address{
							Street:     "Logradouro do Comprador",
							Number:     "1000",
							Complement: "Casa 01",
							ZipCode:    "01001000",
							City:       "Cidade do Comprador",
							District:   "Bairro do Comprador",
							StateCode:  "SC",
						},
					},
					BankNumber: 1,
					RequestKey: "d3296f2d-781f-4caf-85e6-c30a7d85b30f",
				},
				BankID:        1,
				CreateDate:    time.Time{},
				BankNumber:    "001-9",
				DigitableLine: "00190.00009 01014.051005 00066.673179 9 71340000010000",
				OurNumber:     "12345",
				Barcode:       "00199713400000100000000001014051000006667317",
				Barcode64:     "asdasd",
				Links: []models.Link{
					{
						Href:   "http://localhost:3000/boleto?fmt=html&id=6136912227b87a73be9fb9cc&pk=00c9acbc0ad001fd0316c83aeb8ef8d6dccefe2a76698963b305f478f7b0e263",
						Rel:    "html",
						Method: "GET",
					},
					{
						Href:   "http://localhost:3000/boleto?fmt=html&id=6136912227b87a73be9fb9cc&pk=00c9acbc0ad001fd0316c83aeb8ef8d6dccefe2a76698963b305f478f7b0e263",
						Rel:    "html",
						Method: "GET",
					},
				},
			}

			return nil
		},
	}

	// Verify the Provider with Pactflow publish contracts
	//nolint
	pact.VerifyMessageProvider(t, dsl.VerifyMessageRequest{
		PactURLs:  []string{os.Getenv("PACT_URL")},
		BrokerURL: os.Getenv("PACT_BROKER_URL"),
		ConsumerVersionSelectors: []types.ConsumerVersionSelector{
			{
				Tag:         os.Getenv("GIT_BRANCH"),
				FallbackTag: "master",
				Latest:      true,
			},
		},
		BrokerToken:                os.Getenv("PACT_BROKER_TOKEN"),
		PublishVerificationResults: true,
		ProviderVersion:            os.Getenv("GITHUB_COMMIT"),
		ProviderTags:               []string{os.Getenv("GITHUB_BRANCH")},
		MessageHandlers:            functionMappings,
		StateHandlers:              stateMappings,
	})
}

// Configuration / Test Data
var dir, _ = os.Getwd()

var logDir = fmt.Sprintf("%s/../log", dir)

// Setup the Pact client.
func createPact() dsl.Pact {
	return dsl.Pact{
		Provider: "boleto-api",
		LogDir:   logDir,
	}
}
