package service

import (
	"context"

	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/dto"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/repository"
	"go.uber.org/zap"
)

type ReportService struct {
	Repo   *repository.ReportRepository
	Logger *zap.Logger
}

func NewReportService(repo *repository.ReportRepository, log *zap.Logger) *ReportService {
	return &ReportService{
		Repo:   repo,
		Logger: log,
	}
}

func (s *ReportService) GetItemsReport(ctx context.Context) (*dto.ItemsReport, error) {
	return s.Repo.GetItemsReport(ctx)
}

func (s *ReportService) GetSalesReport(ctx context.Context) (*dto.SalesReport, error) {
	return s.Repo.GetSalesReport(ctx)
}

func (s *ReportService) GetRevenueReport(ctx context.Context) (*dto.RevenueReport, error) {
	return s.Repo.GetRevenueReport(ctx)
}
