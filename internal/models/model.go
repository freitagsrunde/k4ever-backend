package models

import "time"

// gorm.Model definition
type Model struct {
	ID        uint       `json:"id" gorm:"primary_key"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`
}

// gorm.Model
type ModelTimes struct {
	ID        uint       `json:"id" gorm:"primary_key"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-"`
}

// swagger:parameters getProducts getUsers
type DefaultParams struct {
	// in: query
	// required: false
	SortBy string `json:"sort_by"`

	// in: query
	// required: false
	Order string `json:"order"`

	// in: query
	// required: false
	Offset int `json:"offset"`

	// in: query
	// required: false
	Limit int `json:"limit"`
}
