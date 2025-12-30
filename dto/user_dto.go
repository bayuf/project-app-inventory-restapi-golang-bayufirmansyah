package dto

import (
	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name" validate:"required"`
	Email    string    `json:"email" validate:"required,email"`
	Password string    `json:"password" validate:"required"`
	RoleID   int       `json:"role_id" validate:"required"`
	// IsActive bool      `json:"is_active"`
}

type UserAdd struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	RoleID   int    `json:"role_id" validate:"required"`
}
