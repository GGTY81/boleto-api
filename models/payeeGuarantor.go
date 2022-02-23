package models

// PayeeGuarantor informações de entrada do lojista(sacador avalista)
type PayeeGuarantor struct {
	Name     string   `json:"name,omitempty"`
	Document Document `json:"document,omitempty"`
}

// HasName diz se o Name está preenchido com algum valor
func (p PayeeGuarantor) HasName() bool {
	return p.Name != ""
}
