package model

import "time"

type User struct {
	ID           int       `json:"id,omitempty"`
	Login        string    `json:"login,omitempty"`
	Email        string    `json:"email,omitempty"`
	Password     string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	LastActivity time.Time `json:"last_activity"`
	IsActive     bool      `json:"is_active,omitempty"`
	PhoneNumber  string    `json:"phone_number,omitempty"`
	Address      string    `json:"address,omitempty"`
	Role         string    `json:"role,omitempty"`
}
