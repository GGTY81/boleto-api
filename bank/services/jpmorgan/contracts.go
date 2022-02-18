package jpmorgan

type arrayDataError struct {
	Error []errorResponse `json:"errors"`
}

type errorResponse struct {
	Code    string `json:"errorCode,omitempty"`
	Message string `json:"errorMsg,omitempty"`
}

type ServiceMessageError struct {
	Status string `xml:"Status"`
	Reason string `xml:"Reason"`
}

const templateRequest = `
{
	"clienteBeneficiario": {
        "cnpjCpfBenfcrioOr": {{toUint .Recipient.Document.Number}},
        "codAgencia": {{toUint .Agreement.Agency}},
        "codBanco": {{.BankNumber}}, 
        "codCartTit": 1,
        "codContaCorrente": {{toUint .Agreement.Account}},
        "tpPessoaBenfcrioOr": {{docType .Recipient.Document}},
        "txtInfCliCed": "{{truncateOnly (.Recipient.Name) 80}}"
    },
	"sacadoOuPagador": {
        "bairroPagdr": "{{truncateOnly (onlyAlphanumerics (onlyOneSpace .Buyer.Address.District)) 15}}",
        "cepPagdr": {{toUint .Buyer.Address.ZipCode}},
		"cidPagdr": "{{truncateOnly (onlyAlphanumerics (onlyOneSpace .Buyer.Address.City)) 15}}",
		"cnpjCpfPagdr": {{toUint .Buyer.Document.Number}},
		"logradPagdr": "{{truncateOnly (onlyAlphanumerics (onlyOneSpace (joinSpace (.Buyer.Address.Street) (.Buyer.Address.Number) (.Buyer.Address.Complement)))) 40}}",
		"nomRzSocPagdr": "{{truncateOnly (onlyAlphanumerics (onlyOneSpace .Buyer.Name)) 40}}",
		"tpPessoaPagdr": {{docType .Buyer.Document}},
		"ufPagdr": "{{truncateOnly (onlyAlphabetics (removeAllSpaces .Buyer.Address.StateCode)) 2}}"
    },
	"titulo": {
		"numDocTit": "{{.Title.DocumentNumber}}",
		"dtVencTit": "{{.Title.ExpireDate}}",
		"vlrTit": {{strToFloat (toFloatStr .Title.AmountInCents)}},
		"codEspTit": {{toUint .Title.BoletoTypeCode}},
		"dtEmsTit":  "{{enDate (today) "-"}}",
		"vlrAbattTit": 0,
		"tpCodPrott": 3,
		"qtdDiaPrott": 0,
        "codMoedaCnab": 9
	},
    "juros": {
        "codJurosTit": 3
    },
    "descontos": {
        "codDesctTit": 0
    }
}`

const templateAPI = `
{
    {{if (hasErrorTags . "errorCode") | (hasErrorTags . "messageError")}}
    "Errors": [
        {
        {{if (hasErrorTags . "errorCode")}}
            "Code": "{{trim .errorCode}}",
        {{end}}
        {{if (eq .messageError "{}") | (eq .messageError "Error")}}
            "Message": "{{trim .errorCode}}"
        {{else}}
            {{if (eq .messageDetails "{}")}}
                "Message": "{{joinSpace (trim .messageError)}}"
            {{else}}
                "Message": "{{joinSpace (trim .messageError) ":" (trim .messageDetails)}}"
            {{end}}
        {{end}}
        }
    ]
    {{else}}
        "DigitableLine": "{{fmtDigitableLine (trim .digitableLine)}}",
        "BarCodeNumber": "{{trim .barCodeNumber}}",
        "OurNumber": "{{padLeft .ourNumber "0" 12}}"
    {{end}}
}
`

const templateResponse = `
{
	"resposta":{
		"numCodBarras": "{{barCodeNumber}}",
		"identdNossoNum": "{{ourNumber}}",
		"linhaDigitavel": "{{digitableLine}}"
	}
}
`

var templateErrorXmltoJson = `
{
	"Status" : "{{messageError}}",
	"Reason" : "{{errorCode}}"
}
`

const templateErrorJson = `
{
	"errorCode": "{{errorCode}}",
	"errorMsg": "{{messageError}}",
    "resposta": "{{messageDetails}}"
}
`

const templateErrorBoletoJson = `
{
    "numCodRetorno": "{{errorCode}}",
    "mensagem": "{{messageError}}",
    "resposta": "{{messageDetails}}"
}
`
