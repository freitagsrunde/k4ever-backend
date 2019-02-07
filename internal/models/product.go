package models

import "github.com/jinzhu/gorm"

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
	gorm.Model

	// The name of the product
	//
	// required: true
	// min: 1
	//
	// example: club mate
	Name        string  `json:"name" gorm:"not_null"`
	Price       float64 `gorm:"not_null"`
	Description string  `gorm:"type:text;"`
	Deposit     float64
	Barcode     string
	Image       string
}
