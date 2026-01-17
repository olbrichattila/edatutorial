package contracts

type OrderRepository interface {
	Cancel(orderUUID string) error
}
