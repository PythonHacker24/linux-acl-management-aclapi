package config

/* daemon config */
type DConfig struct {
	DebugMode 	bool  `yaml:"debug_mode,omitempty"`   
}

/* normalization function */
func (d *DConfig) Normalize() error {
	
	/* 
		debug_mode is false by default
		daemon will run on production mode by default
	*/

	return nil
}
