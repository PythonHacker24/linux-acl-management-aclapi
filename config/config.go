package config

import "fmt"

var (
	APIDConfig apidconfig
)

/* config struct for aclapi */
type apidconfig struct {
	Logging		Logging 	`yaml:"logs,omitempty"`
	Server      Server      `yaml:"server,omitempty"`
}

/* complete config normalizer function */
func (c *apidconfig) Normalize() error {
	
	if err := c.Logging.Normalize(); err != nil {
		return fmt.Errorf("logging configuration error: %w", err)
	}

	if err := c.Server.Normalize(); err != nil {
		return fmt.Errorf("server configuration error: %w", err)
	}

	return nil 
}
