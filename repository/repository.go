package repository

import (
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/db"
	"go.uber.org/zap"
)

type Repository struct {
	*UserRepository
	*AuthRepository

	*WarehouseRepository
	*RackRepository
	*CategoryRepository
	*ItemRepository

	*SaleRepository
}

func NewRepository(db db.DBExecutor, log *zap.Logger) *Repository {
	return &Repository{
		UserRepository: NewUserRepository(db, log),
		AuthRepository: NewAuthRepository(db, log),

		WarehouseRepository: NewWarehousesRepository(db, log),
		RackRepository:      NewRackRepository(db, log),
		CategoryRepository:  NewCategoryRepository(db, log),
		ItemRepository:      NewItemRepository(db, log),

		SaleRepository: NewSaleRepository(db, log),
	}
}
