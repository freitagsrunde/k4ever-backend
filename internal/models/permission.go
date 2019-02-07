package models

import "github.com/jinzhu/gorm"

// A generic Permission
//
// swagger:model
type Permission struct {
	// The id for the product
	//
	// required: true
	// unique: true
	// min: 1
	//
	// example: 1
	gorm.Model

	// The name of the permission
	//
	// required: true
	//
	// example: AdminPermission
	Name string `gorm:"not_null"`

	// The description of the permission
	//
	// example: Can do everything
	Description string
}
