package handler

import (
	"encoding/json"
	"net/http"

	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/dto"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/service"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/utils"
	"go.uber.org/zap"
)

type AuthHandler struct {
	Service *service.AuthService
	Logger  *zap.Logger
}

func NewAuthHandler(service *service.AuthService, log *zap.Logger) *AuthHandler {
	return &AuthHandler{
		Service: service,
		Logger:  log,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	user := dto.UserReq{}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "invalid input format", err)
		return
	}

	// validation
	messageInvalid, err := utils.ValidateInput(&user)
	if err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "invalid input data", messageInvalid)
		return
	}

	// get user id and role(check email and password)
	session, err := h.Service.Login(user)
	if err != nil {
		utils.ResponseFailed(w, http.StatusUnauthorized, "email or password is wrong", err)
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "success", session)
}
