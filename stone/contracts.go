package stone

const templateRequest = `
## Content-Type:application/json
## Authorization:Bearer {{.Authentication.AuthorizationToken}}
{
    "account_id": "{{.Authentication.AccessKey}}",
    "amount": {{.Title.AmountInCents}},
    "expiration_date": "{{.Title.ExpireDate}}",
    {{if .Title.HasRules}} 
        "limit_date": "{{enDate (datePlusDays .Title.ExpireDateTime .Title.Rules.MaxDaysToPayPastDue) "-"}}",
    {{else}}
        "limit_date": "{{enDate (datePlusDays .Title.ExpireDateTime 60) "-"}}",
    {{end}}
    "invoice_type": "{{.Title.BoletoTypeCode}}",
    "customer": {
        "document": "{{.Buyer.Document.Number}}",
        "legal_name": "{{onlyOneSpace .Buyer.Name}}",
	{{if eq .Buyer.Document.Type "CNPJ"}}
        "trade_name": "{{onlyOneSpace .Buyer.Name}}"
	{{else}}
		"trade_name": null
	{{end}}
    }
    {{if .Title.Fees.HasFine}}
        ,"fine": {
            "date": "{{enDate (datePlusDays .Title.ExpireDateTime .Title.Fees.Fine.DaysAfterExpirationDate) "-"}}",
            {{if .Title.Fees.Fine.HasAmountInCents}}
                "value": "{{float64ToStringTruncate "%.2f" 2 (convertAmountInCentsToPercent .Title.AmountInCents .Title.Fees.Fine.AmountInCents)}}"
            {{else}}
                "value": "{{float64ToString "%.2f" .Title.Fees.Fine.PercentageOnTotal}}"
            {{end}}
        }
    {{end}}
    {{if .Title.Fees.HasInterest}}
        ,"interest": {
            "date": "{{enDate (datePlusDays .Title.ExpireDateTime .Title.Fees.Interest.DaysAfterExpirationDate) "-"}}",
            {{if .Title.Fees.Interest.HasAmountPerDayInCents}}
                "value": "{{float64ToStringTruncate "%.2f" 2 (convertAmountInCentsToPercentPerDay .Title.AmountInCents .Title.Fees.Interest.AmountPerDayInCents)}}"
            {{else}}
                "value": "{{float64ToString "%.2f" .Title.Fees.Interest.PercentagePerMonth}}"
            {{end}}
        }
    {{end}}
}`

const templateResponse = `
{
    "barcode": "{{barCodeNumber}}",
    "our_number": "{{ourNumber}}",
    "writable_line": "{{digitableLine}}"
}
`

const templateError = `
{
    "reason": "{{messageError}}",
    "type": "{{errorCode}}"
}
`

const templateAPI = `
{
    {{if (hasErrorTags . "errorCode") | (hasErrorTags . "messageError")}}
    "Errors": [
        {
        {{if (hasErrorTags . "errorCode")}}
            "Code": "{{trim .errorCode}}",
        {{end}}
        {{if (eq .messageError "{}")}}
            "Message": "{{trim .errorCode}}"
        {{else}}
            "Message": "{{trim .messageError}}"
        {{end}}
        }
    ]
    {{else}}
    "DigitableLine": "{{fmtDigitableLine (trim .digitableLine)}}",
    "BarCodeNumber": "{{trim .barCodeNumber}}",
        "OurNumber": "{{.ourNumber}}"
    {{end}}
}
`
