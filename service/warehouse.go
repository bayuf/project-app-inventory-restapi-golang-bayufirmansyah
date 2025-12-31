package service

import (
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/dto"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/model"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/repository"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/utils"
	"go.uber.org/zap"
)

type WarehousesService struct {
	Repo   *repository.WarehouseRepository
	Logger *zap.Logger
	Config *utils.Configuration
}

func NewWarehouseService(repo *repository.WarehouseRepository, log *zap.Logger, config *utils.Configuration) *WarehousesService {
	return &WarehousesService{
		Repo:   repo,
		Logger: log,
		Config: config,
	}
}

func (s *WarehousesService) CreateNewWarehouse(data dto.WarehouseAdd) error {
	if err := s.Repo.Create(model.Warehouse{
		Name:     data.Name,
		Location: data.Location,
	}); err != nil {
		return err
	}

	return nil
}

func (s *WarehousesService) GetAllWarehouses(page, limit int) (*[]dto.WarehouseResponse, *dto.Pagination, error) {
	warehouses, total, err := s.Repo.GetAll(page, limit)
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

func (s *WarehousesService) GetById(id int) (*dto.WarehouseResponse, error) {
	warehouse, err := s.Repo.GetById(id)
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
