package handler

import (
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/service"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/utils"
	"go.uber.org/zap"
)

type Handler struct {
	*UserHandler
	*AuthHandler

	*WarehouseHandler
	*RackHandler
	*CategoryHandler
}

func NewHandler(svc *service.Service, log *zap.Logger, config *utils.Configuration) *Handler {
	return &Handler{
		UserHandler: NewUserHandler(svc.UserService, log, config),
		AuthHandler: NewAuthHandler(svc.AuthService, log),

		WarehouseHandler: NewWarehouseHandler(svc.WarehousesService, log, config),
		RackHandler:      NewRackHandler(svc.RackService, log, config),
		CategoryHandler:  NewCategoryHandler(svc.CategoryService, log, config),
	}
}
