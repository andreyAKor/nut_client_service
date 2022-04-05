package configs

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Config struct {
	// Logging settings.
	Logging struct {
		// Path to the log file.
		File string

		// Logging level, variants levels:
		//  - debug - defines debug log level
		//  - info - defines info log level
		//  - warn - defines warn log level
		//  - error - defines error log level
		//  - fatal - defines fatal log level
		//  - panic - defines panic log level
		//  - no - defines an absent log level
		//  - disabled - disables the logger
		//  - trace - defines trace log level.
		Level string
	}

	// HTTP-server settings
	HTTP struct {
		// Host
		Host string

		// Port
		Port int

		// Maximum content size limit
		BodyLimit int
	}

	Clients struct {
		NUT struct {
			Host     string
			Port     int
			Username string
			Password string
		}
	}

	Metrics struct {
		NUT struct {
			Interval string
		}
	}
}

// Init is using to initialize the current config instance.
func (c *Config) Init(file string) error {
	// read in environment variables that match
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetConfigFile(file)

	if err := viper.ReadInConfig(); err != nil {
		return errors.Wrap(err, "open config file failed")
	}

	if err := viper.Unmarshal(c); err != nil {
		return errors.Wrap(err, "unmarshal config file failed")
	}

	return nil
}
