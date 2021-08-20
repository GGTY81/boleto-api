package api

import (
	"github.com/gin-gonic/gin"
	"github.com/mundipagg/boleto-api/env"
)

const mockPanicRegistrationResponseJSON = `{"errors":[{"code":"MP500","message":"An internal error occurred."}]}`

const mockPanicRegistrationRequestJSON = `{"bankNumber":174,"authentication":{"Username":"altsa","Password":"altsa"},"agreement":{"agreementNumber":267,"wallet":36,"agency":"00000"},"title":{"expireDate":"2050-12-30","amountInCents":200,"ourNumber":1,"instructions":"Nãoreceberapósadatadevencimento.","documentNumber":"1234567890"},"recipient":{"name":"Empresa-Boletos","document":{"type":"CNPJ","number":"29799428000128"},"address":{"street":"AvenidaMiguelEstefno,2394","complement":"ÁguaFunda","zipCode":"04301-002","city":"SãoPaulo","stateCode":"SP"}},"buyer":{"name":"UsuarioTeste","email":"p@p.com","document":{"type":"CNPJ","number":"29.799.428/0001-28"},"address":{"street":"RuaTeste","number":"2","complement":"SALA1","zipCode":"20931-001","district":"Centro","city":"RiodeJaneiro","stateCode":"RJ"}}}`

func mockInstallApi() *gin.Engine {
	env.Config(true, true, true)
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	Base(r)
	V1(r)
	V2(r)
	return r
}

func mockPanicRegistration(c *gin.Context) {
	panic("A panic occurred.")
}
