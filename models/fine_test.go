package models

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testFineParameter struct {
	Line     interface{}
	Input    interface{}
	Expected interface{}
}

var hasFineParameters = []testFineParameter{
	{Line: 4, Input: Fine{DaysAfterExpirationDate: 1, AmountInCents: 1, PercentageOnTotal: 0}, Expected: true},
	{Line: 5, Input: Fine{DaysAfterExpirationDate: 1, AmountInCents: 0, PercentageOnTotal: 1.20}, Expected: true},
	{Line: 6, Input: Fine{DaysAfterExpirationDate: 1, AmountInCents: 200, PercentageOnTotal: 3}, Expected: true},
}

var hasAmountInCentsParameters = []testFineParameter{
	{Line: 1, Input: Fine{AmountInCents: 0}, Expected: false},
	{Line: 2, Input: Fine{AmountInCents: 1}, Expected: true},
	{Line: 3, Input: Fine{AmountInCents: 20034}, Expected: true},
}

var validateFineParameters = []testFineParameter{
	{Line: 1, Input: Fine{DaysAfterExpirationDate: 1, AmountInCents: 0, PercentageOnTotal: 0}, Expected: false},
	{Line: 2, Input: Fine{DaysAfterExpirationDate: 1, AmountInCents: 1, PercentageOnTotal: 0}, Expected: true},
	{Line: 3, Input: Fine{DaysAfterExpirationDate: 1, AmountInCents: 0, PercentageOnTotal: 1}, Expected: true},
	{Line: 4, Input: Fine{DaysAfterExpirationDate: 1, AmountInCents: 1, PercentageOnTotal: 1}, Expected: false},
	{Line: 5, Input: Fine{DaysAfterExpirationDate: 1, PercentageOnTotal: -1}, Expected: false},
	{Line: 6, Input: Fine{DaysAfterExpirationDate: 1, PercentageOnTotal: -2.34}, Expected: false},
	{Line: 7, Input: Fine{DaysAfterExpirationDate: 1, PercentageOnTotal: 0}, Expected: false},
	{Line: 8, Input: Fine{DaysAfterExpirationDate: 1, PercentageOnTotal: 0.0}, Expected: false},
	{Line: 9, Input: Fine{DaysAfterExpirationDate: 1, PercentageOnTotal: 1}, Expected: true},
	{Line: 10, Input: Fine{DaysAfterExpirationDate: 1, PercentageOnTotal: 2.23}, Expected: true},
	{Line: 11, Input: Fine{DaysAfterExpirationDate: 0, AmountInCents: 0, PercentageOnTotal: 0}, Expected: false},
	{Line: 12, Input: Fine{DaysAfterExpirationDate: 0, AmountInCents: 1, PercentageOnTotal: 0}, Expected: false},
	{Line: 13, Input: Fine{DaysAfterExpirationDate: 0, AmountInCents: 0, PercentageOnTotal: 1.2}, Expected: false},
	{Line: 14, Input: Fine{DaysAfterExpirationDate: 1, AmountInCents: 1, PercentageOnTotal: 0.0}, Expected: true},
	{Line: 15, Input: Fine{DaysAfterExpirationDate: 1, AmountInCents: 0, PercentageOnTotal: 1.2}, Expected: true},
	{Line: 16, Input: Fine{DaysAfterExpirationDate: 2, AmountInCents: 0, PercentageOnTotal: 1.2}, Expected: true},
}

var hasExclusiveRateValuesFineParameters = []testFineParameter{
	{Line: 1, Input: Fine{AmountInCents: 0, PercentageOnTotal: 0}, Expected: false},
	{Line: 2, Input: Fine{AmountInCents: 1, PercentageOnTotal: 0}, Expected: true},
	{Line: 3, Input: Fine{AmountInCents: 0, PercentageOnTotal: 1}, Expected: true},
	{Line: 4, Input: Fine{AmountInCents: 1, PercentageOnTotal: 1}, Expected: false},
}

