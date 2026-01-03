package handler

import (
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/service"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/utils"
	"go.uber.org/zap"
)

type ReportHandler struct {
	Service *service.ReportService
	Logger  *zap.Logger
	Config  *utils.Configuration
}

func NewReportHandler(service *service.ReportService, log *zap.Logger, config *utils.Configuration) *ReportHandler {
	return &ReportHandler{
		Service: service,
		Logger:  log,
		Config:  config,
	}
}
