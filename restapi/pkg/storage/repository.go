package storage

import (
	"fmt"

	"github.com/la4ezar/restapi/internal/crypto"
	"github.com/la4ezar/restapi/pkg/log"
)

type Repository interface {
	GetAllCryptos() ([]crypto.Cryptocurrency, error)
	GetSingleCrypto(cryptoID string) (crypto.Cryptocurrency, error)
	AddCrypto(c crypto.Cryptocurrency) error
	UpdateCrypto(oldCryptoID string, c crypto.Cryptocurrency) error
	RemoveCrypto(cryptoID string) error
	Ping() error
}

type RepositoryImpl struct {
	storage Storage
}

func NewRepository(storage Storage) Repository {
	return &RepositoryImpl{
		storage: storage,
	}
}

// GetAllCryptos retrieves all cryptos from a sql.r.storage.DB
func (r *RepositoryImpl) GetAllCryptos() ([]crypto.Cryptocurrency, error) {
	var cryptos []crypto.Cryptocurrency

	cryptoRows, err := r.storage.DB.Query("SELECT * FROM CRYPTOS.CRYPTOCURRENCIES")
	if err != nil {
		return cryptos, fmt.Errorf("an error occurred while querying cryptos from DB: %v", err)
	}
	defer func() {
		if err := cryptoRows.Close(); err != nil {
			log.D().WithError(err).Errorf("an error occurred while closing DB cryptos rows: %v", err)
		}
	}()

	for cryptoRows.Next() {
		cryptocurrency := &crypto.Cryptocurrency{}
		if err := cryptoRows.Scan(&cryptocurrency.Name, &cryptocurrency.CryptoID, &cryptocurrency.Price); err != nil {
			return cryptos, fmt.Errorf("an error occurred while scanning cryptos row: %v", err)
		}
		cryptos = append(cryptos, *cryptocurrency)
	}
	if err := cryptoRows.Err(); err != nil {
		return cryptos, fmt.Errorf("an error occurred while iterating cryptos rows: %v", err)
	}
	if err := cryptoRows.Close(); err != nil {
		return cryptos, fmt.Errorf("an error occurred while closing DB cryptos rows: %v", err)
	}

	authorsRows, err := r.storage.DB.Query("SELECT * FROM CRYPTOS.AUTHORS")
	if err != nil {
		return cryptos, fmt.Errorf("an error occurred while querying authors from DB: %v", err)
	}
	defer func() {
		if err := authorsRows.Close(); err != nil {
			log.D().WithError(err).Errorf("an error occurred while closing DB author rows: %v", err)
		}
	}()
	for authorsRows.Next() {
		currentCrypto := &crypto.Cryptocurrency{}
		author := crypto.Author{}
		if err := authorsRows.Scan(&currentCrypto.CryptoID, &author.Firstname, &author.Lastname); err != nil {
			return cryptos, fmt.Errorf("an error occurred while scanning authors row: %v", err)
		}
		for i := range cryptos {
			if cryptos[i].CryptoID == currentCrypto.CryptoID {
				cryptos[i].Authors = append(cryptos[i].Authors, author)
			}
		}
	}
	if err := authorsRows.Err(); err != nil {
		return cryptos, fmt.Errorf("an error occurred while iterating authors rows: %v", err)
	}
	if err := authorsRows.Close(); err != nil {
		return cryptos, fmt.Errorf("an error occurred while closing DB authors rows: %v", err)
	}

	return cryptos, nil
}

