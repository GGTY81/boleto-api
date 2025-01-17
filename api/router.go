package api

import (
	"github.com/gin-gonic/gin"
	"github.com/mundipagg/boleto-api/healthcheck"
)

func Base(router *gin.Engine) {
	router.StaticFile("/favicon.ico", "./boleto/favicon.ico")
	router.GET("/boleto", getBoletoLogger, getBoleto)
	router.GET("/boleto/memory-check/:unit", memory)
	router.GET("/boleto/memory-check/", memory)
	router.GET("/boleto/confirmation", confirmation)
	router.POST("/boleto/confirmation", confirmation)
	router.GET("/healthcheck", healthcheck.ExecuteOnAPI)
}

//V1 configura as rotas da v1
func V1(router *gin.Engine) {
	v1 := router.Group("v1")
	v1.Use(timingMetrics())
	v1.Use(returnHeaders())
	v1.POST("/boleto/register", authentication, parseBoleto, validateRegisterV1, registerBoletoLogger, errorResponseToClient, panicRecoveryHandler, registerBoleto)
	v1.GET("/boleto/:id", getBoletoByID)
}

//V2 configura as rotas da v2
func V2(router *gin.Engine) {
	v2 := router.Group("v2")
	v2.Use(timingMetrics())
	v2.Use(returnHeaders())
	v2.POST("/boleto/register", authentication, parseBoleto, validateRegisterV2, registerBoletoLogger, handleErrors, panicRecoveryHandler, registerBoleto)
}
