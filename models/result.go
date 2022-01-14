package models

import (
	"regexp"

	"github.com/gin-gonic/gin"
)

//GetBoletoResult Centraliza as informações da operação GetBoleto
type GetBoletoResult struct {
	Id                                string
	Format                            string
	PrivateKey                        string
	URI                               string
	BoletoSource                      string
	TotalElapsedTimeInMilliseconds    int64
	CacheElapsedTimeInMilliseconds    int64
	DatabaseElapsedTimeInMilliseconds int64
	ErrorResponse                     BoletoResponse
	LogSeverity                       string
}

func NewGetBoletoResult(c *gin.Context) *GetBoletoResult {
	g := new(GetBoletoResult)
	g.Id = c.Query("id")
	g.Format = c.Query("fmt")
	g.PrivateKey = c.Query("pk")
	g.URI = c.Request.RequestURI
	g.BoletoSource = "none"
	return g
}

//HasValidPublicKey Verifica se a chave pública para buscar um boleto está presente e se é um hexadecimal
func HasValidPublicKey(g *GetBoletoResult) bool {
	return g.PrivateKey != "" && isValidHex(g.PrivateKey)
}

func HasValidId(g *GetBoletoResult) bool {
	return g.Id != "" && isValidHex(g.Id)
}

func (g *GetBoletoResult) HasValidParameters() bool {
	return HasValidPublicKey(g) && HasValidId(g)
}

//SetErrorResponse Insere as informações de erro para resposta
func (g *GetBoletoResult) SetErrorResponse(c *gin.Context, err ErrorResponse, statusCode int) {
	g.ErrorResponse = BoletoResponse{
		Errors: NewErrors(),
	}
	g.ErrorResponse.Errors.Append(err.Code, err.Message)

	if statusCode > 499 {
		c.JSON(statusCode, ErrorResponseToClient())
	} else {
		c.JSON(statusCode, g.ErrorResponse)
	}
}

func ErrorResponseToClient() BoletoResponse {
	resp := BoletoResponse{
		Errors: NewErrors(),
	}
	resp.Errors.Append("MP500", "Internal Error")
	return resp
}

func isValidHex(id string) bool {
	match, _ := regexp.MatchString("^([0-9A-Fa-f])+$", id)
	return match
}
