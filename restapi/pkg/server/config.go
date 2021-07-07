// Package server contains custom http server for our API
package server // import "github.com/la4ezar/restapi/pkg/server

import (
	"fmt"
	"time"
)

// Config contains Server settings
type Config struct {
	Port            int           `mapstructure:"port" description:"port of the server"`
	ReadTimeout     time.Duration `mapstructure:"readtimeout" description:"read timeout duration for the server"`
	WriteTimeout    time.Duration `mapstructure:"writetimeout" description:"write timeout duration for the server"`
	IdleTimeout     time.Duration `mapstructure:"idletimeout" description:"idle timeout duration for the server"`
	ShutdownTimeout time.Duration `mapstructure:"shutdowntimeout" description:"time to wait for the server to shutdown"`
}

// DefaultConfig returns the default values for configuring the Server
func DefaultConfig() *Config {
	return &Config{
		Port: 8080,

		ReadTimeout:     15 * time.Second,
		WriteTimeout:    15 * time.Second,
		IdleTimeout:     45 * time.Second,
		ShutdownTimeout: 10 * time.Second,
	}
}

// Validate validates the server settings
func (c *Config) Validate() error {
	if c.Port == 0 {
		return fmt.Errorf("validate Server settings: Port missing")
	}
	if c.ReadTimeout == 0 {
		return fmt.Errorf("validate Server settings: ReadTimeout missing")
	}
	if c.WriteTimeout == 0 {
		return fmt.Errorf("validate Server settings: WriteTimeout missing")
	}
	if c.IdleTimeout == 0 {
		return fmt.Errorf("validate Server settings: IdleTimeout missing")
	}
	if c.ShutdownTimeout == 0 {
		return fmt.Errorf("validate Server settings: ShutdownTimeout missing")
	}

	return nil
}
