package notification

import (
	"fmt"
	"net/smtp"
	"strings"

	"github.com/olbrichattila/edatutorial/shared/config"
)

func SendEmail(emailAddress, subject, emailBody string) error {
	// Sanitize inputs to prevent header injection
	emailAddress = sanitizeHeader(emailAddress)
	subject = sanitizeHeader(subject)

	addr := fmt.Sprintf("%s:%s",
		config.MailSmtpHost(),
		config.MailSmtpPort(),
	)

	msg := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/html; charset=\"UTF-8\"\r\n\r\n%s",
		config.MailFrom(),
		emailAddress,
		subject,
		emailBody,
	)

	return smtp.SendMail(
		addr,
		nil, // NO AUTH
		config.MailFrom(),
		[]string{emailAddress},
		[]byte(msg),
	)
}

// sanitizeHeader removes CRLF characters to prevent header injection
func sanitizeHeader(input string) string {
	// Remove carriage return and line feed characters
	input = strings.ReplaceAll(input, "\r", "")
	input = strings.ReplaceAll(input, "\n", "")
	return input
}
