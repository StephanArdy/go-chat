package dto

import (
	"time"
)

type RegisterDataRequest struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginDataRequest struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

type RegisterDataResponse struct {
	ID        string    `json:"_id"`
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

type LoginDataResponse struct {
	ID        string    `json:"_id"`
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Friends   []string  `json:"friends"`
}
