package config

import "os"

const (
	rabbitMQURLKey = "RABBIT_URL"

	dbHostKey     = "DB_HOST"
	dbPortKey     = "DB_PORT"
	dbDatabaseKey = "DB_DATABASE"
	dbUsernameKey = "DB_USERNAME"
	dbPasswordKey = "DB_PASSWORD"

	mailSmtpUserNameKey = "SMTP_USER_NAME"
	mailSmtpPasswordKey = "SMTP_PASSWORD"
	mailFromKey         = "SMTP_MAIL_FROM"
	mailHostKey         = "SMTP_HOST"
	mailPortKey         = "SMTP_PORT"
)

const (
	defaultRabbitMQURL = "amqp://dev:dev@localhost:5672/"

	defaultDbHost     = "127.0.0.1"
	defaultDbPort     = "3306"
	defaultDbDatabase = "eda"
	defaultDbUsername = "eda"
	defaultDbPassword = "eda"

	defaultMailFrom     = "testcompany@test.com"
	defaultMailHost     = "127.0.0.1"
	defaultMailPort     = "1025"
	defaultSmtpUserName = "mailtrap"
	defaultSmtpPassword = "mailtrap"
)

func RabbitMQURL() string {
	return getFromEnv(rabbitMQURLKey, defaultRabbitMQURL)
}

func DBUsername() string {
	return getFromEnv(dbUsernameKey, defaultDbUsername)
}

func DBPassword() string {
	return getFromEnv(dbPasswordKey, defaultDbPassword)
}

func DBHost() string {
	return getFromEnv(dbHostKey, defaultDbHost)
}

func DBPort() string {
	return getFromEnv(dbPortKey, defaultDbPort)
}

func DBDatabase() string {
	return getFromEnv(dbDatabaseKey, defaultDbDatabase)
}

// Email
func MailSmtpUserName() string {
	return getFromEnv(mailSmtpUserNameKey, defaultSmtpUserName)
}

func MailSmtpPassword() string {
	return getFromEnv(mailSmtpPasswordKey, defaultSmtpPassword)
}

func MailSmtpHost() string {
	return getFromEnv(mailHostKey, defaultMailHost)
}

func MailSmtpPort() string {
	return getFromEnv(mailPortKey, defaultMailPort)
}

func MailFrom() string {
	return getFromEnv(mailFromKey, defaultMailFrom)
}

func getFromEnv(key, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}

	return val
}
