package handler

import (
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/service"
	"go.uber.org/zap"
)

type Handler struct {
	UserHandler *UserHandler
	AuthHandler *AuthHandler
}

func NewHandler(svc *service.Service, log *zap.Logger) *Handler {
	return &Handler{
		UserHandler: NewUserHandler(svc.UserService, log),
		AuthHandler: NewAuthHandler(svc.AuthService, log),
	}
}
