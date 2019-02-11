package models

// A generic product
//
// A generic product with some values
//
// swagger:model
type Product struct {
	// The id for the product
	//
	// required: true
	// unique: true
	// min: 1
	//
	// example: 1
	Model

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
	Barcode string `json:"barcode"`

	// Currently the path to the image (tbi)
	//
	// required: false
	Image string `json:"image"`
}
