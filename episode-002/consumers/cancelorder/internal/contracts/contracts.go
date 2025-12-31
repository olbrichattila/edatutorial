package contracts

type OrderRepository interface {
	Cancel(orderId int64) error
}
