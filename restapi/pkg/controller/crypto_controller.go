package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/la4ezar/restapi/internal/routes"

	"github.com/gorilla/mux"
	"github.com/la4ezar/restapi/internal/crypto"
	"github.com/la4ezar/restapi/pkg/log"
	"github.com/la4ezar/restapi/pkg/storage"
)

type Route struct {
	Name    string
	Method  string
	Path    string
	Handler http.HandlerFunc
}

type Routes []Route

type Controller struct {
	repository storage.Repository
}

// getCryptos returns http.HandlerFunc
// which encodes all the cryptos in the http.ResponseWriter
func (c *Controller) getAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := r.Body.Close()
			logOnError("an error occurred while closing request body", err)
		}()

		setHeaders(&w)

		cryptos, err := c.repository.GetAllCryptos()
		logOnError("an error occurred while getting all cryptos from repository", err)

		err = json.NewEncoder(w).Encode(cryptos) // Response with all cryptos
		logOnError("an error occurred while encoding cryptos", err)
	}
}

// getCrypto returns http.HandlerFunc
// which encodes requested crypto or default crypto in the http.ResponseWriter
func (c *Controller) getByCryptoID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := r.Body.Close()
			logOnError("an error occurred while closing request body", err)
		}()

		setHeaders(&w)

		params := mux.Vars(r)

		crypto, err := c.repository.GetSingleCrypto(params["crypto_id"])
		logOnError(fmt.Sprintf("an error occurred while getting crypto with CryptoID=%s from repository", params["crypto_id"]), err)

		err = json.NewEncoder(w).Encode(crypto)
		logOnError("an error occurred while encoding crypto", err)
	}
}

// createCrypto returns http.HandlerFunc
// which appends new crypto to cryptos and
// encodes it in the http.ResponseWriter
func (c *Controller) add() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := r.Body.Close()
			logOnError("an error occurred while closing request body", err)
		}()

		setHeaders(&w)

		var crypto crypto.Cryptocurrency
		_ = json.NewDecoder(r.Body).Decode(&crypto)

		err := c.repository.AddCrypto(crypto)
		logOnError("an error occurred while adding crypto to repository", err)

		err = json.NewEncoder(w).Encode(crypto) // Response with the new crypto
		logOnError("an error occurred while encoding crypto", err)
	}
}

// updateCrypto returns http.HandlerFunc
// which updates existing crypto and
// encodes all cryptos in the http.ResponseWriter
func (c *Controller) update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := r.Body.Close()
			logOnError("an error occurred while closing request body", err)
		}()

		setHeaders(&w)

		params := mux.Vars(r)

		newCrypto := crypto.Cryptocurrency{}
		_ = json.NewDecoder(r.Body).Decode(&newCrypto)

		err := c.repository.UpdateCrypto(params["crypto_id"], newCrypto)
		logOnError("an error occurred while updating crypto in repository", err)

		cryptos, err := c.repository.GetAllCryptos()
		logOnError("an error occurred while getting all cryptos from repository", err)

		err = json.NewEncoder(w).Encode(cryptos)
		logOnError("an error occurred while encoding cryptos", err)

	}
}

// removeCrypto returns http.HandlerFunc
// which removes existing crypto and
// encodes all cryptos in the http.ResponseWriter
func (c *Controller) remove() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := r.Body.Close()
			logOnError("an error occurred while closing request body", err)
		}()
		setHeaders(&w)

		params := mux.Vars(r)

		err := c.repository.RemoveCrypto(params["crypto_id"])
		logOnError("an error occurred while deleting crypto from repository", err)

		cryptos, err := c.repository.GetAllCryptos()
		logOnError("an error occurred while getting all cryptos from repository", err)

		err = json.NewEncoder(w).Encode(cryptos) // Response with all cryptos
		logOnError("an error occurred while encoding cryptos", err)
	}
}

func (c *Controller) healthCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := r.Body.Close()
			logOnError("an error occurred while closing request body", err)
		}()

		if err := c.repository.Ping(); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}
}

// setHeaders sets http.ResponseWriter headers
func setHeaders(w *http.ResponseWriter) {
	// Set Content-Type to accept json
	(*w).Header().Set("Content-Type", "application/json")
}

// logOnError logs error message if err is not nil
func logOnError(msg string, err error) {
	if err != nil {
		log.D().WithError(err).Errorf("%s: %v", msg, err)
	}
}

func NewController(repository storage.Repository) *Controller {
	return &Controller{
		repository: repository,
	}
}

func (c *Controller) Routes() *Routes {
	return &Routes{
		{
			Name:    "Get all cryptos",
			Method:  http.MethodGet,
			Path:    routes.AllCryptosURL,
			Handler: c.getAll(),
		},
		{
			Name:    "Get specific crypto",
			Method:  http.MethodGet,
			Path:    routes.SingleCryptoURL,
			Handler: c.getByCryptoID(),
		},
		{
			Name:    "Create crypto",
			Method:  http.MethodPost,
			Path:    routes.AddCryptoURL,
			Handler: c.add(),
		},
		{
			Name:    "Update existing crypto",
			Method:  http.MethodPut,
			Path:    routes.UpdateCryptoURL,
			Handler: c.update(),
		},
		{
			Name:    "Remove existing crypto",
			Method:  http.MethodDelete,
			Path:    routes.RemoveCryptoURL,
			Handler: c.remove(),
		},
		{
			Name:    "Health Check",
			Method:  http.MethodGet,
			Path:    routes.HealthCheckURL,
			Handler: c.healthCheck(),
		},
	}
}
