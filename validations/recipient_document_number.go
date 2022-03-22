package validations

import "github.com/mundipagg/boleto-api/models"

//ValidateRecipientDocumentNumber Verifica se o número do documento do recebedor é válido
func ValidateRecipientDocumentNumber(b interface{}) error {
	switch t := b.(type) {
	case *models.BoletoRequest:
		if t.Recipient.Document.IsCPF() {
			return t.Recipient.Document.ValidateCPF()
		}
		if t.Recipient.Document.IsCNPJ() {
			return t.Recipient.Document.ValidateCNPJ()
		}
		return models.NewErrorResponse("MPRecipientDocumentType", "Tipo de Documento inválido")
	default:
		return InvalidType(t)
	}
}

//ValidatePayeeGuarantorDocumentNumber Verifica se o número do documento do lojista é válido
func ValidatePayeeGuarantorDocumentNumber(b interface{}) error {
	switch t := b.(type) {
	case *models.BoletoRequest:
		if t.HasPayeeGuarantor() {
			if t.PayeeGuarantor.Document.IsCPF() {
				return t.PayeeGuarantor.Document.ValidateCPF()
			}
			if t.PayeeGuarantor.Document.IsCNPJ() {
				return t.PayeeGuarantor.Document.ValidateCNPJ()
			}
			return models.NewErrorResponse("MPPayeeGuarantorDocumentType", "Tipo de Documento inválido")
		}
		return nil
	default:
		return InvalidType(t)
	}
}
