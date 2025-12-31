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

	r.Mount("/api/v1", Apiv1(handler, service, mw))

	return r
}

func Apiv1(handler *handler.Handler, service *service.Service, mw *middlewareCustom.Middleware) *chi.Mux {
	r := chi.NewRouter()

	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", handler.AuthHandler.Login)
		// r.Post("/logout", handler.AuthHandler.Login)
	})

	r.Group(func(r chi.Router) {
		r.Use(mw.AuthMiddleware.SessionAuthMiddleware())
		r.Post("/create_user", handler.UserHandler.Create)
	})

	return r
}
