package service

import (
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/repository"
	"go.uber.org/zap"
)

type ItemService struct {
	Repo   *repository.ItemRepository
	Logger *zap.Logger
}

func NewItemService(repo *repository.ItemRepository, log *zap.Logger) *ItemService {
	return &ItemService{
		Repo:   repo,
		Logger: log,
	}
}
