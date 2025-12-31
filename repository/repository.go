package repository

import (
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/db"
	"go.uber.org/zap"
)

type Repository struct {
	UserRepository *UserRepository
	AuthRepository *AuthRepository

	WarehouseRepository *WarehouseRepository
}

func NewRepository(db db.PgxIface, log *zap.Logger) *Repository {
	return &Repository{
		UserRepository: NewUserRepository(db, log),
		AuthRepository: NewAuthRepository(db, log),

		WarehouseRepository: NewWarehousesRepository(db, log),
	}
}
