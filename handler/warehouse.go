package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/dto"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/service"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/utils"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type WarehouseHandler struct {
	service *service.WarehousesService
	Logger  *zap.Logger
	Config  *utils.Configuration
}

func NewWarehouseHandler(service *service.WarehousesService, log *zap.Logger, config *utils.Configuration) *WarehouseHandler {
	return &WarehouseHandler{
		service: service,
		Logger:  log,
		Config:  config,
	}
}

func (h *WarehouseHandler) CreateWarehouse(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method != "POST" {
		utils.ResponseFailed(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}

	newWarehouse := dto.WarehouseAdd{}
	if err := json.NewDecoder(r.Body).Decode(&newWarehouse); err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "invalid input format", err)
		return
	}

	// validate
	messageInvalid, err := utils.ValidateInput(&newWarehouse)
	if err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "invalid input data", messageInvalid)
		return
	}

	if err := h.service.CreateNewWarehouse(ctx, newWarehouse); err != nil {
		h.Logger.Error("failed to create new warehouse", zap.Error(err))
		utils.ResponseFailed(w, http.StatusInternalServerError, "failed to create new warehouse", err)
		return
	}

	utils.ResponseSuccess(w, http.StatusCreated, "success", nil)
}

func (h *WarehouseHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method != "GET" {
		utils.ResponseFailed(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	log.Println(page)

	if err != nil {
		h.Logger.Info("invalid page :", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "invalid page", err)
		return
	}

	// page limit
	limit := h.Config.PageLimit

	warehouses, pagination, err := h.service.GetAllWarehouses(ctx, page, limit)
	if err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "cant get all warehouses", err)
		return
	}

	utils.ResponsePagination(w, http.StatusOK, "success", warehouses, *pagination)
}

func (h *WarehouseHandler) GetById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method != "GET" {
		utils.ResponseFailed(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}

	strWarehouse := chi.URLParam(r, "warehouse_id")
	warehouseId, err := strconv.Atoi(strWarehouse)
	if err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "id invalid", nil)
		return
	}

	warehouse, err := h.service.GetById(ctx, warehouseId)
	if err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "cant find warehouse", errors.New("warehouse not found"))
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "success", warehouse)

}

func (h *WarehouseHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method != "PUT" {
		utils.ResponseFailed(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}

	strWarehouse := chi.URLParam(r, "warehouse_id")
	warehouseId, err := strconv.Atoi(strWarehouse)
	if err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "id invalid", nil)
		return
	}

	newWarehouse := dto.Warehouse{
		ID: warehouseId,
	}
	if err := json.NewDecoder(r.Body).Decode(&newWarehouse); err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "invalid input format", err)
		return
	}

	if err := h.service.UpdateWarehouse(ctx, newWarehouse); err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "warehouse not found", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "success", nil)
}

func (h *WarehouseHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method != "DELETE" {
		utils.ResponseFailed(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}

	strWarehouse := chi.URLParam(r, "warehouse_id")
	warehouseId, err := strconv.Atoi(strWarehouse)
	if err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "id invalid", nil)
		return
	}

	if err := h.service.DeleteWarehouseById(ctx, warehouseId); err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "warehouse not found", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "success", nil)
}
