package models

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testInterestParameter struct {
	Line     interface{}
	Input    interface{}
	Expected interface{}
}

var hasInterestParameters = []testInterestParameter{
	{Line: 4, Input: Interest{DaysAfterExpirationDate: 1, AmountPerDayInCents: 1, PercentagePerMonth: 0}, Expected: true},
	{Line: 5, Input: Interest{DaysAfterExpirationDate: 1, AmountPerDayInCents: 0, PercentagePerMonth: 1.20}, Expected: true},
	{Line: 6, Input: Interest{DaysAfterExpirationDate: 1, AmountPerDayInCents: 200, PercentagePerMonth: 3}, Expected: true},
}

var hasAmountPerDayInCentsParameters = []testInterestParameter{
	{Line: 1, Input: Interest{AmountPerDayInCents: 0}, Expected: false},
	{Line: 2, Input: Interest{AmountPerDayInCents: 1}, Expected: true},
	{Line: 3, Input: Interest{AmountPerDayInCents: 20034}, Expected: true},
}

var validateParameters = []testFineParameter{
	{Line: 1, Input: Interest{DaysAfterExpirationDate: 1, AmountPerDayInCents: 0, PercentagePerMonth: 0}, Expected: false},
	{Line: 2, Input: Interest{DaysAfterExpirationDate: 1, AmountPerDayInCents: 1, PercentagePerMonth: 0}, Expected: true},
	{Line: 3, Input: Interest{DaysAfterExpirationDate: 1, AmountPerDayInCents: 0, PercentagePerMonth: 1}, Expected: true},
	{Line: 4, Input: Interest{DaysAfterExpirationDate: 1, AmountPerDayInCents: 1, PercentagePerMonth: 1}, Expected: false},
	{Line: 5, Input: Interest{DaysAfterExpirationDate: 1, PercentagePerMonth: -1}, Expected: false},
	{Line: 6, Input: Interest{DaysAfterExpirationDate: 1, PercentagePerMonth: -2.34}, Expected: false},
	{Line: 7, Input: Interest{DaysAfterExpirationDate: 1, PercentagePerMonth: 0}, Expected: false},
	{Line: 8, Input: Interest{DaysAfterExpirationDate: 1, PercentagePerMonth: 0.0}, Expected: false},
	{Line: 9, Input: Interest{DaysAfterExpirationDate: 1, PercentagePerMonth: 1}, Expected: true},
	{Line: 10, Input: Interest{DaysAfterExpirationDate: 1, PercentagePerMonth: 2.23}, Expected: true},
	{Line: 11, Input: Interest{DaysAfterExpirationDate: 0, AmountPerDayInCents: 0, PercentagePerMonth: 0}, Expected: false},
	{Line: 12, Input: Interest{DaysAfterExpirationDate: 0, AmountPerDayInCents: 1, PercentagePerMonth: 0}, Expected: false},
	{Line: 13, Input: Interest{DaysAfterExpirationDate: 0, AmountPerDayInCents: 0, PercentagePerMonth: 1.2}, Expected: false},
	{Line: 14, Input: Interest{DaysAfterExpirationDate: 1, AmountPerDayInCents: 1, PercentagePerMonth: 0.0}, Expected: true},
	{Line: 15, Input: Interest{DaysAfterExpirationDate: 1, AmountPerDayInCents: 0, PercentagePerMonth: 1.2}, Expected: true},
	{Line: 16, Input: Interest{DaysAfterExpirationDate: 2, AmountPerDayInCents: 0, PercentagePerMonth: 1.2}, Expected: true},
}

var hasExclusiveRateValuesParameters = []testInterestParameter{
	{Line: 1, Input: Interest{AmountPerDayInCents: 0, PercentagePerMonth: 0}, Expected: false},
	{Line: 2, Input: Interest{AmountPerDayInCents: 1, PercentagePerMonth: 0}, Expected: true},
	{Line: 3, Input: Interest{AmountPerDayInCents: 0, PercentagePerMonth: 1}, Expected: true},
	{Line: 4, Input: Interest{AmountPerDayInCents: 1, PercentagePerMonth: 1}, Expected: false},
}

