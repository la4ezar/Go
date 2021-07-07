// Package crypto contains Cryptocurrency structure
package crypto // import "github.com/la4ezar/restapi/internal/server

type Authors []Author

// Cryptocurrency structure with crypto's id, name, crypto_id, current price and authors/innovators
type Cryptocurrency struct {
	Name     string   `json:"name"`
	CryptoID string   `json:"crypto_id"`
	Price    float64  `json:"price"`
	Authors  []Author `json:"authors"`
}

// Author structure with crypto author's first and last name
type Author struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

// DefaultCrypto returns default Cryptocurrency
func DefaultCrypto() *Cryptocurrency {
	return &Cryptocurrency{
		Name:     "",
		CryptoID: "",
		Price:    0.0,
		Authors:  []Author{},
	}
}
