package models

import "github.com/jinzhu/gorm"

type Product struct {
	gorm.Model
	Name        string  `json:"name" gorm:"not_null"`
	Price       float64 `gorm:"not_null"`
	Description string  `gorm:"type:text;"`
	Deposit     float64
	Barcode     string
	Image       string
}
