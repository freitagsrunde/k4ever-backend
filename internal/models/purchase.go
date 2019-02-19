package models

// swagger:model
type PurchaseItem struct {
	Model

	// The amount of products bought
	Amount int `json:"amount"`

	// Information about the bought product
	Product Product `json:"product"`

	ProductID  uint `json:"-"`
	PurchaseID uint `json:"-"`
}

// swagger:model
type Purchase struct {
	Model

	// The total amount of the purchase
	Total float64 `json:"total"`

	// A list of all items from the purchase
	Items []PurchaseItem `json:"items"`

	UserID uint `json:"-"`
}

// swagger:model
type PurchaseArray struct {
	// An array of purchases
	Purchases []Purchase `json:"purchases"`
}
