package service

import (
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/repository"
	"go.uber.org/zap"
)

type Service struct {
	*UserService
	*AuthService

	*WarehousesService
	*RackService
}

func NewService(repo *repository.Repository, log *zap.Logger) *Service {
	return &Service{
		UserService: NewUserService(repo.UserRepository, log),
		AuthService: NewAuthService(repo.AuthRepository, log),

		WarehousesService: NewWarehouseService(repo.WarehouseRepository, log),
		RackService:       NewRackService(repo.RackRepository, log),
	}
}
