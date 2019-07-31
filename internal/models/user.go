package models

type UserDgraph struct {
	User
	Type bool `json:"user"`
}

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

	// A list of user permission
	Permissions []Permission `json:"permissions" gorm:"many2many:user_permissions;"`

	// A list of purchases made by the user
	Histories []History `json:"histories"`
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
