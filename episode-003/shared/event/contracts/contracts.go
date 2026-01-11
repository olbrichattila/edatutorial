package contracts

type CallbackFn func(evt EventManager, message []byte) error

type EventManager interface {
	Publish(topic string, event []byte) error
	Consume(topic, consumerName string, callback CallbackFn) error
}
