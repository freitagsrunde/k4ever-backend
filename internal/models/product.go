package models

import (
	"fmt"
	"time"
)

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

	ProductInformation
}

type ProductInformation struct {
	// The name of the product
	//
	// required: true
	// min: 1
	//
	// example: club mate
	Name string `json:"name" gorm:"not_null;unique"`

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
	Barcode NullString `json:"barcode" gorm:"unique;default: null"`

	// Currently the path to the image (tbi)
	//
	// required: false
	Image string `json:"image"`

	TimesBoughtTotal int `json:"times_bought_total" gorm:"-"`

	TimesBought int `json:"times_bought" gorm:"-"`

	LastBought *time.Time `json:"last_bought" gorm:"-"`
	// A flag to show if the product is currently buyable
	//
	// required: false
	Disabled bool `json:"disabled" gorm:"default:false"`

	// A flag to determine wether the object should be displayed at all
	//
	// reqired: false
	Hidden bool `json:"hidden" gorm:"default:false"`
}

func (p Product) String() string {
	return fmt.Sprintf("%s_%d", p.Name, p.ID)
}
