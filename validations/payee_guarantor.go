package validations

import (
	"github.com/mundipagg/boleto-api/models"
)

//ValidatePayeeGuarantorName Verifica se o nome do lojista é existe
func ValidatePayeeGuarantorName(b interface{}) error {
	switch t := b.(type) {
	case *models.BoletoRequest:
		if t.HasPayeeGuarantor() {
			if !t.PayeeGuarantor.HasName() {
				return models.NewErrorResponse("MPPayeeGuarantorNameType", "Nome do sacador avalista está vazio")
			}
		}
		return nil
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
