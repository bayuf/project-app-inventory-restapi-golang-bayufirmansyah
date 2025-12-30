package middleware

import (
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/service"
	"go.uber.org/zap"
)

type Middleware struct {
	Service *service.Service
	Logger  *zap.Logger
}

func NewCustomMiddleware(service *service.Service, log *zap.Logger) *Middleware {
	return &Middleware{
		Service: service,
		Logger:  log,
	}
}
