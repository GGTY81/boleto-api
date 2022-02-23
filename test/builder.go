package test

import (
	"github.com/google/uuid"
	"github.com/mundipagg/boleto-api/models"
)

type BuilderBoletoRequest struct {
	authentication models.Authentication
	agreement      models.Agreement
	title          models.Title
	recipient      models.Recipient
	PayeeGuarantor *models.PayeeGuarantor
	buyer          models.Buyer
	bank           models.BankNumber
}

func NewBuilderBoletoRequest() BuilderBoletoRequest {
	return BuilderBoletoRequest{}
}

func (b *BuilderBoletoRequest) SetBank(bank models.BankNumber) {
	b.bank = bank
}

func (b *BuilderBoletoRequest) SetAuthentication(authentication models.Authentication) {
	b.authentication = authentication
}

func (b *BuilderBoletoRequest) SetAgreement(agreement models.Agreement) {
	b.agreement = agreement
}

func (b *BuilderBoletoRequest) SetTitle(title models.Title) {
	b.title = title
}

func (b *BuilderBoletoRequest) SetRecipient(recipient models.Recipient) {
	b.recipient = recipient
}

func (b *BuilderBoletoRequest) SetPayeeGuarantor(PayeeGuarantor *models.PayeeGuarantor) {
	b.PayeeGuarantor = PayeeGuarantor
}

func (b *BuilderBoletoRequest) SetBuyer(buyer models.Buyer) {
	b.buyer = buyer
}

func (b *BuilderBoletoRequest) BoletoRequest() *models.BoletoRequest {
	guid, _ := uuid.NewUUID()
	return &models.BoletoRequest{
		BankNumber:     b.bank,
		Authentication: b.authentication,
		Agreement:      b.agreement,
		Title:          b.title,
		Recipient:      b.recipient,
		PayeeGuarantor: b.PayeeGuarantor,
		Buyer:          b.buyer,
		RequestKey:     guid.String(),
	}
}
