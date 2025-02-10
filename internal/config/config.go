package config

import (
	"flag"
	"fmt"
	"os"
)

// ErrInitConfigFailed - config initialization error.
var ErrInitConfigFailed = fmt.Errorf("failed to init config")

// Config - application configuration structure.
type Config struct {
	RunAddress  string // Address and port of HTTP server
	DatabaseURI string // Address for database connection
	SecretKey   string // Authentication secret key

}

// configBuilder - application configuration builder.
type configBuilder struct {
	runAddress  string `env:"RUN_ADDRESS"`
	databaseURI string `env:"DATABASE_URI"`
	secretKey   string `env:"SECRET_KEY"`
}

// newConfigBuilder creates new application configuration builder.
func newConfigBuilder() *configBuilder {
	return &configBuilder{}
}

// setDefaults defines application configuration parameters defaults.
func (cb *configBuilder) setDefaults() error {
	cb.runAddress = "localhost:8080"
	cb.databaseURI = ""
	cb.secretKey = "secret"

	return nil
}

// setFlags sets application configuration parameters from command line parameters.
func (cb *configBuilder) setFlags() error {
	if flag.Lookup("a") == nil {
		flag.StringVar(&cb.runAddress, "a", cb.runAddress, "HTTP server address and port")
	}
	if flag.Lookup("d") == nil {
		flag.StringVar(&cb.databaseURI, "d", cb.databaseURI, "database connection string")
	}
	flag.Parse()

	return nil
}

// setEnvs sets application configuration parameters from environment variables.
func (cb *configBuilder) setEnvs() error {
	ra := os.Getenv("RUN_ADDRESS")
	if ra != "" {
		cb.runAddress = ra
	}

	dbi := os.Getenv("DATABASE_URI")
	if dbi != "" {
		cb.databaseURI = dbi
	}

	sk := os.Getenv("SECRET_KEY")
	if sk != "" {
		cb.secretKey = sk
	}

	return nil
}

// build builds application cofiguration.
func (cb *configBuilder) build() *Config {
	return &Config{
		RunAddress:  cb.runAddress,
		DatabaseURI: cb.databaseURI,
		SecretKey:   cb.secretKey,
	}
}

// Get returns application configuration.
func Get() (*Config, error) {
	cb := newConfigBuilder()

	confSets := []func() error{
		cb.setDefaults,
		cb.setFlags,
		cb.setEnvs,
	}

	for _, confSet := range confSets {
		err := confSet()
		if err != nil {
			return nil, ErrInitConfigFailed
		}
	}

	return cb.build(), nil
}
