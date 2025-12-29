package handler

import (
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/service"
	"go.uber.org/zap"
)

type Handler struct {
}

func NewHandler(svc *service.Service, log *zap.Logger) *Handler {
	return &Handler{}
}
