package repository

import (
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/db"
	"go.uber.org/zap"
)

type ReportRepository struct {
	DB     db.DBExecutor
	Logger *zap.Logger
}

func NewReportRepository(db db.DBExecutor, log *zap.Logger) *ReportRepository {
	return &ReportRepository{
		DB:     db,
		Logger: log,
	}
}
