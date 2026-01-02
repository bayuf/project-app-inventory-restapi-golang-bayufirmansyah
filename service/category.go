package service

import (
	"context"

	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/dto"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/model"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/repository"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/utils"
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

func (s *CategoryService) GetAllCategories(ctx context.Context, page, limit int) (*[]dto.CategoryResponse, *dto.Pagination, error) {
	warehouses, total, err := s.Repo.GetAll(ctx, page, limit)
	if err != nil {
		return nil, nil, err
	}

	pagination := dto.Pagination{
		CurrentPage:  page,
		Limit:        limit,
		TotalPages:   utils.TotalPage(limit, int64(total)),
		TotalRecords: total,
	}
	categoriesRes := []dto.CategoryResponse{}
	for _, v := range *warehouses {
		res := dto.CategoryResponse{
			ID:          v.ID,
			Name:        v.Name,
			Description: v.Description,
			Created_At:  v.Created_At,
			Updated_At:  v.Updated_At,
		}

		categoriesRes = append(categoriesRes, res)
	}

	return &categoriesRes, &pagination, nil
}

func (s *CategoryService) GetCategoryById(ctx context.Context, id int) (*dto.CategoryResponse, error) {
	category, err := s.Repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	return &dto.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		Created_At:  category.Created_At,
		Updated_At:  category.Updated_At,
	}, nil

}

func (s *CategoryService) UpdateCategory(ctx context.Context, new dto.CategoryUpdate) error {
	if err := s.Repo.Update(ctx, model.Category{
		ID:          new.ID,
		Name:        new.Name,
		Description: new.Description,
	}); err != nil {
		return err
	}

	return nil
}

func (s *CategoryService) DeleteCategory(ctx context.Context, id int) error {
	return s.Repo.Delete(ctx, id)
}
