package model

type User struct {
	ModelUser
	Email    string
	Password string
	RoleID   int
	// IsActive bool
}
