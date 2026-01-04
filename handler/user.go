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

type UserHandler struct {
	Service *service.UserService
	Logger  *zap.Logger
	Config  *utils.Configuration
}

func NewUserHandler(service *service.UserService, log *zap.Logger, config *utils.Configuration) *UserHandler {
	return &UserHandler{
		Service: service,
		Logger:  log,
		Config:  config,
	}
}

func (h *UserHandler) ShowMyData(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId, ok := middleware.GetAuthUser(r)
	if !ok {
		utils.ResponseFailed(w, http.StatusUnauthorized, "user not authenticated", nil)
		return
	}

	user, err := h.Service.GetUserData(ctx, userId.UserID)
	if err != nil {
		utils.ResponseFailed(w, http.StatusUnauthorized, "invlaid credentials", err)
		return
	}

	userRes := dto.UserResponse{
		ID:         userId.ID.String(),
		Name:       user.Name,
		RoleName:   user.RoleName,
		Email:      user.Email,
		Created_At: user.Created_At,
	}

	utils.ResponseSuccess(w, http.StatusOK, "success", userRes)
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method != "POST" {
		utils.ResponseFailed(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}

	newUser := dto.UserAdd{}
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "invalid input format", err)
		return
	}

	// validate
	messageInvalid, err := utils.ValidateInput(&newUser)
	if err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "invalid input data", messageInvalid)
		return
	}

	if err := h.Service.AddUser(ctx, newUser); err != nil {
		h.Logger.Error("failed to create new user", zap.Error(err))
		utils.ResponseFailed(w, http.StatusInternalServerError, "failed to create new user", err)
		return
	}

	utils.ResponseSuccess(w, http.StatusCreated, "success", nil)
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method != "GET" {
		utils.ResponseFailed(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}

	page, err := strconv.Atoi(r.URL.Query().Get("page"))

	if err != nil {
		h.Logger.Info("invalid page :", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "invalid page", err.Error())
		return
	}

	// page limit
	limit := h.Config.PageLimit

	items, pagination, err := h.Service.GetAllUsersData(ctx, page, limit)
	if err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "cant get all users", err)
		return
	}

	utils.ResponsePagination(w, http.StatusOK, "success", items, *pagination)
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method != "POST" {
		utils.ResponseFailed(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}

	newUser := dto.UserRegister{}
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "invalid input format", err)
		return
	}

	// validate
	messageInvalid, err := utils.ValidateInput(&newUser)
	if err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "invalid input data", messageInvalid)
		return
	}

	if err := h.Service.RegisterUser(ctx, newUser); err != nil {
		h.Logger.Error("failed to create new user", zap.Error(err))
		utils.ResponseFailed(w, http.StatusInternalServerError, "register failed", err)
		return
	}

	utils.ResponseSuccess(w, http.StatusCreated, "success", nil)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	strUserId := chi.URLParam(r, "user_id")
	userId, err := uuid.Parse(strUserId)
	if err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "invalid user id", err.Error())
		return
	}

	if err := h.Service.DeleteUser(ctx, userId); err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "error", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "success", nil)
}

func (h *UserHandler) SuspendUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	strUserId := chi.URLParam(r, "user_id")
	userId, err := uuid.Parse(strUserId)
	if err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "invalid user id", err.Error())
		return
	}

	suspend := dto.UserSuspend{ID: userId}
	if err := json.NewDecoder(r.Body).Decode(&suspend); err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "invalid input format", err)
		return
	}

	// validate
	messageInvalid, err := utils.ValidateInput(&suspend)
	if err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "invalid input data", messageInvalid)
		return
	}

	if err := h.Service.SuspendUser(ctx, suspend); err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "error", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "success", nil)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	strUserId := chi.URLParam(r, "user_id")
	userId, err := uuid.Parse(strUserId)
	if err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "invalid user id", err.Error())
		return
	}

	newData := dto.UserUpdate{ID: userId}
	if err := json.NewDecoder(r.Body).Decode(&newData); err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "invalid input format", err)
		return
	}

	// validate
	messageInvalid, err := utils.ValidateInput(&newData)
	if err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "invalid input data", messageInvalid)
		return
	}

	if err := h.Service.UpdateUser(ctx, newData); err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "error", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "success", nil)
}

func (h *UserHandler) UpdateMyUserData(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, ok := middleware.GetAuthUser(r)
	if !ok {
		utils.ResponseFailed(w, http.StatusUnauthorized, "user not authenticated", nil)
		return
	}

	newData := dto.UserSelfUpdate{ID: user.UserID}
	if err := json.NewDecoder(r.Body).Decode(&newData); err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "invalid input format", err)
		return
	}

	// validate
	messageInvalid, err := utils.ValidateInput(&newData)
	if err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "invalid input data", messageInvalid)
		return
	}

	if err := h.Service.UpdateMyUserData(ctx, newData); err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "error", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "success", nil)
}
