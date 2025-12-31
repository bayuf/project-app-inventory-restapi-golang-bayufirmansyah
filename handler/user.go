package handler

import (
	"encoding/json"
	"net/http"

	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/dto"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/middleware"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/service"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/utils"
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

	userId, ok := middleware.GetAuthUser(r)
	if !ok {
		utils.ResponseFailed(w, http.StatusUnauthorized, "user not authenticated", nil)
		return
	}

	user, err := h.Service.GetUserData(userId.UserID)
	if err != nil {
		utils.ResponseFailed(w, http.StatusUnauthorized, "invlaid credentials", err)
		return
	}

	userRes := dto.UserResponse{
		ID:       userId.ID.String(),
		Name:     user.Name,
		RoleName: user.RoleName,
	}

	utils.ResponseSuccess(w, http.StatusOK, "success", userRes)

}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
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

	if err := h.Service.AddUser(newUser); err != nil {
		h.Logger.Error("failed to create new user", zap.Error(err))
		utils.ResponseFailed(w, http.StatusInternalServerError, "failed to create new user", err)
		return
	}

	utils.ResponseSuccess(w, http.StatusCreated, "success", nil)
}
