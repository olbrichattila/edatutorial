package actions

type OrderCreatedAction struct {
	UserID string      `json:"userID"`
	Email  string      `json:"email"`
	Items  []OrderItem `json:"items"`
}

type OrderItem struct {
	ProductID string `json:"productID"`
	Quantity  uint   `json:"quantity"`
}

type OrderPersistedAction struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}
