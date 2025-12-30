package event

import (
	"os"

	"github.com/olbrichattila/edatutorial/shared/event/contracts"
	"github.com/olbrichattila/edatutorial/shared/event/rabbitmq"
)

func New() contracts.EventManager {
	switch os.Getenv("QUEUE") {
	case "AWS":
		panic("not implemented")
	case "KAFKA":
		panic("not implemented")
	case "MQ":
		return rabbitmq.New()
	default:
		return rabbitmq.New()
	}
}
