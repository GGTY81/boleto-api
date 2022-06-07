package api

import (
	"github.com/gin-gonic/gin"
	"github.com/mundipagg/boleto-api/models"
)

//validateRegisterV1 Middleware de validação das requisições de registro de boleto na rota V1
func validateRegisterV1(c *gin.Context) {
	rules := getBoletoFromContext(c).Title.Rules
	bn := getBankFromContext(c).GetBankNumber()

	errorResponse := models.BoletoResponse{
		Errors: models.NewErrors(),
	}

	if rules != nil {
		errorResponse.Errors.Append("MP400", "title.rules not available in this version")
		c.AbortWithStatusJSON(400, errorResponse)
		return
	}

	if bn == models.Stone {
		errorResponse.Errors.Append("MP400", "bank Stone not available in this version")
		c.AbortWithStatusJSON(400, errorResponse)
		return
	}
}

//validateRegisterV2 Middleware de validação das requisições de registro de boleto na rota V2
func validateRegisterV2(c *gin.Context) {
	t := getBoletoFromContext(c).Title
	bn := getBankFromContext(c).GetBankNumber()

	errorResponse := models.BoletoResponse{
		Errors: models.NewErrors(),
	}

	if t.HasFees() && !isBankNumberAcceptFees(bn) {
		errorResponse.Errors.Append("MP400", "title.fees not available for this bank")
		c.AbortWithStatusJSON(400, errorResponse)
		return
	}

	if t.HasRules() && !isBankNumberAcceptRules(bn) {
		errorResponse.Errors.Append("MP400", "title.rules not available for this bank")
		c.AbortWithStatusJSON(400, errorResponse)
		return
	}
}

func isBankNumberAcceptFees(b models.BankNumber) bool {
	return b == models.Caixa || b == models.Stone
}

func isBankNumberAcceptRules(b models.BankNumber) bool {
	return b == models.Caixa || b == models.Stone
}
