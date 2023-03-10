package types

import (
	"fmt"
	"net/smtp"

	"github.com/AlecAivazis/survey/v2"
	"github.com/idestis/feathr-cli/helpers"
)

type SMTPConfig struct {
	Host     string `survey:"host" yaml:"host"`         // The hostname of the SMTP server.
	Port     int    `survey:"port" yaml:"port"`         // The port number to use for the SMTP server.
	Username string `survey:"username" yaml:"username"` // The username to use for authentication with the SMTP server.
	Password string `survey:"password" yaml:"password"` // The password to use for authentication with the SMTP server.
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

func PromptSMTP() (SMTPConfig, error) {
	smtp := SMTPConfig{}
	qs := []*survey.Question{
		{
			Name: "host",
			Prompt: &survey.Input{
				Message: "SMTP Host",
				Default: "smtp.gmail.com",
			},
			Validate: survey.Required,
		},
		{
			Name: "port",
			Prompt: &survey.Input{
				Message: "SMTP Port",
				Default: "587",
			},
			Validate: survey.Required,
		},
		{
			Name: "username",
			Prompt: &survey.Input{
				Message: "SMTP Username",
			},
			Validate: survey.Required,
		},
		{
			Name: "password",
			Prompt: &survey.Password{
				Message: "SMTP Password",
			},
			Validate: survey.Required,
		},
	}
	if err := survey.Ask(qs, &smtp); err != nil {
		return smtp, err
	}

	if encryptedPassword, err := helpers.EncryptString(smtp.Password); err != nil {
		return smtp, fmt.Errorf("failed to encrypt password")
	} else {
		smtp.Password = encryptedPassword
	}

	return smtp, nil
}
