package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/dto"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/service"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ItemHandler struct {
	Service *service.ItemService
	Logger  *zap.Logger
	Config  *utils.Configuration
}

func NewItemHandler(service *service.ItemService, log *zap.Logger, config *utils.Configuration) *ItemHandler {
	return &ItemHandler{
		Service: service,
		Logger:  log,
		Config:  config,
	}
}

func (h *ItemHandler) InputNewItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method != "POST" {
		utils.ResponseFailed(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}

	newItem := dto.ItemAdd{}
	if err := json.NewDecoder(r.Body).Decode(&newItem); err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "invalid input format", err)
		return
	}

	// validate
	messageInvalid, err := utils.ValidateInput(&newItem)
	if err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "invalid input data", messageInvalid)
		return
	}

	if err := h.Service.InputNewItem(ctx, newItem); err != nil {
		h.Logger.Error("failed to create new item", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "failed to create new item", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusCreated, "success", nil)
}

func (h *ItemHandler) GetAllItems(w http.ResponseWriter, r *http.Request) {
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

	items, pagination, err := h.Service.GetAllItems(ctx, page, limit)
	if err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "cant get all items", err)
		return
	}

	utils.ResponsePagination(w, http.StatusOK, "success", items, *pagination)
}

func (h *ItemHandler) GetItemById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method != "GET" {
		utils.ResponseFailed(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}

	strItemId := chi.URLParam(r, "item_id")

	uuidItemId, err := uuid.Parse(strItemId)
	if err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "id not valid", err.Error())
	}

	items, err := h.Service.GetItem(ctx, uuidItemId)
	if err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "cant get item", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "success", items)
}

func (h *ItemHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method != "DELETE" {
		utils.ResponseFailed(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}

	strItemId := chi.URLParam(r, "item_id")

	uuidItemId, err := uuid.Parse(strItemId)
	if err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "id not valid", err.Error())
	}

	if err := h.Service.DeleteItem(ctx, uuidItemId); err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "cant delete item", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "success", nil)

}
