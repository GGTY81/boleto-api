package models

// Buyer informações de entrada do comprador
type Buyer struct {
	Name     string   `json:"name,omitempty"`
	Email    string   `json:"email,omitempty"`
	Document Document `json:"document,omitempty"`
	Address  Address  `json:"address,omitempty"`
}

//HasAddress verify if Adress is not empty
func (b Buyer) HasAddress() bool {
	return b.Address != Address{}
}
