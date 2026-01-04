package service

import (
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/db"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/repository"
	"go.uber.org/zap"
)

type Service struct {
	*UserService
	*AuthService

	*WarehousesService
	*RackService
	*CategoryService
	*ItemService

	*SaleService
	*ReportService
}

func NewService(repo *repository.Repository, log *zap.Logger, tx db.TxManager) *Service {
	return &Service{
		UserService: NewUserService(repo.UserRepository, log),
		AuthService: NewAuthService(repo.AuthRepository, log),

		WarehousesService: NewWarehouseService(repo.WarehouseRepository, log),
		RackService:       NewRackService(repo.RackRepository, log),
		CategoryService:   NewCategoryService(repo.CategoryRepository, log),
		ItemService:       NewItemService(repo.ItemRepository, log, tx),

		SaleService:   NewSaleService(repo.SaleRepository, log, tx),
		ReportService: NewReportService(repo.ReportRepository, log),
	}
}
