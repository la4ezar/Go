// Making requests to the server with custom client
package main // import "github.com/la4ezar/restapi"

import (
	"github.com/la4ezar/restapi/internal/config"
	"github.com/la4ezar/restapi/pkg/client"
	"github.com/la4ezar/restapi/pkg/log"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.D().WithError(err).Error()
	}

	if err := cfg.Validate(); err != nil {
		log.D().WithError(err).Error()
	}

	// Init Client
	c := client.New(cfg.Client)

	// Health check
	c.HealthCheck()

	// GET requests
	c.GetCryptos()
	c.GetCrypto("BTC")
	c.GetCrypto("LTC")

	//POST request
	crypto := client.Cryptocurrency{
		Name:     "LachoCoin",
		CryptoID: "LCN",
		Price:    1.54,
		Authors:  client.CryptoAuthors{{Firstname: "Lachezar", Lastname: "Bogomilov"}},
	}
	c.PostCrypto(crypto)

	//PUT request
	updatedCrypto := client.Cryptocurrency{
		Name:     "Bitcoin",
		CryptoID: "BTC",
		Price:    45000.3,
		Authors:  client.CryptoAuthors{{Firstname: "Satoshi", Lastname: "Nakamoto"}},
	}
	c.PutCrypto(updatedCrypto)

	// DELETE request
	c.DeleteCrypto("LCN")
}
