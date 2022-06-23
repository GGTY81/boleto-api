package bb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const requestExpected = `
 ## SOAPACTION:registrarBoleto
 ##	Authorization:Bearer {{.Authentication.AuthorizationToken}}
 ## Content-Type:text/xml; charset=utf-8

 <soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:sch="http://www.tibco.com/schemas/bws_registro_cbr/Recursos/XSD/Schema.xsd">
 <soapenv:Header/>
 <soapenv:Body>
<sch:requisicao>
 <sch:numeroConvenio>{{.Agreement.AgreementNumber}}</sch:numeroConvenio>
 <sch:numeroCarteira>17</sch:numeroCarteira>
 <sch:numeroVariacaoCarteira>{{.Agreement.WalletVariation}}</sch:numeroVariacaoCarteira>
 <sch:codigoModalidadeTitulo>1</sch:codigoModalidadeTitulo>
 <sch:dataEmissaoTitulo>{{replace (today | brdate) "/" "."}}</sch:dataEmissaoTitulo>
 <sch:dataVencimentoTitulo>{{replace (.Title.ExpireDateTime | brdate) "/" "."}}</sch:dataVencimentoTitulo>
 <sch:valorOriginalTitulo>{{toFloatStr .Title.AmountInCents}}</sch:valorOriginalTitulo>
 <sch:codigoTipoDesconto>0</sch:codigoTipoDesconto> 
 <sch:codigoTipoMulta>0</sch:codigoTipoMulta> 
 <sch:codigoAceiteTitulo>N</sch:codigoAceiteTitulo>
 <sch:codigoTipoTitulo>{{.Title.BoletoTypeCode}}</sch:codigoTipoTitulo>
 <sch:textoDescricaoTipoTitulo></sch:textoDescricaoTipoTitulo>
 <sch:indicadorPermissaoRecebimentoParcial>N</sch:indicadorPermissaoRecebimentoParcial>
 <sch:textoNumeroTituloBeneficiario>{{.Title.DocumentNumber}}</sch:textoNumeroTituloBeneficiario>
 <sch:textoNumeroTituloCliente>000{{padLeft (toString .Agreement.AgreementNumber) "0" 7}}{{padLeft (toString .Title.OurNumber) "0" 10}}</sch:textoNumeroTituloCliente>
 <sch:textoMensagemBloquetoOcorrencia>Pagamento disponível até a data de vencimento</sch:textoMensagemBloquetoOcorrencia>
 <sch:codigoTipoInscricaoPagador>{{docType .Buyer.Document}}</sch:codigoTipoInscricaoPagador>
 <sch:numeroInscricaoPagador>{{clearString (truncate .Buyer.Document.Number 15)}}</sch:numeroInscricaoPagador>
 <sch:nomePagador>{{clearString (truncate .Buyer.Name 60)}}</sch:nomePagador>
 <sch:textoEnderecoPagador>{{clearString (truncate .Buyer.Address.Street 60)}}</sch:textoEnderecoPagador>
 <sch:numeroCepPagador>{{extractNumbers .Buyer.Address.ZipCode}}</sch:numeroCepPagador>
 <sch:nomeMunicipioPagador>{{clearString (truncate .Buyer.Address.City 20)}}</sch:nomeMunicipioPagador>
 <sch:nomeBairroPagador>{{clearString (truncate .Buyer.Address.District 20)}}</sch:nomeBairroPagador>
 <sch:siglaUfPagador>{{clearString (truncate .Buyer.Address.StateCode 2)}}</sch:siglaUfPagador> 
 <sch:codigoChaveUsuario>1</sch:codigoChaveUsuario>
 <sch:codigoTipoCanalSolicitacao>5</sch:codigoTipoCanalSolicitacao>
 </sch:requisicao>
 </soapenv:Body>
</soapenv:Envelope>
 `

func Test_GivenTheGetRequestMethodWasCalled_ThenItShouldCorrectlyGetTheRequestTemplate(t *testing.T) {
	result := getRequest()

	assert.Equal(t, requestExpected, result, "Deve trazer corretamente o template de request")
}
