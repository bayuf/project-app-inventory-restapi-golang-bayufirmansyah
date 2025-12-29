package router

import (
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/handler"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/service"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func NewRouter(handler *handler.Handler, service *service.Service, log *zap.Logger) *chi.Mux {
	r := chi.NewRouter()

	return r
}
