module producer.example

go 1.25.5

replace eda.event => ../../shared/event

require eda.event v0.0.0-00010101000000-000000000000

require github.com/rabbitmq/amqp091-go v1.10.0 // indirect
