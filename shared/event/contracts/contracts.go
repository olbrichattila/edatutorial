package contracts

type EventManager interface {
	Publish(queName string, event []byte) error
	Consume(consumerName string) error
}
