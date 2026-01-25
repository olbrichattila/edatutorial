package actions

// Episode 004 Separate payment actions
type PaymentSucceededAction struct {
	OrderUUID string `json:"orderUUID"`
	Email     string `json:"email"`
}

type PaymentFailedActon struct {
	OrderUUID string `json:"orderUUID"`
	Email     string `json:"email"`
}
