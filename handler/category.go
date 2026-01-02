package handler

import (
	"encoding/json"
	"net/http"

	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/dto"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/service"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/utils"
	"go.uber.org/zap"
)

type CategoryHandler struct {
	Service *service.CategoryService
	Logger  *zap.Logger
	Config  *utils.Configuration
}

func NewCategoryHandler(service *service.CategoryService, log *zap.Logger, config *utils.Configuration) *CategoryHandler {
	return &CategoryHandler{
		Service: service,
		Logger:  log,
		Config:  config,
	}
}

func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method != "POST" {
		utils.ResponseFailed(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}

	newCategory := dto.CategoryAdd{}
	if err := json.NewDecoder(r.Body).Decode(&newCategory); err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "invalid input format", err.Error())
		return
	}

	// validate
	messageInvalid, err := utils.ValidateInput(&newCategory)
	if err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "invalid input data", messageInvalid)
		return
	}

	if err := h.Service.CreateNewCategory(ctx, newCategory); err != nil {
		h.Logger.Error("failed to create new category", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "failed to create new category", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusCreated, "success", nil)
}
