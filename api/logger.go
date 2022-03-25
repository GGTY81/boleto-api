package api

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mundipagg/boleto-api/log"
	"github.com/mundipagg/boleto-api/metrics"
	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/util"
)

const (
	resultGetBoletoKey = "result"
)

//registerBoletoLogger Middleware de log do request e response do registro da BoletoAPI
func registerBoletoLogger(c *gin.Context) {
	boleto := getBoletoFromContext(c)
	bank := getBankFromContext(c)

	l := loadBankLog(c)

	l.RequestApplication(boleto, c.Request.URL.RequestURI(), util.HeaderToMap(c.Request.Header))

	c.Next()

	l = loadBankLog(c)

	resp, _ := c.Get(responseKey)

	if hasPanic(c) {
		l.ResponseApplicationFatal(resp, c.Request.URL.RequestURI(), getErrorCodeToLog(c))
	} else {
		l.ResponseApplication(resp, c.Request.URL.RequestURI(), getErrorCodeToLog(c))
	}

	tag := bank.GetBankNameIntegration() + "-status"
	metrics.PushBusinessMetric(tag, c.Writer.Status())
}

//getBoletoLogger Middleware de log da operação de GetBoleto
func getBoletoLogger(c *gin.Context) {
	start := time.Now()
	c.Next()
	elapsedTimeInMilliseconds := time.Since(start).Milliseconds()

	result := getResultFromContext(c)
	result.TotalElapsedTimeInMilliseconds = elapsedTimeInMilliseconds

	log := log.CreateLog()
	log.Operation = "GetBoleto"
	log.IPAddress = c.ClientIP()
	log.RequestKey = getRequestKeyFromContext(c)

	log.GetBoleto(result, result.LogSeverity)
}

func loadBankLog(c *gin.Context) *log.Log {
	boleto := getBoletoFromContext(c)
	bank := getBankFromContext(c)
	l := bank.Log()
	l.Operation = "RegisterBoleto"
	l.NossoNumero = getNossoNumeroFromContext(c)
	l.Recipient = boleto.Recipient.Name
	if boleto.HasPayeeGuarantor() {
		l.PayeeGuarantor = boleto.PayeeGuarantor.Name
	}
	l.RequestKey = boleto.RequestKey
	l.BankName = bank.GetBankNameIntegration()
	l.IPAddress = c.ClientIP()
	l.ServiceUser = getUserFromContext(c)
	return l
}

func getResultFromContext(c *gin.Context) *models.GetBoletoResult {
	if result, exists := c.Get(resultGetBoletoKey); exists {
		return result.(*models.GetBoletoResult)
	}
	return nil
}

func getRequestKeyFromContext(c *gin.Context) string {
	var requestKey string
	requestKey = c.Request.Header.Get("RequestKey")
	if requestKey == "" {
		uid, _ := uuid.NewUUID()
		requestKey = uid.String()
	}

	return strings.ToLower(requestKey)
}
