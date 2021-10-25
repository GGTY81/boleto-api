package issuer

import "regexp"

type Issuer struct {
	barCode       string
	digitableLine string
}

const (
	barCodePattern       = `^\d+$`
	digitableLinePattern = `^\d{5}\.\d{5}\s\d{5}\.\d{6}\s\d{5}\.\d{6}\s\d{1}\s\d+$`
)

//NewIssuer Instancia um objeto Issuer com código de barras e linha digitável
func NewIssuer(barcode, digitableLine string) *Issuer {
	issuer := new(Issuer)
	issuer.barCode = barcode
	issuer.digitableLine = digitableLine

	return issuer
}

//IsValidBarCode Verifica se o barCode contém apenas dígitos numéricos
func (i *Issuer) IsValidBarCode() bool {
	isValid, _ := regexp.Match(barCodePattern, []byte(i.barCode))

	return isValid
}

//IsValidDigitableLine Verifica se o digitableLine está de acordo com o padrão 99999.99999 99999.999999 99999.999999 9 99999999999999
func (i *Issuer) IsValidDigitableLine() bool {
	isValid, _ := regexp.Match(digitableLinePattern, []byte(i.digitableLine))

	return isValid
}
