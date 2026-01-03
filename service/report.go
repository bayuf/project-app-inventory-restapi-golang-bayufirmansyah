package service

import (
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
