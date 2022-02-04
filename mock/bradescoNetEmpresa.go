package mock

import (
	"io/ioutil"
	"strings"

	"github.com/gin-gonic/gin"
)

func registerBoletoBradescoNetEmpresa(c *gin.Context) {

	const respOk = `<?xml version="1.0" encoding="UTF-8"?>
	<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/">
		<soapenv:Body>
			<ns2:registrarTituloResponse xmlns:ns2="http://ws.registrotitulo.ibpj.web.bradesco.com.br/">
				<return>
				{
					"cdErro": "0",
					"msgErro": "Solicita&amp;ccedil;&amp;atilde;o atendida",
					"idProduto": "9",
					"negociacao": "265600000000053937",
					"clubBanco": "2269651",
					"tpContrato": "48",
					"nuSequenciaContrato": "2078071",
					"cdProduto": "1730",
					"nuTituloGerado": "99999999999",
					"agenciaCreditoBeneficiario": "0",
					"contaCreditoBeneficiario": "0",
					"digCreditoBeneficiario": "00",
					"cdCipTitulo": "0",
					"statusTitulo": "1",
					"descStatusTitulo": "A VENCER/VENCIDO",
					"nomeBeneficiario": "NOME DO CEDENTE",
					"logradouroBeneficiario": "ENDERECO CEDENTE",
					"nuLogradouroBeneficiario": "",
					"complementoLogradouroBeneficiario": "",
					"bairroBeneficiario": "BAIRRO CEDENTE",
					"cepBeneficiario": "99999",
					"cepComplementoBeneficiario": "999",
					"municipioBeneficiario": "MUNICIPIO DO CEDENTE",
					"ufBeneficiario": "UF",
					"razaoContaBeneficiario": "0",
					"nomePagador": "Cliente Teste",
					"cpfcnpjPagador": "123322797000017",
					"enderecoPagador": "rua Teste",
					"bairroPagador": "bairro Teste",
					"municipioPagador": "Teste",
					"ufPagador": "RJ",
					"cepPagador": "21510",
					"cepComplementoPagador": "013",
					"endEletronicoPagador": "",
					"nomeSacadorAvalista": "",
					"cpfcnpjSacadorAvalista": "0",
					"enderecoSacadorAvalista": "",
					"municipioSacadorAvalista": "",
					"ufSacadorAvalista": "",
					"cepSacadorAvalista": "0",
					"cepComplementoSacadorAvalista": "0",
					"numeroTitulo": "1234567890",
					"dtRegistro": "14012018",
					"especieDocumentoTitulo": "DM",
					"descEspecie": "",
					"vlIOF": "0",
					"dtEmissao": "14012018",
					"dtVencimento": "30.06.2018",
					"vlTitulo": "100",
					"vlAbatimento": "0",
					"dtInstrucaoProtestoNegativacao": "",
					"diasInstrucaoProtestoNegativacao": "0",
					"dtMulta": "",
					"vlMulta": "0",
					"qtdeCasasDecimaisMulta": "0",
					"cdValorMulta": "0",
					"descCdMulta": "",
					"dtJuros": "",
					"vlJurosAoDia": "0",
					"dtDesconto1Bonificacao": "",
					"vlDesconto1Bonificacao": "0",
					"qtdeCasasDecimaisDesconto1Bonificacao": "0",
					"cdValorDesconto1Bonificacao": "0",
					"descCdDesconto1Bonificacao": "",
					"dtDesconto2": "",
					"vlDesconto2": "0",
					"qtdeCasasDecimaisDesconto2": "0",
					"cdValorDesconto2": "0",
					"descCdDesconto2": "",
					"dtDesconto3": "",
					"vlDesconto3": "0",
					"qtdeCasasDecimaisDesconto3": "0",
					"cdValorDesconto3": "0",
					"descCdDesconto3": "",
					"diasDispensaMulta": "0",
					"diasDispensaJuros": "0",	
					"cdBarras":"&lt;NWnnwnNnWwnwwNNnNwnWwwNNnnnWWnnnWWnnnWWnnwNNwnNwwNNwWnnnwWNnnwNWnnnWWnnnWWnNnwwNNWnnwNnWnwnnWWnWNwnnNWnwnnnNWw>",
					"linhaDigitavel": "23792.65602 90000.001231 45005.393702 6 74230000000200",
					"cdAcessorioEscrituralEmpresa": "0",
					"tpVencimento": "0",
					"indInstrucaoProtesto": "0",
					"tipoAbatimentoTitulo": "0",
					"cdValorJuros": "0",
					"tpDesconto1": "0",
					"tpDesconto2": "0",
					"tpDesconto3": "0",
					"nuControleParticipante": "",
					"diasJuros": "0",
					"cdJuros": "0",
					"vlJuros": "0",
					"cpfcnpjBeneficiario": "",
					"vlTituloEmitidoBoleto": "0",
					"dtVencimentoBoleto": "30.06.2018",
					"indTituloPertenceBaseTitulos": "",
					"dtLimitePagamentoBoleto": "30.06.2018",
					"cdIdentificacaoTituloDDACIP": "0",
					"indPagamentoParcial": "N",
					"qtdePagamentoParciais": "0"
				}
				</return>
			</ns2:registrarTituloResponse>
		</soapenv:Body>
	</soapenv:Envelope>
`

	const respError = `
<?xml version="1.0" encoding="UTF-8"?>
<soapenv:Envelope 
xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/">
<soapenv:Body>
	<ns2:registrarTituloResponse 
		xmlns:ns2="http://ws.registrotitulo.ibpj.web.bradesco.com.br/">
		<return>
		{
			"cdErro": "40",
			"msgErro": "N&amp;uacute;mero CGC/CPF inv&amp;aacute;lido",
			"idProduto": "0",
			"negociacao": "0",
			"clubBanco": "0",
			"tpContrato": "0",
			"nuSequenciaContrato": "0",
			"cdProduto": "0",
			"nuTituloGerado": "0",
			"agenciaCreditoBeneficiario": "0",
			"contaCreditoBeneficiario": "0",
			"digCreditoBeneficiario": "",
			"cdCipTitulo": "0",
			"statusTitulo": "0",
			"descStatusTitulo": "",
			"nomeBeneficiario": "",
			"logradouroBeneficiario": "",
			"nuLogradouroBeneficiario": "",
			"complementoLogradouroBeneficiario": "",
			"bairroBeneficiario": "",
			"cepBeneficiario": "0",
			"cepComplementoBeneficiario": "0",
			"municipioBeneficiario": "",
			"ufBeneficiario": "",
			"razaoContaBeneficiario": "0",
			"nomePagador": "",
			"cpfcnpjPagador": "0",
			"enderecoPagador": "",
			"bairroPagador": "",
			"municipioPagador": "",
			"ufPagador": "",
			"cepPagador": "0",
			"cepComplementoPagador": "",
			"endEletronicoPagador": "",
			"nomeSacadorAvalista": "",
			"cpfcnpjSacadorAvalista": "0",
			"enderecoSacadorAvalista": "",
			"municipioSacadorAvalista": "",
			"ufSacadorAvalista": "",
			"cepSacadorAvalista": "0",
			"cepComplementoSacadorAvalista": "0",
			"numeroTitulo": "",
			"dtRegistro": "",
			"especieDocumentoTitulo": "",
			"descEspecie": "",
			"vlIOF": "0",
			"dtEmissao": "",
			"dtVencimento": "",
			"vlTitulo": "0",
			"vlAbatimento": "0",
			"dtInstrucaoProtestoNegativacao": "",
			"diasInstrucaoProtestoNegativacao": "0",
			"dtMulta": "",
			"vlMulta": "0",
			"qtdeCasasDecimaisMulta": "0",
			"cdValorMulta": "0",
			"descCdMulta": "",
			"dtJuros": "",
			"vlJurosAoDia": "0",
			"dtDesconto1Bonificacao": "",
			"vlDesconto1Bonificacao": "0",
			"qtdeCasasDecimaisDesconto1Bonificacao": "0",
			"cdValorDesconto1Bonificacao": "0",
			"descCdDesconto1Bonificacao": "",
			"dtDesconto2": "",
			"vlDesconto2": "0",
			"qtdeCasasDecimaisDesconto2": "0",
			"cdValorDesconto2": "0",
			"descCdDesconto2": "",
			"dtDesconto3": "",
			"vlDesconto3": "0",
			"qtdeCasasDecimaisDesconto3": "0",
			"cdValorDesconto3": "0",
			"descCdDesconto3": "",
			"diasDispensaMulta": "0",
			"diasDispensaJuros": "0",
			"cdBarras": "",
			"linhaDigitavel": "",
			"cdAcessorioEscrituralEmpresa": "0",
			"tpVencimento": "0",
			"indInstrucaoProtesto": "0",
			"tipoAbatimentoTitulo": "0",
			"cdValorJuros": "0",
			"tpDesconto1": "0",
			"tpDesconto2": "0",
			"tpDesconto3": "0",
			"nuControleParticipante": "",
			"diasJuros": "0",
			"cdJuros": "0",
			"vlJuros": "0",
			"cpfcnpjBeneficiario": "",
			"vlTituloEmitidoBoleto": "0",
			"dtVencimentoBoleto": "",
			"indTituloPertenceBaseTitulos": "",
			"dtLimitePagamentoBoleto": "",
			"cdIdentificacaoTituloDDACIP": "0",
			"indPagamentoParcial": "",
			"qtdePagamentoParciais": "0"
		}           
		</return>
	</ns2:registrarTituloResponse>
</soapenv:Body>
</soapenv:Envelope>
`

	const respErrFail = `
<?xml version="1.0" encoding="UTF-8"?>
<S:Envelope 
	xmlns:S="http://schemas.xmlsoap.org/soap/envelope/" 
	xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/">
	<S:Body>
		<ns2:registrarTituloResponse 
			xmlns:ns2="http://ws.registrotitulo.ibpj.web.bradesco.com.br/">
			<return>{"cdErro":"810", "msgErro":"Erro Certificado / Formatacao dos campos da mensagem invalida [0x02430001]"}</return>
		</ns2:registrarTituloResponse>
	</S:Body>
</S:Envelope>
`

	const resErrSpecialCharacter = `<?xml version="1.0" encoding="UTF-8"?>
	<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/">
		 <soapenv:Body>
				<ns2:registrarTituloResponse xmlns:ns2="http://ws.registrotitulo.ibpj.web.bradesco.com.br/">
					 <return>{
						 "cdErro":"0",
						  "msgErro":"SOLICITACAO ATENDIDA (0)", 
							"idProduto":"9", 
							"negociacao":"364500000000002877", 
							"clubBanco":"2269651", 
							"tpContrato":"48", 
							"nuSequenciaContrato":"2091580", 
							"cdProduto":"1730", 
							"nuTituloGerado":"99999999999", 
							"agenciaCreditoBeneficiario":"0", 
							"contaCreditoBeneficiario":"0", 
							"digCreditoBeneficiario":"00", 
							"cdCipTitulo":"0", 
							"statusTitulo":"1", 
							"descStatusTitulo":"A VENCER/VENCIDO", 
							"nomeBeneficiario":"NOME DO CEDENTE",
							"logradouroBeneficiario":"ENDERECO CEDENTE", 
							"nuLogradouroBeneficiario":"", 
							"complementoLogradouroBeneficiario":"", 
							"bairroBeneficiario":"BAIRRO CEDENTE", 
							"cepBeneficiario":"99999", 
							"cepComplementoBeneficiario":"999", 
							"municipioBeneficiario":"MUNICIPIO DO CEDENTE", 
							"ufBeneficiario":"UF", 
							"razaoContaBeneficiario":"0", 
							"nomePagador":"Nome do Comprador Cliente", 
							"cpfcnpjPagador":"397340220000059", 
							"enderecoPagador":"Logradouro 	do Comprador", 
							"bairroPagador":"Bairro do Comprador",
							"municipioPagador":"Cidade do Comprador",
							"ufPagador":"SC", 
							"cepPagador":"1001", 
							"cepComplementoPagador":"000", 
							"endEletronicoPagador":"compradorgmail.com", 
							"nomeSacadorAvalista":"", 
							"cpfcnpjSacadorAvalista":"0", 
							"enderecoSacadorAvalista":"", 
							"municipioSacadorAvalista":"", 
							"ufSacadorAvalista":"", 
							"cepSacadorAvalista":"0", 
							"cepComplementoSacadorAvalista":"0", 
							"numeroTitulo":"00022072766", 
							"dtRegistro":"04022022", 
							"especieDocumentoTitulo":"DM", 
							"descEspecie":"", 
							"vlIOF":"0", 
							"dtEmissao":"04022022", 
							"dtVencimento":"09.02.2022", 
							"vlTitulo":"200", 
							"vlAbatimento":"0", 
							"dtInstrucaoProtestoNegativacao":"", 
							"diasInstrucaoProtestoNegativacao":"0", 
							"dtMulta":"", 
							"vlMulta":"0",
							"qtdeCasasDecimaisMulta":"0", 
							"cdValorMulta":"0", 
							"descCdMulta":"", 
							"dtJuros":"", 
							"vlJurosAoDia":"0", 
							"dtDesconto1Bonificacao":"",
							"vlDesconto1Bonificacao":"0", 
							"qtdeCasasDecimaisDesconto1Bonificacao":"0",
							"cdValorDesconto1Bonificacao":"0", 
							"descCdDesconto1Bonificacao":"", 
							"dtDesconto2":"", 
							"vlDesconto2":"0", 
							"qtdeCasasDecimaisDesconto2":"0", 
							"cdValorDesconto2":"0", 
							"descCdDesconto2":"", 
							"dtDesconto3":"", 
							"vlDesconto3":"0", 
							"qtdeCasasDecimaisDesconto3":"0", 
							"cdValorDesconto3":"0", 
							"descCdDesconto3":"", 
							"diasDispensaMulta":"0", 
							"diasDispensaJuros":"0", 
							"cdBarras":"&lt;WWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWW&gt;",
							"linhaDigitavel":"99999.99999 99999.999999 99999.999999 9 99999999999999",
							"cdAcessorioEscrituralEmpresa":"0",
							"tpVencimento":"0",
							"indInstrucaoProtesto":"0",
							"tipoAbatimentoTitulo":"0",
							"cdValorJuros":"0",
							"tpDesconto1":"0",
							"tpDesconto2":"0",
							"tpDesconto3":"0",
							"nuControleParticipante":"",
							"diasJuros":"0", 
							"cdJuros":"0", 
							"vlJuros":"0", 
							"cpfcnpjBeneficiario":"",
							"vlTituloEmitidoBoleto":"0",
							"dtVencimentoBoleto":"09.02.2022", 
							"indTituloPertenceBaseTitulos":"", 
							"dtLimitePagamentoBoleto":"09.02.2022",
							"cdIdentificacaoTituloDDACIP":"0", 
							"indPagamentoParcial":"", 
							"qtdePagamentoParciais":"0"
							}</return>
				</ns2:registrarTituloResponse>
		 </soapenv:Body>
	</soapenv:Envelope>`

	d, _ := ioutil.ReadAll(c.Request.Body)
	json := string(d)
	if strings.Contains(json, `"vlNominalTitulo": "200"`) {
		c.Data(200, "text/xml", []byte(respOk))
	} else if strings.Contains(json, `"vlNominalTitulo": "201"`) {
		c.Data(200, "text/xml", []byte(respError))
	} else if strings.Contains(json, `"vlNominalTitulo": "204"`) {
		c.Data(200, "text/xml", []byte(resErrSpecialCharacter))
	} else {
		c.Data(500, "text/xml", []byte(respErrFail))
	}
}
