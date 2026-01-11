package contracts

type OrderRepository interface {
	Cancel(orderID int64) error
}
