package types

import (
	"fmt"
	"net/smtp"
)

type SMTPConfig struct {
	Host     string `survey:"host"`     // The hostname of the SMTP server.
	Port     int    `survey:"port"`     // The port number to use for the SMTP server.
	Username string `survey:"username"` // The username to use for authentication with the SMTP server.
	Password string `survey:"password"` // The password to use for authentication with the SMTP server.
}

// NewSMTPClient returns a new SMTP client using the provided configuration.
func NewSMTPClient(config SMTPConfig) (*smtp.Client, error) {
	// Connect to the SMTP server.
	client, err := smtp.Dial(fmt.Sprintf("%s:%d", config.Host, config.Port))
	if err != nil {
		return nil, err
	}

	// Perform the SMTP handshake.
	err = client.Hello("localhost")
	if err != nil {
		return nil, err
	}

	// Authenticate with the SMTP server.
	auth := smtp.PlainAuth("", config.Username, config.Password, config.Host)
	err = client.Auth(auth)
	if err != nil {
		return nil, err
	}

	return client, nil
}
