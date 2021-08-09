package config

type Validator interface {
	Validate() error
}

type ConfigFile struct {
	Name     string
	Location string
	Format   string
}

func DefaultConfigFile() *ConfigFile {
	return &ConfigFile{
		Name:     "application",
		Location: ".",
		Format:   "yml",
	}
}
