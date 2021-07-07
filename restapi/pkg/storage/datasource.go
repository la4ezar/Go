package storage

import "fmt"

type DataSource struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func (ds DataSource) String() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		ds.Host, ds.Port, ds.User, ds.Password, ds.DBName, ds.SSLMode)
}

func DefaultDataSource() DataSource {
	return DataSource{
		Host:     "localhost",
		Port:     "5432",
		User:     "postgres",
		Password: "postgres",
		DBName:   "postgres",
		SSLMode:  "disable",
	}
}

func NewDSN(opts ...Option) *DataSource {
	dsn := &DataSource{}
	for _, opt := range opts {
		opt(dsn)
	}
	return dsn
}

func (ds *DataSource) Validate() error {
	if len(ds.Host) == 0 {
		return fmt.Errorf("validate DataSource settings: Host missing")
	}
	if len(ds.Port) == 0 {
		return fmt.Errorf("validate DataSource settings: Port missing")
	}
	if len(ds.User) == 0 {
		return fmt.Errorf("validate DataSource settings: User missing")
	}
	if len(ds.Password) == 0 {
		return fmt.Errorf("validate DataSource settings: Password missing")
	}
	if len(ds.DBName) == 0 {
		return fmt.Errorf("validate DataSource settings: Database name missing")
	}
	if len(ds.SSLMode) == 0 {
		return fmt.Errorf("validate DataSource settings: SSL Mode missing")
	}

	return nil
}

// Functional options

type Option func(dsn *DataSource)

func SetHost(host string) Option {
	return func(dsn *DataSource) {
		dsn.Host = host
	}
}

func SetPort(port string) Option {
	return func(dsn *DataSource) {
		dsn.Port = port
	}
}

func SetUser(user string) Option {
	return func(dsn *DataSource) {
		dsn.User = user
	}
}

func SetPassword(password string) Option {
	return func(dsn *DataSource) {
		dsn.Password = password
	}
}
func SetDBName(dbname string) Option {
	return func(dsn *DataSource) {
		dsn.DBName = dbname
	}
}
func SetSSLMode(sslmode string) Option {
	return func(dsn *DataSource) {
		dsn.SSLMode = sslmode
	}
}
