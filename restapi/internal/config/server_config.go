package config

import (
	"fmt"

	"github.com/la4ezar/restapi/pkg/log"
	"github.com/la4ezar/restapi/pkg/server"
	"github.com/la4ezar/restapi/pkg/storage"
	"github.com/spf13/viper"
)

type ServerConfig struct {
	Server  *server.Config
	Storage *storage.Config
	Logger  *log.Config
}

func (c *ServerConfig) Validate() error {
	validatable := []Validator{c.Server, c.Logger, c.Storage}

	for _, v := range validatable {
		if err := v.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func DefaultServerConfig() *ServerConfig {
	return &ServerConfig{
		Server:  server.DefaultConfig(),
		Storage: storage.DefaultConfig(),
		Logger:  log.DefaultConfig(),
	}
}

func NewDefaultServerConfig() (*ServerConfig, error) {
	serverConfig := DefaultServerConfig()
	configFile := DefaultConfigFile()

	v := viper.New()

	v.AddConfigPath(configFile.Location)
	v.SetConfigName(configFile.Name)
	v.SetConfigType(configFile.Format)

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("could not read configuration: %s", err)
		}
	}

	if err := v.Unmarshal(serverConfig); err != nil {
		return nil, fmt.Errorf("error loading configuration: %s", err)
	}

	return serverConfig, nil
}