// GetSingleCrypto retrieves single crypto from sql.r.storage.DB
func (r *RepositoryImpl) GetSingleCrypto(cryptoID string) (crypto.Cryptocurrency, error) {
	var cryptocurrency crypto.Cryptocurrency

	if err := r.storage.DB.QueryRow("SELECT * FROM CRYPTOS.CRYPTOCURRENCIES WHERE CRYPTOID = $1",
		cryptoID).Scan(&cryptocurrency.Name, &cryptocurrency.CryptoID, &cryptocurrency.Price); err != nil {
		return cryptocurrency, fmt.Errorf("an error occurred while querying cryptos from DB: %v", err)
	}

	// Get crypto authors
	authorsRows, err := r.storage.DB.Query("SELECT * FROM CRYPTOS.AUTHORS WHERE CRYPTOID = $1", cryptoID)
	if err != nil {
		return cryptocurrency, fmt.Errorf("an error occurred while querying authors from DB: %v", err)
	}
	defer func() {
		if err := authorsRows.Close(); err != nil {
			log.D().WithError(err).Errorf("an error occurred while closing DB author rows: %v", err)
		}
	}()
	for authorsRows.Next() {
		currentCrypto := &crypto.Cryptocurrency{}
		author := crypto.Author{}
		if err := authorsRows.Scan(&currentCrypto.CryptoID, &author.Firstname, &author.Lastname); err != nil {
			return cryptocurrency, fmt.Errorf("an error occurred while scanning authors row: %v", err)
		}
		if cryptocurrency.CryptoID == currentCrypto.CryptoID {
			cryptocurrency.Authors = append(cryptocurrency.Authors, author)
		}
	}
	if err := authorsRows.Err(); err != nil {
		return cryptocurrency, fmt.Errorf("an error occurred while iterating authors rows: %v", err)
	}
	if err := authorsRows.Close(); err != nil {
		return cryptocurrency, fmt.Errorf("an error occurred while closing DB authors rows: %v", err)
	}

	return cryptocurrency, nil
}

func (r *RepositoryImpl) AddCrypto(c crypto.Cryptocurrency) error {
	if _, err := r.storage.DB.Exec("INSERT INTO CRYPTOS.CRYPTOCURRENCIES(NAME, CRYPTOID, PRICE) VALUES ($1, $2, $3)", c.Name, c.CryptoID, c.Price); err != nil {
		return fmt.Errorf("an error occurred while inserting crypto in DB: %v", err)
	}

	for _, a := range c.Authors {
		if _, err := r.storage.DB.Exec("INSERT INTO CRYPTOS.AUTHORS(CRYPTOID, FIRSTNAME, LASTNAME) VALUES ($1, $2, $3)", c.CryptoID, a.Firstname, a.Lastname); err != nil {
			return fmt.Errorf("an error occurred while inserting author in DB: %v", err)
		}
	}

	return nil
}

func (r *RepositoryImpl) UpdateCrypto(oldCryptoID string, c crypto.Cryptocurrency) error {
	if _, err := r.storage.DB.Exec("UPDATE CRYPTOS.CRYPTOCURRENCIES SET NAME = $1, CRYPTOID = $2, PRICE = $3 WHERE CRYPTOID = $4", c.Name, c.CryptoID, c.Price, oldCryptoID); err != nil {
		return fmt.Errorf("an error occurred while updating crypto in DB: %v", err)
	}

	if _, err := r.storage.DB.Exec("DELETE FROM CRYPTOS.AUTHORS WHERE CRYPTOID = $1", c.CryptoID); err != nil {
		return fmt.Errorf("an error occurred while deleting authors in DB: %v", err)
	}

	for _, a := range c.Authors {
		if _, err := r.storage.DB.Exec("INSERT INTO CRYPTOS.AUTHORS(CRYPTOID, FIRSTNAME, LASTNAME) VALUES ($1, $2, $3)", c.CryptoID, a.Firstname, a.Lastname); err != nil {
			return fmt.Errorf("an error occurred while inserting author in DB: %v", err)
		}
	}

	return nil
}

func (r *RepositoryImpl) RemoveCrypto(cryptoID string) error {
	if _, err := r.storage.DB.Exec("DELETE FROM CRYPTOS.CRYPTOCURRENCIES WHERE CRYPTOID = $1", cryptoID); err != nil {
		return fmt.Errorf("an error occurred while deleting crypto in DB: %v", err)
	}
	return nil
}

func (r *RepositoryImpl) Ping() error {
	return r.storage.DB.Ping()
	// TODO check pingcontext
}
