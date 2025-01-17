package test

import (
	"testing"
	"time"

	"github.com/mundipagg/boleto-api/models"
	"github.com/stretchr/testify/assert"
)

func Test_StubBoletoRequest_WhenCreateAndSetBoletoRequest_ReturnBoletoRequestSuccessful(t *testing.T) {

	expectedAgreementNumber := uint(1234)
	expectedAgreementAgency := "1234"
	expectedAgreementAccount := "123456"
	expectedAmountInCents := uint64(100)
	expectedExpirationDate := time.Now()
	expectedOurNumber := uint(1234567890)
	expectedRecipientDocument := "12.123.123-0001/11"
	expectedAcceptDivergentAmount := true
	expectedInstructions := "campo de instruções"
	expectedBuyerName := "Nome Do Comprador"
	expectedZipCode := "12345-123"

	s := NewStubBoletoRequest(models.BancoDoBrasil)
	s.WithAgreementNumber(expectedAgreementNumber)
	s.WithAgreementAgency(expectedAgreementAgency)
	s.WithAgreementAccount(expectedAgreementAccount)
	s.WithAmountInCents(uint64(expectedAmountInCents))
	s.WithExpirationDate(expectedExpirationDate)
	s.WithOurNumber(expectedOurNumber)
	s.WithRecipientDocumentNumber(expectedRecipientDocument)
	s.WithAcceptDivergentAmount(expectedAcceptDivergentAmount)
	s.WithInstructions(expectedInstructions)
	s.WithBuyerName(expectedBuyerName)
	s.WithBuyerZipCode(expectedZipCode)

	b := s.Build()

	assert.NotEmpty(t, b)
	assert.Equal(t, expectedAgreementNumber, s.Agreement.AgreementNumber)
	assert.Equal(t, expectedAgreementAgency, s.Agreement.Agency)
	assert.Equal(t, expectedAgreementAccount, s.Agreement.Account)
	assert.Equal(t, expectedAmountInCents, s.Title.AmountInCents)
	assert.Equal(t, expectedExpirationDate, s.Title.ExpireDateTime)
	assert.Equal(t, expectedOurNumber, s.Title.OurNumber)
	assert.Equal(t, expectedAcceptDivergentAmount, s.Title.Rules.AcceptDivergentAmount)
	assert.Equal(t, expectedRecipientDocument, s.Recipient.Document.Number)
	assert.Equal(t, expectedInstructions, s.Title.Instructions)
	assert.Equal(t, expectedBuyerName, s.Buyer.Name)
	assert.Equal(t, expectedZipCode, s.Buyer.Address.ZipCode)
}
