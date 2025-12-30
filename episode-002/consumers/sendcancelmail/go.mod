module producer.example

go 1.25.5

replace github.com/olbrichattila/edatutorial/shared => ../../shared

require github.com/olbrichattila/edatutorial/shared v0.0.0-00010101000000-000000000000

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/go-sql-driver/mysql v1.9.3 // indirect
	github.com/rabbitmq/amqp091-go v1.10.0 // indirect
)
