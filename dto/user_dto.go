package dto

type UserAdd struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	RoleID   int    `json:"role_id" validate:"required,eq=2|eq=3"`
}

type UserReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserResponse struct {
	ID       string `json:"user_id"`
	Name     string `json:"name"`
	RoleName string `json:"role"`
}
