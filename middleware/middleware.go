package middleware

import (
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/service"
	"go.uber.org/zap"
)

type Middleware struct {
	AuthMiddleware *AuthMiddleware
}

func NewCustomMiddleware(service *service.Service, log *zap.Logger) *Middleware {
	return &Middleware{
		AuthMiddleware: NewAuthMiddleware(service.AuthService, log),
	}
}
