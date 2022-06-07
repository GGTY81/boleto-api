package validations

import (
	"github.com/mundipagg/boleto-api/models"
)

func ValidateInterest(b interface{}) error {
	switch t := b.(type) {
	case *models.BoletoRequest:
		if t.Title.Fees.HasInterest() {
			if err := t.Title.Fees.Interest.Validate(); err != nil {
				return err
			}
		}
		return nil
	default:
		return InvalidType(t)
	}
}

func ValidateFine(b interface{}) error {
	switch t := b.(type) {
	case *models.BoletoRequest:
		if t.Title.Fees.HasFine() {
			if err := t.Title.Fees.Fine.Validate(); err != nil {
				return err
			}
		}
		return nil
	default:
		return InvalidType(t)
	}
}
