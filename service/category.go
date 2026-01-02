package service

import (
	"context"

	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/dto"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/model"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/repository"
	"go.uber.org/zap"
)

type CategoryService struct {
	Repo   *repository.CategoryRepository
	Logger *zap.Logger
}

func NewCategoryService(repo *repository.CategoryRepository, log *zap.Logger) *CategoryService {
	return &CategoryService{
		Repo:   repo,
		Logger: log,
	}
}

func (s *CategoryService) CreateNewCategory(ctx context.Context, newData dto.CategoryAdd) error {
	return s.Repo.Create(ctx, model.Category{
		Name:        newData.Name,
		Description: newData.Description,
	})
}
