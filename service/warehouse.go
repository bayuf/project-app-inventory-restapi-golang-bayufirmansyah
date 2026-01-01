package service

import (
	"context"

	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/dto"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/model"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/repository"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/utils"
	"go.uber.org/zap"
)

type WarehousesService struct {
	Repo   *repository.WarehouseRepository
	Logger *zap.Logger
}

func NewWarehouseService(repo *repository.WarehouseRepository, log *zap.Logger) *WarehousesService {
	return &WarehousesService{
		Repo:   repo,
		Logger: log,
	}
}

func (s *WarehousesService) CreateNewWarehouse(ctx context.Context, data dto.WarehouseAdd) error {
	if err := s.Repo.Create(ctx, model.Warehouse{
		Name:     data.Name,
		Location: data.Location,
	}); err != nil {
		return err
	}

	return nil
}

func (s *WarehousesService) GetAllWarehouses(ctx context.Context, page, limit int) (*[]dto.WarehouseResponse, *dto.Pagination, error) {
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
	warehousesRes := []dto.WarehouseResponse{}
	for _, v := range *warehouses {
		res := dto.WarehouseResponse{
			ID:         v.ID,
			Name:       v.Name,
			Location:   v.Location,
			Created_at: v.Created_at,
			Updated_at: v.Updated_at,
		}

		warehousesRes = append(warehousesRes, res)
	}

	return &warehousesRes, &pagination, nil
}

func (s *WarehousesService) GetById(ctx context.Context, id int) (*dto.WarehouseResponse, error) {
	warehouse, err := s.Repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	return &dto.WarehouseResponse{
		ID:         warehouse.ID,
		Name:       warehouse.Name,
		Location:   warehouse.Location,
		Created_at: warehouse.Created_at,
		Updated_at: warehouse.Updated_at,
	}, nil
}

func (s *WarehousesService) UpdateWarehouse(ctx context.Context, newData dto.Warehouse) error {
	if err := s.Repo.Update(ctx, model.Warehouse{
		ID:       newData.ID,
		Name:     newData.Name,
		Location: newData.Location,
	}); err != nil {
		return err
	}

	return nil
}

func (s *WarehousesService) DeleteWarehouseById(ctx context.Context, id int) error {
	return s.Repo.Delete(ctx, id)
}
