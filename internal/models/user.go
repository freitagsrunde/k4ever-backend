package models

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	UserName    string `json:"name" gorm:"not_null;unique;"`
	DisplayName string `json:"display_name" gorm:"not_null;"`
	Balance     float64
	Permissions []Permission `gorm:"many2many:user_permissions;"`
}
