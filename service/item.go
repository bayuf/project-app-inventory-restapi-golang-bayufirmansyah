package service

import (
	"context"
	"fmt"

	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/dto"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/model"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/repository"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/utils"
	"github.com/google/uuid"
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

func (s *ItemService) InputNewItem(ctx context.Context, new dto.ItemAdd) error {
	return s.Repo.Create(ctx, model.Item{
		ID:         uuid.New(),
		Name:       new.Name,
		CategoryID: new.CategoryId,
		RackID:     new.RackId,
		Stock:      new.Stock,
		MinStock:   new.MinStock,
		Price:      new.Price,
	})
}

func (s *ItemService) GetAllItems(ctx context.Context, page, limit int) (*[]dto.ItemResponse, *dto.Pagination, error) {
	items, total, err := s.Repo.GetAllItems(ctx, page, limit)
	if err != nil {
		return nil, nil, err
	}

	pagination := dto.Pagination{
		CurrentPage:  page,
		Limit:        limit,
		TotalPages:   utils.TotalPage(limit, int64(total)),
		TotalRecords: total,
	}

	itemsRes := []dto.ItemResponse{}
	for _, v := range *items {
		res := dto.ItemResponse{
			ID:         v.ID,
			Name:       v.Name,
			CategoryId: v.CategoryID,
			Category:   v.CategoryName,
			RackID:     v.RackID,
			Rack:       v.RackCode,
			Location:   fmt.Sprintf("%s %s", v.RackName, v.Location),
			Stock:      v.Stock,
			Price:      v.Price,
			Created_At: v.Created_At,
			Updated_At: v.Updated_At,
		}

		itemsRes = append(itemsRes, res)
	}

	return &itemsRes, &pagination, nil
}

func (s *ItemService) GetItem(ctx context.Context, id uuid.UUID) (*dto.ItemResponse, error) {
	v, err := s.Repo.GetItem(ctx, id)
	if err != nil {
		return nil, err
	}

	itemsRes := dto.ItemResponse{
		ID:         v.ID,
		Name:       v.Name,
		CategoryId: v.CategoryID,
		Category:   v.CategoryName,
		RackID:     v.RackID,
		Rack:       v.RackCode,
		Location:   fmt.Sprintf("%s %s", v.RackName, v.Location),
		Stock:      v.Stock,
		Price:      v.Price,
		Created_At: v.Created_At,
		Updated_At: v.Updated_At,
	}

	return &itemsRes, nil
}

func (s *ItemService) DeleteItem(ctx context.Context, id uuid.UUID) error {
	return s.Repo.Delete(ctx, id)
}
