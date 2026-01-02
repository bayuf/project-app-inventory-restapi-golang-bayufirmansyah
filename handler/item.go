package handler

import (
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/service"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/utils"
	"go.uber.org/zap"
)

type ItemHandler struct {
	Service *service.ItemService
	Logger  *zap.Logger
	Config  *utils.Configuration
}

func NewItemHandler(service *service.ItemService, log *zap.Logger, config *utils.Configuration) *ItemHandler {
	return &ItemHandler{
		Service: service,
		Logger:  log,
		Config:  config,
	}
}
