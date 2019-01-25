package models

import "github.com/jinzhu/gorm"

type Item struct {
	gorm.Model
	Amount     int
	Product    Product
	ProductID  uint
	PurchaseID uint
}

type Purchase struct {
	gorm.Model
	Amount float64
	Items  []Item
	UserID uint
}
