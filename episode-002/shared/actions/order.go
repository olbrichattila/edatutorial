package actions

type OrderSentAction struct {
	UserID string      `json:"userId"`
	Email  string      `json:"email"`
	Items  []OrderItem `json:"items"`
}

type OrderItem struct {
	ProductID string `json:"productId"`
	Quantity  int    `json:"quantity"`
}

type OrderStoredAction struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}
