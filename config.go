package smtprelay

import (
	"fmt"
	"os"
	"strconv"
)

var (
	defaultMetricsAddress = ":9090"
	defaultSMTPHost       = "localhost"
	defaultSMTPPort       = "25"
)

// Config contains all the config for serving SMTP Relay
type Config struct {
	MetricsAddress string
	AccessLog      bool
	HandleAsync    bool
	SentryDSN      string

	// SMTP server config for SMTP Relay
	SMTP struct {
		Host string
		Port int
	}

	// Remote SMTP to which emails will be relayed
	RemoteSMTP struct {
		Host       string
		Port       int
		User       string
		Password   string
		DisableTLS bool
	}
}

// BuildConfigFromEnv populates a config from env variables
func BuildConfigFromEnv() *Config {
	config := &Config{}

	// SMTP Relay SMTP settings
	config.SMTP.Host = getEnv("SMTP_HOST", defaultSMTPHost)

	port, err := strconv.Atoi(getEnv("SMTP_PORT", defaultSMTPPort))
	if err != nil {
		config.SMTP.Port = 0
	} else {
		config.SMTP.Port = port
	}

	// Remote SMTP server settings
	config.RemoteSMTP.Host = getEnv("REMOTE_SMTP_HOST", "")

	port, err = strconv.Atoi(getEnv("REMOTE_SMTP_PORT", "0"))
	if err != nil {
		config.RemoteSMTP.Port = 0
	} else {
		config.RemoteSMTP.Port = port
	}

	config.RemoteSMTP.User = getEnv("REMOTE_SMTP_USER", "")
	config.RemoteSMTP.Password = getEnv("REMOTE_SMTP_PASSWORD", "")
	if getEnv("REMOTE_SMTP_DISABLE_TLS", "0") == "1" {
		config.RemoteSMTP.DisableTLS = true
	}

	// Access log
	accessLog := getEnv("ACCESS_LOG", "1")
	if accessLog == "0" {
		config.AccessLog = false
	} else {
		config.AccessLog = true
	}

	// Handle Async
	handleAsync := getEnv("HANDLE_ASYNC", "0")
	if handleAsync == "1" {
		config.HandleAsync = true
	} else {
		config.HandleAsync = false
	}

	// Metrics
	config.MetricsAddress = getEnv("METRICS_ADDRESS", defaultMetricsAddress)

	// Sentry DSN
	config.SentryDSN = getEnv("SENTRY_DSN", "")

	return config
}

// Validate validates whether all config is set and valid
func (config *Config) Validate() error {

	// SMTP config
	if config.SMTP.Host == "" {
		return fmt.Errorf("SMTP_HOST must be set")
	}
	if config.SMTP.Port == 0 {
		return fmt.Errorf("SMTP_PORT must be set")
	}

	// Remote SMTP config
	if config.RemoteSMTP.Host == "" {
		return fmt.Errorf("REMOTE_SMTP_HOST must be set")
	}
	if config.RemoteSMTP.Port == 0 {
		return fmt.Errorf("REMOTE_SMTP_PORT must be set")
	}

	// Metrics
	if config.MetricsAddress == "" {
		return fmt.Errorf("METRICS_ADDRESS cannot be empty")
	}

	return nil
}

// getEnv gets the env variable with the given key if the key exists
// else it falls back to the fallback value
func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
