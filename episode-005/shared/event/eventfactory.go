package event

import (
	"fmt"
	"os"

	"github.com/olbrichattila/edatutorial/shared/event/contracts"
	"github.com/olbrichattila/edatutorial/shared/event/rabbitmq"
)

func New() (contracts.EventManager, error) {
	switch os.Getenv("QUEUE") {
	case "AWS":
		return nil, fmt.Errorf("AWS not implemented")
	case "KAFKA":
		return nil, fmt.Errorf("KAFKA not implemented")
	case "MQ":
		return rabbitmq.New(), nil
	default:
		return rabbitmq.New(), nil
	}
}
