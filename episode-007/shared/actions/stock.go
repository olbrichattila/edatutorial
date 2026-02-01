package actions

// Episode 005  handle business scenario, the product is out of stock
type OutOfStockAction struct {
	OrderUUID string `json:"orderUUID"`
	Email     string `json:"email"`
}
