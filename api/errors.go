package api

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/mundipagg/boleto-api/models"
)

var validate = map[string]int{
	"MP400":                   http.StatusBadRequest,
	"MPAmountInCents":         http.StatusBadRequest,
	"MPExpireDate":            http.StatusBadRequest,
	"MPBuyerDocumentType":     http.StatusBadRequest,
	"MPDocumentNumber":        http.StatusBadRequest,
	"MPRecipientDocumentType": http.StatusBadRequest,
	"MPTimeout":               http.StatusGatewayTimeout,
	"MPOurNumberFail":         http.StatusBadGateway,
}

func handleErrors(c *gin.Context) {
	c.Next()

	var status int
	var exist bool

	response := getResponseFromContext(c)

	if !qualifiedForNewErrorHandling(c, response) {
		return
	}

	bankcode := response.Errors[0].Code
	if status, exist = validate[bankcode]; !exist {
		status = getBankFromContext(c).GetErrorsMap()[bankcode]
	}

	switch status {
	case http.StatusBadRequest:
		response.Errors[0].Code = "MP400"
		c.JSON(http.StatusBadRequest, response)
	case http.StatusBadGateway:
		response.Errors[0].Code = "MP502"
		c.JSON(http.StatusBadGateway, response)
	case http.StatusGatewayTimeout:
		response.Errors[0].Code = "MP504"
		c.JSON(http.StatusGatewayTimeout, response)
	default:
		response.Errors[0].Code = "MP500"
		clientResponse := getResponseError("MP500", "An internal error occurred.")
		c.JSON(http.StatusInternalServerError, clientResponse)
	}

	c.Set(responseKey, response)
}

func qualifiedForNewErrorHandling(c *gin.Context, response models.BoletoResponse) bool {
	if (getBankFromContext(c).GetErrorsMap() != nil && response.HasErrors()) || hasPanic(c) {
		return true
	}
	return false
}

func hasPanic(c *gin.Context) bool {
	_, exists := c.Get("hasPanic")

	return exists
}

func getErrorCodeToLog(c *gin.Context) string {
	response := getResponseFromContext(c)
	if response.HasErrors() {
		return response.Errors[0].ErrorCode()
	}
	return ""
}

func panicRecoveryHandler(c *gin.Context) {
	defer func() {
		if rec := recover(); rec != nil {
			err := fmt.Errorf("an internal error occurred: %v.\ninner exception: %s", rec, string(debug.Stack()))

			errorResponse := getResponseError("MP500", err.Error())

			c.Set(responseKey, errorResponse)
			c.Set("hasPanic", true)
		}
	}()

	c.Next()
}

func errorResponseToClient(c *gin.Context) {
	c.Next()

	if hasPanic(c) {
		errorResponse := getResponseError("MP500", "An internal error occurred.")

		c.JSON(http.StatusInternalServerError, errorResponse)
	}
}

func getResponseError(code string, message string) models.BoletoResponse {
	errorResponse := models.BoletoResponse{
		Errors: models.NewErrors(),
	}

	errorResponse.Errors.Append(code, message)

	return errorResponse
}
