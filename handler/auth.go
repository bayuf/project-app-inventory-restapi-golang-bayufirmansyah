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
	ctx := r.Context()

	if r.Method != "POST" {
		utils.ResponseFailed(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}

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
	session, err := h.Service.Login(ctx, user)
	if err != nil {
		utils.ResponseFailed(w, http.StatusUnauthorized, "email or password is wrong", err)
		return
	}

	respSession := dto.ResponseSession{
		ID:        session.ID,
		UserID:    session.UserID,
		Username:  session.Username,
		RoleId:    session.RoleId,
		RoleName:  session.RoleName,
		CreatedAt: session.CreatedAt,
		ExpiresAt: session.ExpiresAt,
	}

	utils.ResponseSuccess(w, http.StatusOK, "success", respSession)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method != "POST" {
		utils.ResponseFailed(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}

	sessionId, ok := middleware.GetAuthUser(r)
	if !ok {
		utils.ResponseFailed(w, http.StatusUnauthorized, "user not authenticated", nil)
		return
	}

	if err := h.Service.Logout(ctx, sessionId.ID); err != nil {
		utils.ResponseFailed(w, http.StatusUnauthorized, "cant log out", err)
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "success logout", nil)

}
