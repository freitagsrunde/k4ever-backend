package models

// Make custom type to add custom marshal interface to it

// A generic User
//
// swagger:model
type User struct {
	Model

	// The username
	UserName string `json:"name" gorm:"not_null;unique;"`

	Password string `json:"-" gorm:"not_null"`

	// The displayname
	DisplayName string `json:"display_name" gorm:"not_null;"`

	// The users current balance
	Balance float64 `json:"balance"`

	// The Role the user has (as an integer)
	Role int `json:"role"`

	// A list of user permission
	Permissions []Permission `json:"permissions,omitempty" gorm:"many2many:user_permissions;"`

	// A list of purchases made by the user
	Histories []History `json:"histories,omitempty"`
}

// This is just for swagger

// The returned token
//
// swagger:model Token
type Token struct {
	Code   string `json:"code"`
	Expire string `json:"expire"`
	Token  string `json:"token"`
}
