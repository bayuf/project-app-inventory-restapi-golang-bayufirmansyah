package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/dto"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/middleware"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/service"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type SaleHandler struct {
	Service *service.SaleService
	Logger  *zap.Logger
	Config  *utils.Configuration
}

func NewSaleHandler(service *service.SaleService, log *zap.Logger, config *utils.Configuration) *SaleHandler {
	return &SaleHandler{
		Service: service,
		Logger:  log,
		Config:  config,
	}
}

func (h *SaleHandler) GetSaleInfo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method != "GET" {
		utils.ResponseFailed(w, http.StatusMethodNotAllowed, "method not allowed", nil)
	}

	strSaleId := chi.URLParam(r, "sale_id")

	uuidSaleId, err := uuid.Parse(strSaleId)
	if err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "id not valid", err.Error())
	}

	sale, err := h.Service.GetSaleDetailById(ctx, uuidSaleId)
	if err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "cant get sale detail", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "success", sale)
}

func (h *SaleHandler) GetSaleInfoStaff(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method != "GET" {
		utils.ResponseFailed(w, http.StatusMethodNotAllowed, "method not allowed", nil)
	}

	authUser, ok := middleware.GetAuthUser(r)
	if !ok {
		utils.ResponseFailed(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	strSaleId := chi.URLParam(r, "sale_id")

	uuidSaleId, err := uuid.Parse(strSaleId)
	if err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "id not valid", err.Error())
	}

	sale, err := h.Service.GetStaffSaleDetailById(ctx, uuidSaleId, authUser.UserID)
	if err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "cant get sale detail", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "success", sale)
}

func (h *SaleHandler) InsertNewSale(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authUser, ok := middleware.GetAuthUser(r)
	if !ok {
		utils.ResponseFailed(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	if r.Method != "POST" {
		utils.ResponseFailed(w, http.StatusMethodNotAllowed, "method not allowed", nil)
	}

	Sale := dto.NewSale{}
	if err := json.NewDecoder(r.Body).Decode(&Sale); err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "invalid input format", err)
		return
	}

	// validate
	messageInvalid, err := utils.ValidateInput(&Sale)
	if err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "invalid input data", messageInvalid)
		return
	}

	saleInfo, err := h.Service.NewSaleTX(ctx, Sale, authUser.UserID)
	if err != nil {
		h.Logger.Error("failed to create new item", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "failed to create new item", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusCreated, "success", saleInfo)
}

func (h *SaleHandler) GetAllSales(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method != "GET" {
		utils.ResponseFailed(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}

	authUser, ok := middleware.GetAuthUser(r)
	if !ok {
		utils.ResponseFailed(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	page, err := strconv.Atoi(r.URL.Query().Get("page"))

	if err != nil {
		h.Logger.Info("invalid page :", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "invalid page", err)
		return
	}

	// page limit
	limit := h.Config.PageLimit

	items, pagination, err := h.Service.GetAllSales(ctx, page, limit, authUser.Role, authUser.UserID)
	if err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "cant get all sales", err.Error())
		return
	}

	utils.ResponsePagination(w, http.StatusOK, "success", items, *pagination)
}

func (h *SaleHandler) UpdateSaleStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method != "PATCH" {
		utils.ResponseFailed(w, http.StatusMethodNotAllowed, "method not allowed", nil)
	}

	status := dto.NewSaleStatus{}
	if err := json.NewDecoder(r.Body).Decode(&status); err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "invalid input format", err)
		return
	}

	// validate
	messageInvalid, err := utils.ValidateInput(&status)
	if err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "invalid input data", messageInvalid)
		return
	}

	strSaleId := chi.URLParam(r, "sale_id")
	saleId, err := uuid.Parse(strSaleId)
	if err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "invalid input data", err.Error())
		return
	}

	newStatus := status.Status

	if err := h.Service.UpdateSaleStatus(ctx, saleId, newStatus); err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "failed update sale status", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "succes", nil)
}

func (h *SaleHandler) DeleteSale(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method != "DELETE" {
		utils.ResponseFailed(w, http.StatusMethodNotAllowed, "method not allowed", nil)
	}

	strSaleId := chi.URLParam(r, "sale_id")
	saleId, err := uuid.Parse(strSaleId)
	if err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "invalid input data", err.Error())
		return
	}

	if err := h.Service.DeteleSale(ctx, saleId); err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "failed to delete sale", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "succes", nil)
}
