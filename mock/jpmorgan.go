package mock

import (
	"io/ioutil"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const successResponseJp = `
{
	"resposta":{
	   "clienteBeneficiario":{
		  "codBanco":376,
		  "codAgencia":98,
		  "codContaCorrente":6000164,
		  "tpPessoaBenfcrioOr":2,
		  "cnpjCpfBenfcrioOr":"33172537000198",
		  "codCartTit":1,
		  "txtInfCliCed":"TEST PAGARME"
	   },
	   "titulo":{
		  "numDocTit":"1001",
		  "dtVencTit":"2021-10-13",
		  "vlrTit":2,
		  "codEspTit":2,
		  "dtEmsTit":"2021-10-08",
		  "vlrAbattTit":0,
		  "tpCodPrott":3,
		  "qtdDiaPrott":0,
		  "codMoedaCnab":9
	   },
	   "juros":{
		  "codJurosTit":2,
		  "vlrPercJurosTit":7,
		  "dtJurosTit":"2021-10-14"
	   },
	   "descontos":{
		  "codDesctTit":0,
		  "dtDesctTit":null,
		  "vlrPercDesctTit":0
	   },
	   "sacadoOuPagador":{
		  "tpPessoaPagdr":1,
		  "cnpjCpfPagdr":"37303489819",
		  "nomRzSocPagdr":"Nome do Comprador Cliente",
		  "logradPagdr":"Logradouro",
		  "bairroPagdr":"Bairro",
		  "cepPagdr":"15050466",
		  "cidPagdr":"Cidade",
		  "ufPagdr":"SP"
	   },
	   "multa":{
		  "dtMultaTit":null,
		  "codMultaTit":3,
		  "vlrPercMultaTit":0
	   },
	   "numCodBarras":"37691877200000002000098000600016400000504858",
	   "linhaDigitavel":"37690098080060001640600005048582187720000000200",
	   "identdNossoNum":"123456789012"
	}
 }
`

const successResponseJpShortOurNumber = `
{
	"resposta":{
	   "clienteBeneficiario":{
		  "codBanco":376,
		  "codAgencia":98,
		  "codContaCorrente":6000164,
		  "tpPessoaBenfcrioOr":2,
		  "cnpjCpfBenfcrioOr":"33172537000198",
		  "codCartTit":1,
		  "txtInfCliCed":"TEST PAGARME"
	   },
	   "titulo":{
		  "numDocTit":"1001",
		  "dtVencTit":"2021-10-13",
		  "vlrTit":2,
		  "codEspTit":2,
		  "dtEmsTit":"2021-10-08",
		  "vlrAbattTit":0,
		  "tpCodPrott":3,
		  "qtdDiaPrott":0,
		  "codMoedaCnab":9
	   },
	   "juros":{
		  "codJurosTit":2,
		  "vlrPercJurosTit":7,
		  "dtJurosTit":"2021-10-14"
	   },
	   "descontos":{
		  "codDesctTit":0,
		  "dtDesctTit":null,
		  "vlrPercDesctTit":0
	   },
	   "sacadoOuPagador":{
		  "tpPessoaPagdr":1,
		  "cnpjCpfPagdr":"37303489819",
		  "nomRzSocPagdr":"Nome do Comprador Cliente",
		  "logradPagdr":"Logradouro",
		  "bairroPagdr":"Bairro",
		  "cepPagdr":"15050466",
		  "cidPagdr":"Cidade",
		  "ufPagdr":"SP"
	   },
	   "multa":{
		  "dtMultaTit":null,
		  "codMultaTit":3,
		  "vlrPercMultaTit":0
	   },
	   "numCodBarras":"37691877200000002000098000600016400000504858",
	   "linhaDigitavel":"37690098080060001640600005048582187720000000200",
	   "identdNossoNum":"123456"
	}
 }
`

const successResponseJpWithoutOurNumber = `
{
	"resposta":{
	   "clienteBeneficiario":{
		  "codBanco":376,
		  "codAgencia":98,
		  "codContaCorrente":6000164,
		  "tpPessoaBenfcrioOr":2,
		  "cnpjCpfBenfcrioOr":"33172537000198",
		  "codCartTit":1,
		  "txtInfCliCed":"TEST PAGARME"
	   },
	   "titulo":{
		  "numDocTit":"1001",
		  "dtVencTit":"2021-10-13",
		  "vlrTit":2,
		  "codEspTit":2,
		  "dtEmsTit":"2021-10-08",
		  "vlrAbattTit":0,
		  "tpCodPrott":3,
		  "qtdDiaPrott":0,
		  "codMoedaCnab":9
	   },
	   "juros":{
		  "codJurosTit":2,
		  "vlrPercJurosTit":7,
		  "dtJurosTit":"2021-10-14"
	   },
	   "descontos":{
		  "codDesctTit":0,
		  "dtDesctTit":null,
		  "vlrPercDesctTit":0
	   },
	   "sacadoOuPagador":{
		  "tpPessoaPagdr":1,
		  "cnpjCpfPagdr":"37303489819",
		  "nomRzSocPagdr":"Nome do Comprador Cliente",
		  "logradPagdr":"Logradouro",
		  "bairroPagdr":"Bairro",
		  "cepPagdr":"15050466",
		  "cidPagdr":"Cidade",
		  "ufPagdr":"SP"
	   },
	   "multa":{
		  "dtMultaTit":null,
		  "codMultaTit":3,
		  "vlrPercMultaTit":0
	   },
	   "numCodBarras":"37691877200000002000098000600016400000504858",
	   "linhaDigitavel":"37690098080060001640600005048582187720000000200",
	   "identdNossoNum":""
	}
 }
`

var withoutCert = `
<ServiceMessage>
    <Status>Error</Status>
    <Reason>Internal Error - Contact Service Provider</Reason>
</ServiceMessage>
`

const accountNotFound = `
{
    "errors" : [
        {
            "errorCode" : "GCA-010",
            "errorMsg" : "The account was not found."
        } ]
}
`

const missingParameter = `
{
    "errors": [
        {
            "errorCode": "GCA-111",
            "errorMsg": "codContaCorrente/codAgencia is required"
        }
    ]
}
`

const expirationDateBeforeIssue = `
{
    "numCodRetorno": "BOL-3",
    "mensagem": "Não foi possível processar a requisição por inconsistencia nos campos abaixo",
    "resposta": {
        "titulo": [
            {
                "nomCampo": "dtVencTit",
                "numCodRetorno": 108,
                "valorRecebido": "2021-10-17",
                "mensgRetorno": "Data de vencimento inválida"
            },
            {
                "nomCampo": "dtEmsTit",
                "numCodRetorno": 111,
                "valorRecebido": "2021-10-18",
                "mensgRetorno": "Data de emissão inválida"
            }
        ]
    }
}
`

const boletoDuplicated = `
{
    "numCodRetorno": "BOL-144",
    "mensagem": "Boleto já Existente. Detalhes abaixo:",
    "resposta": {
        "boletoBeneficiaryNumber": "4102374718",
        "boletoBarCode": "37691878200000002010098000600016400000504893",
        "boletoPaymentLine": "37690098080060001640600005048939187820000000201"
    }
}
`

const boletoDuplicatedWithoutDetails = `
{
    "numCodRetorno": "BOL-144",
    "mensagem": "Boleto já Existente. Detalhes abaixo:"
}
`

const signatureFailed = `
<ServiceMessage>
    <Status>Error</Status>
    <Reason>Signature Verification Failure </Reason>
</ServiceMessage>
`

const authenticationFailed = `
<ServiceMessage>
    <Status>Error</Status>
    <Reason>Authentication Failure. </Reason>
</ServiceMessage>
`

func registerJpMorgan(c *gin.Context) {
	d, _ := ioutil.ReadAll(c.Request.Body)
	json := string(d)

	if strings.Contains(json, `"vlrTit": 2,`) {
		c.Data(200, contentApplication, []byte(successResponseJp))
	} else if strings.Contains(json, `"vlrTit": 2.05,`) {
		c.Data(200, "application/xml", []byte(successResponseJpWithoutOurNumber))
	} else if strings.Contains(json, `"vlrTit": 2.11,`) {
		c.Data(200, "application/xml", []byte(successResponseJpShortOurNumber))
	} else if strings.Contains(json, `"vlrTit": 2.01,`) {
		c.Data(500, "application/xml", []byte(withoutCert))
	} else if strings.Contains(json, `"vlrTit": 2.02,`) {
		c.Data(404, "application/json", []byte(accountNotFound))
	} else if strings.Contains(json, `"vlrTit": 2.03,`) {
		c.Data(422, "application/json", []byte(expirationDateBeforeIssue))
	} else if strings.Contains(json, `"vlrTit": 2.04,`) {
		c.Data(409, "application/json", []byte(boletoDuplicated))
	} else if strings.Contains(json, `"vlrTit": 2.06,`) {
		c.Data(400, "application/json", []byte(missingParameter))
	} else if strings.Contains(json, `"vlrTit": 2.07,`) {
		c.Data(401, "application/xml", []byte(signatureFailed))
	} else if strings.Contains(json, `"vlrTit": 2.08,`) {
		c.Data(401, "application/xml", []byte(authenticationFailed))
	} else if strings.Contains(json, `"vlrTit": 2.09,`) {
		c.Data(409, "application/xml", []byte(boletoDuplicatedWithoutDetails))
	} else {
		time.Sleep(35 * time.Second)
		c.Data(504, contentApplication, []byte("timeout"))
	}
}