var hasPercentageOnTotalParameters = []testFineParameter{
	{Line: 1, Input: Fine{PercentageOnTotal: -1}, Expected: false},
	{Line: 2, Input: Fine{PercentageOnTotal: -2.34}, Expected: false},
	{Line: 3, Input: Fine{PercentageOnTotal: 0}, Expected: false},
	{Line: 4, Input: Fine{PercentageOnTotal: 0.0}, Expected: false},
	{Line: 5, Input: Fine{PercentageOnTotal: 1}, Expected: true},
	{Line: 6, Input: Fine{PercentageOnTotal: 2.23}, Expected: true},
}

var hasDaysAfterExpirationDateFineParameters = []testFineParameter{
	{Line: 1, Input: Fine{DaysAfterExpirationDate: 0, AmountInCents: 0, PercentageOnTotal: 0}, Expected: false},
	{Line: 2, Input: Fine{DaysAfterExpirationDate: 0, AmountInCents: 1, PercentageOnTotal: 0}, Expected: false},
	{Line: 3, Input: Fine{DaysAfterExpirationDate: 0, AmountInCents: 0, PercentageOnTotal: 1.2}, Expected: false},
	{Line: 4, Input: Fine{DaysAfterExpirationDate: 1, AmountInCents: 1, PercentageOnTotal: 0.0}, Expected: true},
	{Line: 5, Input: Fine{DaysAfterExpirationDate: 1, AmountInCents: 0, PercentageOnTotal: 1.2}, Expected: true},
	{Line: 6, Input: Fine{DaysAfterExpirationDate: 2, AmountInCents: 0, PercentageOnTotal: 1.2}, Expected: true},
}

func TestHasFine(t *testing.T) {
	for _, fact := range hasFineParameters {
		input := fact.Input.(Fine)

		result := input.HasFine()

		assert.Equal(t, fact.Expected, result, fmt.Sprintf("HasFine - Linha %d: Deve verificar corretamente se possui multa", fact.Line.(int)))
	}
}

func TestHasAmountInCents(t *testing.T) {
	for _, fact := range hasAmountInCentsParameters {
		input := fact.Input.(Fine)

		result := input.HasAmountInCents()

		assert.Equal(t, fact.Expected, result, fmt.Sprintf("HasAmountInCents - Linha %d: Deve verificar corretamente se possui multa por valor", fact.Line.(int)))
	}
}

func TestFineValidate(t *testing.T) {
	for _, fact := range validateFineParameters {
		input := fact.Input.(Fine)

		result := input.Validate() == nil

		assert.Equal(t, fact.Expected, result, fmt.Sprintf("Validate - Linha %d: Deve validar a struct de multa", fact.Line.(int)))
	}
}

func TestFineHasExclusiveRateValues(t *testing.T) {
	for _, fact := range hasExclusiveRateValuesFineParameters {
		input := fact.Input.(Fine)

		result := input.HasExclusiveRateValues()

		assert.Equal(t, fact.Expected, result, fmt.Sprintf("HasExclusiveRateValues - Linha %d: Deve validar corretamente as regras de valores", fact.Line.(int)))
	}
}

func TestHasPercentageOnTotal(t *testing.T) {
	for _, fact := range hasPercentageOnTotalParameters {
		input := fact.Input.(Fine)

		result := input.HasPercentageOnTotal()

		assert.Equal(t, fact.Expected, result, fmt.Sprintf("HasPercentageOnTotal - Linha %d: Deve validar corretamente o campo Percentage", fact.Line.(int)))
	}
}

func TestFineHasDaysAfterExpirationDate(t *testing.T) {
	for _, fact := range hasDaysAfterExpirationDateFineParameters {
		input := fact.Input.(Fine)

		result := input.HasDaysAfterExpirationDate()

		assert.Equal(t, fact.Expected, result, fmt.Sprintf("HasDaysAfterExpirationDate - Linha %d: Deve validar corretamente o campo Days", fact.Line.(int)))
	}
}
