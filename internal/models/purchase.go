package models

// swagger:model
type PurchaseItem struct {
	ModelTimes

	// The amount of products bought
	Amount int `json:"amount"`

	// Information about the bought product
	PurchaseItemInformation

	ProductID  uint `json:"product_id"`
	PurchaseID uint `json:"-"`
}

type PurchaseItemInformation struct {
	// The name of the product
	//
	// required: true
	// min: 1
	//
	// example: club mate
	Name string `json:"name" gorm:"not_null"`

	// The price of the product
	//
	// required: true
	//
	// example: 1.00
	Price float64 `json:"price" gorm:"not_null"`

	// The description of the product
	//
	// required: false
	Description string `json:"description" gorm:"type:text;"`

	// An optional deposit meant for drinks
	//
	// required: false
	Deposit float64 `json:"deposit"`

	// The barcode of the product
	//
	// required: false
	Barcode string `json:"barcode" gorm:"default: null"`

	// Currently the path to the image (tbi)
	//
	// required: false
	Image string `json:"image"`
}

// swagger:model
type Purchase struct {
	ModelTimes

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
