package dto

type Order struct {
	UserID string `json:"userId"`
	Email  string `json:"email"`
	Items  []Item `json:"items"`
}

type Item struct {
	ProductID string `json:"productId"`
	Quantity  int    `json:"quantity"`
}

type OrderAction struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}
