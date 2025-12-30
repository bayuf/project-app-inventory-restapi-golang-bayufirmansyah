package dto

import (
	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name" validate:"required"`
	Email    string    `json:"email" validate:"required,email"`
	Password string    `json:"password" validate:"required"`
	RoleID   int       `json:"role_id" validate:"required,eq=2|eq=3"`
}

type UserAdd struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	RoleID   int    `json:"role_id" validate:"required,eq=2|eq=3"`
}
