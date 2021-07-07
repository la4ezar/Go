package storage

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Storage struct {
	DB *sql.DB
}

func (s *Storage) Close() error {
	return s.DB.Close()
}

func New(c *Config) (*Storage, error) {
	db, err := sql.Open(c.Type, c.DataSource.String())
	if err != nil {
		return nil, fmt.Errorf("unable to open db connection: %s", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("database is not alive: %s", err)
	}

	log.Println("Database is up-to-date")

	return &Storage{
		DB: db,
	}, nil
}
