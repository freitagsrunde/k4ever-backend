package models

import "github.com/jinzhu/gorm"

type Product struct {
	gorm.Model
	ProductName        string  `json:"name" gorm:"not_null"`
	ProductPrice       float64 `gorm:"not_null"`
	ProductDescription string  `gorm:"type:text;"`
	ProductDeposit     float64
	Barcode            string
	ProductImage       string
}
