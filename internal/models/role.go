package models

type Role struct {
	Model

	Name string `gorm:"not_null"`

	Description string

	Permissions []Permission `json:"permissions" gorm:"many2many:roles_permissions;"`
}
