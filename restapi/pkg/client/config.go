package client

import (
	"fmt"
	"time"
)

type Config struct {
	Endpoints         map[string]string `mapstructure:"endpoints" description:"All client http requests endpoints"`
	Timeout           time.Duration     `mapstructure:"timeout" description:"Client timeout"`
	DisableKeepAlives bool              `mapstructure:"disable_keep_alives" description:"Whether to disable http keep-alives"`
}

func DefaultConfig() *Config {
	return &Config{
		Endpoints:         make(map[string]string),
		Timeout:           15 * time.Second,
		DisableKeepAlives: false,
	}
}

func (c *Config) Validate() error {
	if len(c.Endpoints) == 0 {
		return fmt.Errorf("validate Client settings: Endpoints missing")
	}
	if c.Timeout <= 0*time.Second {
		return fmt.Errorf("validate Client settings: Timeout missing")
	}

	return nil
}
