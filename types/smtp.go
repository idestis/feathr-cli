package types

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/smtp"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/idestis/feathr-cli/helpers"
)

type SMTPConfig struct {
	Host     string `survey:"host" yaml:"host"`         // The hostname of the SMTP server.
	Port     int    `survey:"port" yaml:"port"`         // The port number to use for the SMTP server.
	Username string `survey:"username" yaml:"username"` // The username to use for authentication with the SMTP server.
	Password string `survey:"password" yaml:"password"` // The password to use for authentication with the SMTP server.
}

func (cfg *SMTPConfig) SendEmailWithAttachment(to []string, subject, body, attachmentPath string) error {
	// Set up authentication information.
	password, err := helpers.DecryptString(cfg.Password)
	if err != nil {
		return fmt.Errorf("failed to decrypt password: %w", err)
	}
	auth := smtp.PlainAuth("", cfg.Username, password, cfg.Host)

	// Set up the attachment.
	// attachment, err := os.Open(attachmentPath)
	// if err != nil {
	// 	return err
	// }
	// defer attachment.Close()
	attachmentBytes, err := ioutil.ReadFile(attachmentPath)
	if err != nil {
		return fmt.Errorf("failed to read attachment: %w", err)
	}
	attachmentName := filepath.Base(attachmentPath)

	// Encode the attachment as base64.
	b := make([]byte, base64.StdEncoding.EncodedLen(len(attachmentBytes)))
	// encoded := base64.StdEncoding.EncodeToString(attachmentBytes)
	base64.StdEncoding.Encode(b, attachmentBytes)

	// Set up the email headers.
	headers := make(map[string]string)
	headers["From"] = cfg.Username
	headers["To"] = strings.Join(to, ",")
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"

	// Build the email body.
	var emailBody bytes.Buffer
	emailBody.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=\"%s\"\n", "boundarystring"))
	emailBody.WriteString(fmt.Sprintf("--%s\n", "boundarystring"))
	emailBody.WriteString(fmt.Sprintf("%s\n", body))
	emailBody.WriteString(fmt.Sprintf("--%s\n", "boundarystring"))
	emailBody.WriteString(fmt.Sprintf("Content-Type: application/pdf; name=\"%v\"\r\n", attachmentName))
	emailBody.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=\"%v\"\n", attachmentName))
	emailBody.WriteString("Content-Transfer-Encoding: base64\r\n")
	emailBody.Write(b)
	emailBody.WriteString(fmt.Sprintf("--%s--\n", "boundarystring"))

	// Construct the email message.
	msg := []byte{}
	for k, v := range headers {
		msg = append(msg, []byte(fmt.Sprintf("%s: %s\n", k, v))...)
	}
	msg = append(msg, emailBody.Bytes()...)

	// Send the email.
	err = smtp.SendMail(fmt.Sprintf("%v:%v", cfg.Host, cfg.Port), auth, cfg.Username, to, msg)
	if err != nil {
		return err
	}
	return nil
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
