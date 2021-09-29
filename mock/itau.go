package mock

import (
	"io/ioutil"
	"strings"

	"github.com/gin-gonic/gin"
)

func getTokenItau(c *gin.Context) {
	b, _ := ioutil.ReadAll(c.Request.Body)
	const tok = `{
		"access_token": "5f1cb9512fe587763ea33a3fb31e62cb",
		"expires_in": 14400,
		"token_type": "Bearer"
	}`
	if strings.Contains(string(b), `clientId=&`) {
		c.Data(500, "text/json", []byte(`{"Message":"An error has occurred."}`))
	} else {
		c.Data(200, "text/json", []byte(tok))
	}

}

func registerItau(c *gin.Context) {
	b, _ := ioutil.ReadAll(c.Request.Body)
	const resp = `{
		"beneficiario": {
			"codigo_banco_beneficiario": "341",
			"digito_verificador_banco_beneficiario": "7",
			"agencia_beneficiario": "0407",
			"conta_beneficiario": "55292",
			"digito_verificador_conta_beneficiario": "6",
			"cpf_cnpj_beneficiario": "00123456789012",
			"nome_razao_social_beneficiario": "NOME BENEFICIARIO",
			"logradouro_beneficiario": "RUA TESTE",
			"bairro_beneficiario": "TESTE",
			"complemento_beneficiario": "",
			"cidade_beneficiario": "RIO DE JANEIRO",
			"uf_beneficiario": "RJ",
			"cep_beneficiario": "22330000"
		},
		"pagador": {
			"cpf_cnpj_pagador": "00001234567890",
			"nome_razao_social_pagador": "NOME TESTE",
			"logradouro_pagador": "RUA TESTE",
			"complemento_pagador": "",
			"bairro_pagador": "BAIRRO TESTE",
			"cidade_pagador": "RIO DE JANEIRO",
			"uf_pagador": "RJ",
			"cep_pagador": "22555000"
		},
		"sacador_avalista": {
			"cpf_cnpj_sacador_avalista": "00000000000000",
			"nome_razao_social_sacador_avalista": ""
		},
		"moeda": {
			"sigla_moeda": "R$",
			"quantidade_moeda": 0,
			"cotacao_moeda": 0
		},
		"especie_documento": "DM",
		"vencimento_titulo": "2017-12-31",
		"tipo_carteira_titulo": "109",
		"nosso_numero": "079499759",
		"seu_numero": "000001234567890",
		"codigo_barras": "34199739000000010001090794997590407552926000",
		"numero_linha_digitavel": "34191090739499759040475529260004973900000001000",
		"local_pagamento": "ATE O VENCIMENTO PAGUE EM QUALQUER BANCO OU CORRESPONDENTE NAO BANCARIO. APOS O VENCIMENTO, ACESSE ITAU.COM.BR/BOLETOS E PAGUE EM QUALQUER BANCO OU CORRESPONDENTE NAO BANCARIO.",
		"data_processamento": "2017-10-26",
		"data_emissao": "2017-09-22",
		"uso_banco": "",
		"valor_titulo": 10,
		"valor_desconto": 0,
		"valor_outra_deducao": 0,
		"valor_juro_multa": 0,
		"valor_outro_acrescimo": 0,
		"valor_total_cobrado": 0,
		"lista_texto_informacao_cliente_beneficiario": [
			{
				"texto_informacao_cliente_beneficiario": ""
			},
			{
				"texto_informacao_cliente_beneficiario": ""
			},
			{
				"texto_informacao_cliente_beneficiario": ""
			},
			{
				"texto_informacao_cliente_beneficiario": ""
			},
			{
				"texto_informacao_cliente_beneficiario": ""
			},
			{
				"texto_informacao_cliente_beneficiario": ""
			},
			{
				"texto_informacao_cliente_beneficiario": ""
			},
			{
				"texto_informacao_cliente_beneficiario": ""
			},
			{
				"texto_informacao_cliente_beneficiario": ""
			}
		]
	}`
	if strings.Contains(string(b), `"valor_cobrado": "0000000000000200"`) {
		c.Data(200, "text/json", []byte(resp))
	} else if strings.Contains(string(b), `"valor_cobrado": "0000000000000500"`) {
		c.Data(500, "text/json", []byte("500"))
	} else if strings.Contains(string(b), `"valor_cobrado": "0000000000000300"`) {
		c.Data(300, "text/json", []byte("500"))
	} else {
		c.Data(400, "text/json", []byte(`
			{
				"codigo":"error_code_mock",
				"mensagem":"error mock message"	
			}
		`))
	}

}
