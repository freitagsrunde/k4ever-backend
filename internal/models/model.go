package models

import (
	"database/sql"
	"encoding/json"
	"time"
)

// gorm.Model definition
type Model struct {
	Uid       string     `json:"uid"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`
}

// gorm.Model
type ModelTimes struct {
	Uid       string     `json:"uid"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
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

type NullString struct {
	sql.NullString
}

func (ns NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("\"\""), nil
	}
	return json.Marshal(ns.String)
}

func (ns *NullString) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &ns.String)
	ns.Valid = (err == nil)
	return err
}
