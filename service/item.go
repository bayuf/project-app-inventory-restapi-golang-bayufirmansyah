package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/db"
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
	TX     db.TxManager
}

func NewItemService(repo *repository.ItemRepository, log *zap.Logger, tx db.TxManager) *ItemService {
	return &ItemService{
		Repo:   repo,
		Logger: log,
		TX:     tx,
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

func (s *ItemService) UpdateItem(ctx context.Context, new dto.ItemUpdate) error {
	return s.Repo.Update(ctx, model.Item{
		ID:         new.ID,
		Name:       new.Name,
		CategoryID: new.CategoryId,
		RackID:     new.RackId,
		MinStock:   new.MinStock,
		Price:      new.Price,
	})
}

func (s *ItemService) DeleteItem(ctx context.Context, id uuid.UUID) error {
	return s.Repo.Delete(ctx, id)
}

func (s *ItemService) StockAdjustment(ctx context.Context, data dto.StockAdjustment) error {
	tx, err := s.TX.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// init repo with tx
	repo := repository.NewItemRepository(tx, s.Logger)

	itemStock, err := repo.GetItemStock(ctx, data.ItemID)
	if err != nil {
		return err
	}

	// Update Stock
	if data.Change < 0 {
		if itemStock < data.Change {
			return errors.New("insufficient stock")
		}
		if err := repo.UpdateStock(ctx, data.ItemID, data.Change); err != nil {
			return err
		}

	} else {
		if err := repo.UpdateStock(ctx, data.ItemID, data.Change); err != nil {
			return err
		}
	}

	if err := repo.StockAdjustments(ctx, dto.StockAdjustment{
		ID:     uuid.New(),
		ItemID: data.ItemID,
		UserID: data.UserID,
		Change: data.Change,
		Reason: data.Reason,
	}); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		s.Logger.Error("failed to update stock ajustment", zap.Error(err))
		return err
	}

	s.Logger.Info("update stock succes", zap.Any("Item ID", data.ItemID))
	return nil
}
