package models

// swagger:model
type PurchaseItem struct {
	ModelTimes

	// The amount of products bought
	Amount int `json:"amount"`

	// Information about the bought product
	PurchaseItemInformation

	ProductID uint `json:"product_id" sql:"default: null"`
	HistoryID uint `json:"-"`
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
type History struct {
	ModelTimes

	// The total amount of the purchase or balance update
	Total float64 `json:"total"`

	// The type of the history item
	Type string `json:"type"`

	// A list of all items from the purchase
	Items []PurchaseItem `json:"items,omitempty"`

	// The recipient of the tranfer
	Recipient string `json:"recipient,omitempty"`

	UserID uint `json:"-"`
}

const PurchaseHistory = "purchase"
const BalanceHistory = "balance"
const TransferHistory = "transfer"

// swagger:model
type HistoryArray struct {
	// An array of purchases
	Histories []History `json:"purchases"`
}
