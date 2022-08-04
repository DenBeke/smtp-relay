package smtprelay

// SMTPRelay holds the config.
type SMTPRelay struct {
	config *Config
}

// New creates a new SMTP Relay with the given config
func New(config *Config) (*SMTPRelay, error) {

	return &SMTPRelay{config: config}, nil

}
