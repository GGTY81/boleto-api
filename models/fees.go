package models

type Fees struct {
	Fine     *Fine     `json:"fine,omitempty"`
	Interest *Interest `json:"interest,omitempty"`
}

//HasFine Verifica se o nó de fine está preenchido
func (f *Fees) HasFine() bool {
	return f != nil && f.Fine != nil
}

//HasInterest Verifica se o nó de interest está preenchido
func (f *Fees) HasInterest() bool {
	return f != nil && f.Interest != nil
}
