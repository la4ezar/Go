// Package client contains the custom http client for our API
package client // import "github.com/la4ezar/restapi/pkg/client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/la4ezar/restapi/internal/crypto"
	"github.com/la4ezar/restapi/pkg/log"

	"github.com/sirupsen/logrus"
)

type Cryptocurrency crypto.Cryptocurrency
type CryptoAuthors []crypto.Author

// Client wrapped http.Client
type Client struct {
	*http.Client

	Endpoints map[string]string
}

func New(c *Config) *Client {
	return &Client{
		Client: &http.Client{
			Timeout: c.Timeout,
			Transport: &http.Transport{
				DisableKeepAlives: c.DisableKeepAlives,
			},
		},
		Endpoints: c.Endpoints,
	}
}

// GetCryptos sends GET HTTP request to URL
// and prints out all the cryptocurrencies in the server on stdout
func (c *Client) GetCryptos() {
	url := c.Endpoints["getcryptos"]

	response, err := c.Get(url)
	logOnError(fmt.Sprintf("An error occurred while making GET request to %s", url), err)

	defer func() {
		err := response.Body.Close()
		logOnError("An error occurred while closing response body", err)
	}()

	logAllCryptos(response)
}

// GetCrypto sends GET HTTP request to URL/{cryptoId}
// and if crypto with such cryptoID exists it will be printed on the stdout
// if not - default crypto will be printed
func (c *Client) GetCrypto(cryptoID string) {
	url := c.Endpoints["getcrypto"] + cryptoID

	response, err := c.Get(url)
	logOnError(fmt.Sprintf("An error occurred while making GET request to %s", url), err)

	defer func() {
		err := response.Body.Close()
		logOnError("An error occurred while closing response body", err)
	}()

	logSingleCrypto(response)
}

// PostCrypto sends POST HTTP request to URL and creates new crypto - crypto
func (c *Client) PostCrypto(crypto Cryptocurrency) {
	url := c.Endpoints["postcrypto"]

	requestBody, err := json.Marshal(crypto)
	logOnError("An error occurred while marshalling crypto", err)

	response, err := c.Post(url, "application/json", bytes.NewBuffer(requestBody))
	logOnError(fmt.Sprintf("An error occurred while making POST request to %s", url), err)
	defer func() {
		err := response.Body.Close()
		logOnError("An error occurred while closing response body", err)
	}()

	logSingleCrypto(response)
}

// PutCrypto sends a PUT HTTP request to URL/{crypto.cryptoID}
// and if crypto with such cryptoID exists it will be replaced with crypto
func (c *Client) PutCrypto(crypto Cryptocurrency) {
	url := c.Endpoints["putcrypto"] + crypto.CryptoID

	requestBody, err := json.Marshal(crypto)
	logOnError("An error occurred while marshalling crypto", err)

	request, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(requestBody))
	logOnError("An error occurred while creating GET request.", err)

	setHeaders(request)
	response, err := c.Do(request)
	logOnError(fmt.Sprintf("An error occured while making PUT request to %s", url), err)

	// Handle Closer error
	defer func() {
		err := response.Body.Close()
		logOnError("An error occurred while closing response body", err)
	}()

	logAllCryptos(response)
}

// DeleteCrypto sends a DELETE HTTP request to URL/{cryptoID}
// and if crypto with such cryptoID exists it will be deleted
func (c *Client) DeleteCrypto(cryptoID string) {
	url := c.Endpoints["deletecrypto"] + cryptoID

	request, err := http.NewRequest(http.MethodDelete, url, nil)
	logOnError("An error occurred while creating DELETE request.", err)

	setHeaders(request)

	response, err := c.Do(request)
	logOnError(fmt.Sprintf("An error occured while making DELETE request to %s", url), err)
	defer func() {
		err := response.Body.Close()
		logOnError("An error occurred while closing response body", err)
	}()

	logAllCryptos(response)
}

// HealthCheck sends a GET HTTP request to URL/health
// and returns the server status code.
func (c *Client) HealthCheck() {
	url := c.Endpoints["healthcheck"]

	response, err := http.Get(url)
	logOnError(fmt.Sprintf("An error occurred while making GET request to %s", url), err)

	defer func() {
		err := response.Body.Close()
		logOnError("An error occurred while closing response body", err)
	}()

	log.D().WithFields(logrus.Fields{
		"status code": response.StatusCode,
	}).Info("Health Check...")
}

// logOnError logs error message if err is not nil
func logOnError(msg string, err error) {
	if err != nil {
		log.D().WithError(err).Errorf("%s: %v", msg, err)
	}
}

// setHeaders sets http.Request headers
func setHeaders(r *http.Request) {
	r.Header.Set("Content-Type", "application/json")
}

// logAllCryptos decodes the cryptos in http.Response Body and logs them
func logAllCryptos(response *http.Response) {
	log.D().Info("Server response...")

	var cryptos []Cryptocurrency
	_ = json.NewDecoder(response.Body).Decode(&cryptos)
	for _, crypto := range cryptos {
		log.D().WithFields(logrus.Fields{
			"name":     crypto.Name,
			"cryptoid": crypto.CryptoID,
			"price":    crypto.Price,
			"authors":  crypto.Authors,
		}).Info()
	}
}

// logSingleCrypto decodes the crypto in http.Response Body and logs it
func logSingleCrypto(response *http.Response) {
	log.D().Info("Server response...")

	var crypto Cryptocurrency
	_ = json.NewDecoder(response.Body).Decode(&crypto)
	log.D().WithFields(logrus.Fields{
		"name":     crypto.Name,
		"cryptoid": crypto.CryptoID,
		"price":    crypto.Price,
		"authors":  crypto.Authors,
	}).Info()
}
