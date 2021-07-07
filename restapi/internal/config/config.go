package config

import (
	"fmt"

	"github.com/la4ezar/restapi/pkg/log"

	"github.com/la4ezar/restapi/pkg/client"
	"github.com/la4ezar/restapi/pkg/server"
	"github.com/la4ezar/restapi/pkg/storage"
	"github.com/spf13/viper"
)

type Validator interface {
	Validate() error
}

type ConfigFile struct {
	Name     string
	Location string
	Format   string
}

type Configs struct {
	Server  *server.Config
	Client  *client.Config
	Storage *storage.Config
	Logger  *log.Config
}

func (c *Configs) Validate() error {
	validatable := []Validator{c.Server, c.Client, c.Storage}

	for _, v := range validatable {
		if err := v.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func DefaultConfigs() *Configs {
	return &Configs{
		Server:  server.DefaultConfig(),
		Client:  client.DefaultConfig(),
		Storage: storage.DefaultConfig(),
		Logger:  log.DefaultConfig(),
	}
}

func DefaultConfigFile() *ConfigFile {
	return &ConfigFile{
		Name:     "application",
		Location: ".",
		Format:   "yml",
	}
}

func New() (*Configs, error) {
	configs := DefaultConfigs()
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

	if err := v.Unmarshal(configs); err != nil {
		return nil, fmt.Errorf("error loading configuration: %s", err)
	}

	return configs, nil
}
