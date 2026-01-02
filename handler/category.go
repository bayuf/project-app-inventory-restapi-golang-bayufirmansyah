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

func (h *CategoryHandler) List(w http.ResponseWriter, r *http.Request) {
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

	categories, pagination, err := h.Service.GetAllCategories(ctx, page, limit)
	if err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "cant get all categories", err)
		return
	}

	utils.ResponsePagination(w, http.StatusOK, "success", categories, *pagination)
}

func (h *CategoryHandler) GetById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method != "GET" {
		utils.ResponseFailed(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}

	strCategory := chi.URLParam(r, "category_id")
	categoryId, err := strconv.Atoi(strCategory)
	if err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "id invalid", nil)
		return
	}

	category, err := h.Service.GetCategoryById(ctx, categoryId)
	if err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "cant find category", errors.New("category not found"))
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "success", category)
}

func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method != "PUT" {
		utils.ResponseFailed(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}

	strCategory := chi.URLParam(r, "category_id")
	categoryId, err := strconv.Atoi(strCategory)
	if err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "id invalid", nil)
		return
	}

	newCategory := dto.CategoryUpdate{
		ID: categoryId,
	}
	if err := json.NewDecoder(r.Body).Decode(&newCategory); err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "invalid input format", err)
		return
	}

	// validate
	messageInvalid, err := utils.ValidateInput(&newCategory)
	if err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "invalid input data", messageInvalid)
		return
	}

	if err := h.Service.UpdateCategory(ctx, newCategory); err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "category not found", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "success", nil)
}

func (h *CategoryHandler) DeleteById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method != "DELETE" {
		utils.ResponseFailed(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}

	strCategory := chi.URLParam(r, "category_id")
	categoryId, err := strconv.Atoi(strCategory)
	if err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "id invalid", nil)
		return
	}

	err = h.Service.DeleteCategory(ctx, categoryId)
	if err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "cant find category", errors.New("category not found"))
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "success", nil)
}
