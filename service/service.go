package service

import (
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/repository"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/utils"
	"go.uber.org/zap"
)

type Service struct {
	UserService *UserService
}

func NewService(repo *repository.Repository, log *zap.Logger, config utils.Configuration) *Service {
	return &Service{
		UserService: NewUserService(repo.UserRepository, log),
	}
}
