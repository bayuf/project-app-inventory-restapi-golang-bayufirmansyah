package handler

import (
	"net/http"

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

func (h *ReportHandler) GetItemsReport(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method != "GET" {
		utils.ResponseFailed(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}

	report, err := h.Service.GetItemsReport(ctx)
	if err != nil {
		utils.ResponseFailed(w, http.StatusInternalServerError, "failed", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "succes", report)
}

func (h *ReportHandler) GetSalesReport(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method != "GET" {
		utils.ResponseFailed(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}

	report, err := h.Service.GetSalesReport(ctx)
	if err != nil {
		utils.ResponseFailed(w, http.StatusInternalServerError, "failed", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "succes", report)
}

func (h *ReportHandler) GetRevenueReport(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method != "GET" {
		utils.ResponseFailed(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}

	report, err := h.Service.GetRevenueReport(ctx)
	if err != nil {
		utils.ResponseFailed(w, http.StatusInternalServerError, "failed", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "succes", report)
}
