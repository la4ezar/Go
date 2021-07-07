package log

import (
	"fmt"
	"os"
)

type Config struct {
	Level  string `mapstructure:"level" description:"Logger level"` // Only log entries with same level or above it will be logged
	Format string `mapstructure:"format" description:"The format that the logger will use"`
	Output string `mapstructure:"output" description:"Where the logs will be outputted"`
}

func DefaultConfig() *Config {
	return &Config{
		Level:  "info",
		Format: "text",
		Output: os.Stdout.Name(),
	}
}

func (c *Config) Validate() error {
	if len(c.Level) == 0 {
		return fmt.Errorf("validate Logger settings: Level missing")
	}
	if len(c.Format) == 0 {
		return fmt.Errorf("validate Logger settings: Format missing")
	}
	if len(c.Output) == 0 {
		return fmt.Errorf("validate Logger settings: Output missing")
	}

	return nil
}
