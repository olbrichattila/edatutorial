package notification

import (
	"fmt"
	"net/smtp"

	"github.com/olbrichattila/edatutorial/shared/config"
)

func SendEmail(emailAddress, subject, emailBody string) error {
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
