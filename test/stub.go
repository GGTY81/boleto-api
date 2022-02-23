package test

import (
	"time"

	"github.com/mundipagg/boleto-api/models"
)

//StubBoletoRequest Stub base para criação de BoletoRequest
type StubBoletoRequest struct {
	BuilderBoletoRequest
	Authentication models.Authentication
	Agreement      models.Agreement
	Title          models.Title
	Recipient      models.Recipient
	PayeeGuarantor *models.PayeeGuarantor
	Buyer          models.Buyer
	bank           models.BankNumber
}

func NewStubBoletoRequest(bank models.BankNumber) *StubBoletoRequest {
	s := &StubBoletoRequest{
		BuilderBoletoRequest: NewBuilderBoletoRequest(),
	}

	s.bank = bank

	s.Authentication = models.Authentication{}
	s.Agreement = models.Agreement{}
	s.Title = models.Title{}
	s.Recipient = models.Recipient{}
	s.Buyer = models.Buyer{}

	return s
}

func (s *StubBoletoRequest) WithAgreementNumber(number uint) *StubBoletoRequest {
	s.Agreement.AgreementNumber = number
	return s
}

func (s *StubBoletoRequest) WithAgreementAgency(agency string) *StubBoletoRequest {
	s.Agreement.Agency = agency
	return s
}

func (s *StubBoletoRequest) WithAgreementAccount(account string) *StubBoletoRequest {
	s.Agreement.Account = account
	return s
}

func (s *StubBoletoRequest) WithAmountInCents(amount uint64) *StubBoletoRequest {
	s.Title.AmountInCents = amount
	return s
}

func (s *StubBoletoRequest) WithOurNumber(ourNumber uint) *StubBoletoRequest {
	s.Title.OurNumber = ourNumber
	return s
}

func (s *StubBoletoRequest) WithExpirationDate(expiredAt time.Time) *StubBoletoRequest {
	s.Title.ExpireDateTime = expiredAt
	s.Title.ExpireDate = expiredAt.Format("2006-01-02")
	return s
}

func (s *StubBoletoRequest) WithDocumentNumber(documentNumber string) *StubBoletoRequest {
	s.Title.DocumentNumber = documentNumber
	return s
}

func (s *StubBoletoRequest) WithInstructions(instructions string) *StubBoletoRequest {
	s.Title.Instructions = instructions
	return s
}

func (s *StubBoletoRequest) WithAcceptDivergentAmount(accepted bool) *StubBoletoRequest {
	if !s.Title.HasRules() {
		s.Title.Rules = &models.Rules{}
	}

	s.Title.Rules.AcceptDivergentAmount = accepted
	return s
}

func (s *StubBoletoRequest) WithRecipientDocumentNumber(docNumber string) *StubBoletoRequest {
	s.Recipient.Document.Number = docNumber
	return s
}

func (s *StubBoletoRequest) WithBuyerName(buyerName string) *StubBoletoRequest {
	s.Buyer.Name = buyerName
	return s
}

func (s *StubBoletoRequest) WithBuyerZipCode(zipcode string) *StubBoletoRequest {
	s.Buyer.Address.ZipCode = zipcode
	return s
}

func (s *StubBoletoRequest) WithRecipientDocumentType(documentType string) *StubBoletoRequest {
	s.Recipient.Document.Type = documentType
	return s
}

func (s *StubBoletoRequest) WithRecipientName(recipientName string) *StubBoletoRequest {
	s.Recipient.Name = recipientName
	return s
}

func (s *StubBoletoRequest) WithPayeeGuarantorName(PayeeGuarantorName string) *StubBoletoRequest {
	s.createStubPayeeGuarantor()
	s.PayeeGuarantor.Name = PayeeGuarantorName
	return s
}

func (s *StubBoletoRequest) WithPayeeGuarantorDocumentNumber(docNumber string) *StubBoletoRequest {
	s.createStubPayeeGuarantor()
	s.PayeeGuarantor.Document.Number = docNumber
	return s
}

func (s *StubBoletoRequest) WithPayeeGuarantorDocumentType(documentType string) *StubBoletoRequest {
	s.createStubPayeeGuarantor()
	s.PayeeGuarantor.Document.Type = documentType
	return s
}

func (s *StubBoletoRequest) createStubPayeeGuarantor() {
	if !s.hasStubPayeeGuarantor() {
		s.PayeeGuarantor = &models.PayeeGuarantor{}
	}
}

func (s *StubBoletoRequest) hasStubPayeeGuarantor() bool {
	return s.PayeeGuarantor != nil
}

func (s *StubBoletoRequest) Build() *models.BoletoRequest {
	s.SetAuthentication(s.Authentication)
	s.SetAgreement(s.Agreement)
	s.SetTitle(s.Title)
	s.SetRecipient(s.Recipient)
	s.SetPayeeGuarantor(s.PayeeGuarantor)
	s.SetBuyer(s.Buyer)
	s.SetBank(s.bank)
	return s.BoletoRequest()
}
