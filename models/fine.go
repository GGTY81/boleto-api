package models

import "fmt"

const (
	defaultFineAmountInCents     = 0
	defaultFinePercentageOnTotal = 0.0
	minDaysToStartChargingFine   = 1
)

//Fine Representa as informações sobre Multa
type Fine struct {
	DaysAfterExpirationDate uint    `json:"daysAfterExpirationDate,omitempty"`
	AmountInCents           uint64  `json:"amountInCents,omitempty"`
	PercentageOnTotal       float64 `json:"percentageOnTotal,omitempty"`
}

//HasAmountInCents Verifica se há AmountInCents
func (fine *Fine) HasAmountInCents() bool {
	return fine.AmountInCents > defaultFineAmountInCents
}

//HasPercentageOnTotal Verifica se há PercentageOnTotal
func (fine *Fine) HasPercentageOnTotal() bool {
	return fine.PercentageOnTotal > defaultFinePercentageOnTotal
}

//HasDaysAfterExpirationDate Verifica se há DaysAfterExpirationDate
func (fine *Fine) HasDaysAfterExpirationDate() bool {
	return fine.DaysAfterExpirationDate >= minDaysToStartChargingFine
}

//HasExclusiveRateValues Verifica se foram informados os valores reverente a multa de forma exclusiva
func (fine *Fine) HasExclusiveRateValues() bool {
	return (fine.HasAmountInCents() || fine.HasPercentageOnTotal()) && !(fine.HasAmountInCents() && fine.HasPercentageOnTotal())
}

//HasFine Verifica algum dado de multa está preenchido
func (fine *Fine) HasFine() bool {
	return fine != nil
}

//Validate Valida as regras de negócio da struct Fine
//Caso haja alguma violação retorna o erro caso a regra violada, caso contrário retorna nulo
func (fine *Fine) Validate() error {
	if fine.HasFine() {
		if !fine.HasExclusiveRateValues() {
			return NewErrorResponse("MP400", "Para o campo Fine deve ser informado exclusivamente o parâmetro AmountInCents ou PercentageOnTotal maiores que zero")
		}
		if !fine.HasDaysAfterExpirationDate() {
			return NewErrorResponse("MP400", fmt.Sprintf("Para o campo Fine o parâmetro DaysAfterExpirationDate precisa ser no mínimo %d", minDaysToStartChargingFine))
		}
	}

	return nil
}
