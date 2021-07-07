package storage

import "fmt"

type Config struct {
	Type       string     `mapstructure:"type" description:"Type of the storage"`
	DataSource DataSource `mapstructure:"data_source" description:"Data source name of the storage"`
}

func DefaultConfig() *Config {
	return &Config{
		Type:       "postgres",
		DataSource: DefaultDataSource(),
	}
}

func (c *Config) Validate() error {
	if len(c.Type) == 0 {
		return fmt.Errorf("validate Storage settings: Type missing")
	}
	if err := c.DataSource.Validate(); err != nil {
		return fmt.Errorf("validate Storage settings: %v", err.Error())
	}

	return nil
}
