package models

import "fmt"

const (
	defaultInterestAmountInCents      = 0
	defaultInterestPercentagePerMonth = 0.0
	minDaysToStartChargingInterest    = 1
)

//Interest Representa as informações sobre juros
type Interest struct {
	DaysAfterExpirationDate uint    `json:"daysAfterExpirationDate,omitempty"`
	AmountPerDayInCents     uint64  `json:"amountPerDayInCents,omitempty"`
	PercentagePerMonth      float64 `json:"percentagePerMonth,omitempty"`
}

//HasAmountPerDayInCents Verifica se há AmountPerDayInCents
func (interest *Interest) HasAmountPerDayInCents() bool {
	return interest.AmountPerDayInCents > defaultInterestAmountInCents
}

//HasPercentagePerMonth Verifica se há PercentagePerMonth
func (interest *Interest) HasPercentagePerMonth() bool {
	return interest.PercentagePerMonth > defaultInterestPercentagePerMonth
}

//HasDaysAfterExpirationDate Verifica se há DaysAfterExpirationDate
func (interest *Interest) HasDaysAfterExpirationDate() bool {
	return interest.DaysAfterExpirationDate >= minDaysToStartChargingInterest
}

//HasExclusiveRateValues Verifica se foi informados os valores reverente aos juros de forma exclusiva
func (interest *Interest) HasExclusiveRateValues() bool {
	return (interest.HasAmountPerDayInCents() || interest.HasPercentagePerMonth()) && !(interest.HasAmountPerDayInCents() && interest.HasPercentagePerMonth())
}

//HasInterest Verifica se algum dado de juros está preenchido
func (interest *Interest) HasInterest() bool {
	return interest != nil
}

//Validate Valida as regras de negócio da struct Interest
//Caso haja alguma violação retorna o erro caso a regra infrigida, caso contrário retorna nulo
func (interest *Interest) Validate() error {
	if interest.HasInterest() {
		if !interest.HasExclusiveRateValues() {
			return NewErrorResponse("MP400", "Para o campo Interest deve ser informado exclusivamente o parâmetro AmountInCents ou PercentagePerMonth maiores que zero")
		}
		if !interest.HasDaysAfterExpirationDate() {
			return NewErrorResponse("MP400", fmt.Sprintf("Para o campo Interest o parâmetro DaysAfterExpirationDate precisa ser no mínimo %d", minDaysToStartChargingInterest))
		}
	}

	return nil
}