var hasPercentagePerMonthParameters = []testInterestParameter{
	{Line: 1, Input: Interest{PercentagePerMonth: -1}, Expected: false},
	{Line: 2, Input: Interest{PercentagePerMonth: -2.34}, Expected: false},
	{Line: 3, Input: Interest{PercentagePerMonth: 0}, Expected: false},
	{Line: 4, Input: Interest{PercentagePerMonth: 0.0}, Expected: false},
	{Line: 5, Input: Interest{PercentagePerMonth: 1}, Expected: true},
	{Line: 6, Input: Interest{PercentagePerMonth: 2.23}, Expected: true},
}

var hasDaysAfterExpirationDateParameters = []testInterestParameter{
	{Line: 1, Input: Interest{DaysAfterExpirationDate: 0, AmountPerDayInCents: 0, PercentagePerMonth: 0}, Expected: false},
	{Line: 2, Input: Interest{DaysAfterExpirationDate: 0, AmountPerDayInCents: 1, PercentagePerMonth: 0}, Expected: false},
	{Line: 3, Input: Interest{DaysAfterExpirationDate: 0, AmountPerDayInCents: 0, PercentagePerMonth: 1.2}, Expected: false},
	{Line: 4, Input: Interest{DaysAfterExpirationDate: 1, AmountPerDayInCents: 1, PercentagePerMonth: 0.0}, Expected: true},
	{Line: 5, Input: Interest{DaysAfterExpirationDate: 1, AmountPerDayInCents: 0, PercentagePerMonth: 1.2}, Expected: true},
	{Line: 6, Input: Interest{DaysAfterExpirationDate: 2, AmountPerDayInCents: 0, PercentagePerMonth: 1.2}, Expected: true},
}

func TestHasInterest(t *testing.T) {
	for _, fact := range hasInterestParameters {
		input := fact.Input.(Interest)

		result := input.HasInterest()

		assert.Equal(t, fact.Expected, result, fmt.Sprintf("HasInterest - Linha %d: Os juros não foram validados corretamente", fact.Line.(int)))
	}
}

func TestHasAmountPerDayInCents(t *testing.T) {
	for _, fact := range hasAmountPerDayInCentsParameters {
		input := fact.Input.(Interest)

		result := input.HasAmountPerDayInCents()

		assert.Equal(t, fact.Expected, result, fmt.Sprintf("HasAmountPerDayInCents - Linha %d: Os juros não foram validados corretamente", fact.Line.(int)))
	}
}

func TestValidate(t *testing.T) {
	for _, fact := range validateParameters {
		input := fact.Input.(Interest)

		result := input.Validate() == nil

		assert.Equal(t, fact.Expected, result, fmt.Sprintf("Validate - Linha %d: Deve validar a struct de juros", fact.Line.(int)))
	}
}

func TestHasExclusiveRateValues(t *testing.T) {
	for _, fact := range hasExclusiveRateValuesParameters {
		input := fact.Input.(Interest)

		result := input.HasExclusiveRateValues()

		assert.Equal(t, fact.Expected, result, fmt.Sprintf("HasExclusiveRateValues - Linha %d: Deve validar corretamente as regras de valores", fact.Line.(int)))
	}
}

func TestHasPercentagePerMonth(t *testing.T) {
	for _, fact := range hasPercentagePerMonthParameters {
		input := fact.Input.(Interest)

		result := input.HasPercentagePerMonth()

		assert.Equal(t, fact.Expected, result, fmt.Sprintf("HasPercentagePerMonth - Linha %d: Deve validar corretamente o campo Percentage", fact.Line.(int)))
	}
}

func TestHasDaysAfterExpirationDate(t *testing.T) {
	for _, fact := range hasDaysAfterExpirationDateParameters {
		input := fact.Input.(Interest)

		result := input.HasDaysAfterExpirationDate()

		assert.Equal(t, fact.Expected, result, fmt.Sprintf("HasDaysAfterExpirationDate - Linha %d: Deve validar corretamente o campo Days", fact.Line.(int)))
	}
}
