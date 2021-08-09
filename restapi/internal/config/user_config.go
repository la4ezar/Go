package config

import (
	"fmt"

	"github.com/la4ezar/restapi/pkg/client"
	"github.com/spf13/viper"
)

type ClientConfig struct {
	Client *client.Config
}

func (c *ClientConfig) Validate() error {
	validatable := []Validator{c.Client}

	for _, v := range validatable {
		if err := v.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func DefaultClientConfig() *ClientConfig {
	return &ClientConfig{
		Client: client.DefaultConfig(),
	}
}

func NewDefaultClientConfig() (*ClientConfig, error) {
	ClientConfig := DefaultClientConfig()
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

	if err := v.Unmarshal(ClientConfig); err != nil {
		return nil, fmt.Errorf("error loading configuration: %s", err)
	}

	return ClientConfig, nil
}
