package service

import (
	"context"

	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/dto"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/model"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/repository"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/utils"
	"go.uber.org/zap"
)

type RackService struct {
	Repo   *repository.RackRepository
	Logger *zap.Logger
}

func NewRackService(repo *repository.RackRepository, log *zap.Logger) *RackService {
	return &RackService{
		Repo:   repo,
		Logger: log,
	}
}

func (s *RackService) CreateNewRack(ctx context.Context, newData dto.RackAdd) error {
	if err := s.Repo.Create(ctx, model.Rack{
		WarehouseID: newData.WarehouseId,
		Code:        newData.RackCode,
		Description: newData.Description,
	}); err != nil {
		return err
	}

	return nil
}

func (s *RackService) GetAllRacks(ctx context.Context, page, limit int) (*[]dto.RackResponse, *dto.Pagination, error) {
	racks, total, err := s.Repo.GetAllRacks(ctx, page, limit)
	if err != nil {
		return nil, nil, err
	}

	pagination := dto.Pagination{
		CurrentPage:  page,
		Limit:        limit,
		TotalPages:   utils.TotalPage(limit, int64(total)),
		TotalRecords: total,
	}

	racksRes := []dto.RackResponse{}
	for _, v := range *racks {
		res := dto.RackResponse{
			ID:                v.ID,
			WarehouseName:     v.WarehouseName,
			WarehouseLocation: v.WarehouseLocation,
			Code:              v.Code,
			Description:       v.Description,
			Created_At:        v.Created_At,
			Updated_At:        v.Updated_At,
		}

		racksRes = append(racksRes, res)
	}

	return &racksRes, &pagination, nil
}

func (s *RackService) GetRackById(ctx context.Context, id int) (*dto.RackResponse, error) {
	rack, err := s.Repo.GetRackById(ctx, id)
	if err != nil {
		return nil, err
	}

	return &dto.RackResponse{
		ID:                rack.ID,
		Code:              rack.Code,
		Description:       rack.Description,
		WarehouseName:     rack.WarehouseName,
		WarehouseLocation: rack.WarehouseLocation,
		Created_At:        rack.Created_At,
		Updated_At:        rack.Updated_At,
	}, nil
}

func (s *RackService) DeleteRackById(ctx context.Context, id int) error {
	return s.Repo.Delete(ctx, id)
}

func (s *RackService) UpdateRackById(ctx context.Context, newRack dto.RackUpdate) error {
	return s.Repo.Update(ctx, model.Rack{
		ID:          newRack.ID,
		Code:        newRack.NewRackCode,
		Description: newRack.NewRackDescription,
	})
}
