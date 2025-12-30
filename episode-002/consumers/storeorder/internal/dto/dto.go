package dto

type Order struct {
	UserID string `json:"userId"`
	Items  []Item `json:"items"`
}

type Item struct {
	ProductID string `json:"productId"`
	Quantity  int    `json:"quantity"`
}
