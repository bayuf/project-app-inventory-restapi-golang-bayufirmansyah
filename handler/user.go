package handler

import (
	"encoding/json"
	"net/http"

	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/dto"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/service"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/utils"
	"go.uber.org/zap"
)

type UserHandler struct {
	Service *service.UserService
	Logger  *zap.Logger
}

func NewUserHandler(service *service.UserService, log *zap.Logger) *UserHandler {
	return &UserHandler{
		Service: service,
		Logger:  log,
	}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
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
