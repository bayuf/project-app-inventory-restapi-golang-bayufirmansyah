package dto

import (
	"time"

	"github.com/google/uuid"
)

type UserAdd struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	RoleID   int    `json:"role_id" validate:"required,eq=2|eq=3"`
}

type UserUpdate struct {
	ID       uuid.UUID `validate:"required"`
	Name     string    `json:"name" validate:"required"`
	Password string    `json:"password" validate:"required"`
	RoleID   int       `json:"role_id" validate:"required,eq=2|eq=3"`
}

type UserSelfUpdate struct {
	ID       uuid.UUID `validate:"required"`
	Name     string    `json:"name" validate:"required"`
	Password string    `json:"password" validate:"required"`
}

type UserRegister struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserSuspend struct {
	ID      uuid.UUID `validate:"required"`
	Suspend bool      `json:"suspend" validate:"required"`
}

type UserResponse struct {
	ID         string    `json:"user_id,omitempty"`
	Name       string    `json:"name"`
	RoleName   string    `json:"role"`
	Email      string    `json:"email"`
	Created_At time.Time `json:"created_at"`
}
