package config

import "fmt"

var (
	APIDConfig ADConfig
)

/* config struct for aclapi */
type ADConfig struct {
	DConfig	DConfig	`yaml:"daemon,omitempty"`
	Logging	Logging	`yaml:"logs,omitempty"`
	Server	Server	`yaml:"server,omitempty"`
}

/* complete config normalizer function */
func (c *ADConfig) Normalize() error {

	if err := c.DConfig.Normalize(); err != nil {
		return fmt.Errorf("daemon configuration error: %w", err)
	}

	if err := c.Logging.Normalize(); err != nil {
		return fmt.Errorf("logging configuration error: %w", err)
	}

	if err := c.Server.Normalize(); err != nil {
		return fmt.Errorf("server configuration error: %w", err)
	}

	return nil
}
