package router

import (
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/handler"
	middlewareCustom "github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/middleware"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func NewRouter(handler *handler.Handler, service *service.Service, log *zap.Logger) *chi.Mux {
	r := chi.NewRouter()
	mw := middlewareCustom.NewCustomMiddleware(service, log)

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Mount("/api", Api(handler, mw))

	return r
}

func Api(handler *handler.Handler, mw *middlewareCustom.Middleware) *chi.Mux {
	r := chi.NewRouter()

	r.Post("/create_user", handler.UserHandler.Create)

	return r
}
