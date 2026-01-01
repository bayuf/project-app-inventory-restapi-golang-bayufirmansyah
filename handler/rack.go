package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/dto"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/service"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/utils"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type RackHandler struct {
	Service *service.RackService
	Logger  *zap.Logger
	Config  *utils.Configuration
}

func NewRackHandler(service *service.RackService, log *zap.Logger, config *utils.Configuration) *RackHandler {
	return &RackHandler{
		Service: service,
		Logger:  log,
		Config:  config,
	}
}

func (h *RackHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method != "POST" {
		utils.ResponseFailed(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}

	rack := dto.RackAdd{}
	if err := json.NewDecoder(r.Body).Decode(&rack); err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "invalid input format", err.Error())
		return
	}

	// validate
	messageInvalid, err := utils.ValidateInput(&rack)
	if err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "invalid input data", messageInvalid)
		return
	}

	if err := h.Service.CreateNewRack(ctx, rack); err != nil {
		h.Logger.Error("failed to create new rack", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "failed to create new rack", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusCreated, "success", nil)
}

func (h *RackHandler) GetAllRacks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method != "GET" {
		utils.ResponseFailed(w, http.StatusMethodNotAllowed, "method not allowed", nil)
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

	racks, pagination, err := h.Service.GetAllRacks(ctx, page, limit)
	if err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "cant get all racks", err)
		return
	}

	utils.ResponsePagination(w, http.StatusOK, "success", racks, *pagination)
}

func (h *RackHandler) GetRackById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method != "GET" {
		utils.ResponseFailed(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}

	strRack := chi.URLParam(r, "rack_id")
	rackId, err := strconv.Atoi(strRack)
	if err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "id invalid", nil)
		return
	}

	rack, err := h.Service.GetRackById(ctx, rackId)
	if err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "cant find rack", errors.New("rack not found"))
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "success", rack)
}

func (h *RackHandler) DeleteRack(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method != "DELETE" {
		utils.ResponseFailed(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}

	strRack := chi.URLParam(r, "rack_id")
	rackId, err := strconv.Atoi(strRack)
	if err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "id invalid", nil)
		return
	}

	err = h.Service.DeleteRackById(ctx, rackId)
	if err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "cant find rack", errors.New("rack not found"))
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "success", nil)
}

func (h *RackHandler) UpdateRack(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != "PUT" {
		utils.ResponseFailed(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}

	strRack := chi.URLParam(r, "rack_id")
	rackId, err := strconv.Atoi(strRack)
	if err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "id invalid", nil)
		return
	}

	newRack := dto.RackUpdate{ID: rackId}
	if err := json.NewDecoder(r.Body).Decode(&newRack); err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "invalid input format", err.Error())
		return
	}

	// validate
	messageInvalid, err := utils.ValidateInput(&newRack)
	if err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "invalid input data", messageInvalid)
		return
	}

	if err := h.Service.UpdateRackById(ctx, newRack); err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "rack not found", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "success", nil)
}
