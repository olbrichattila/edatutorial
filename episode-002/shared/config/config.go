package config

import "os"

const (
	rabbitMqURLKey = "RABBIT_URL"

	dbHostKey     = "DB_HOST"
	dbPortKey     = "DB_PORT"
	dbDatabaseKey = "DB_DATABASE"
	dbUsernameKey = "DB_USERNAME"
	dbPasswordKey = "DB_PASSWORD"
)

const (
	defaultRabbitMqURL = "amqp://dev:dev@localhost:5672/"

	defaultDbHost     = "127.0.0.1"
	defaultDbPort     = "3306"
	defaultDbDatabase = "eda"
	defaultDbUsername = "eda"
	defaultDbPassword = "eda"
)

func RabbitMqURL() string {
	return getFromEnv(rabbitMqURLKey, defaultRabbitMqURL)
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

func getFromEnv(key, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}

	return val
}
