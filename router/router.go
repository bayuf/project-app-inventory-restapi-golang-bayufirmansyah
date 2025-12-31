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

	// auth routes
	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", handler.AuthHandler.Login)
		// protected
		r.Group(func(r chi.Router) {
			r.Use(mw.AuthMiddleware.SessionAuthMiddleware())
			r.Get("/me", handler.UserHandler.ShowMyData)
			r.Post("/logout", handler.AuthHandler.Logout)
		})
	})

	// protected routes
	r.Group(func(r chi.Router) {
		r.Use(mw.AuthMiddleware.SessionAuthMiddleware())

		// users routes
		r.Route("/users", func(r chi.Router) {
			r.With(mw.AuthMiddleware.RequireRoles("super_admin", "admin")).
				Post("/", handler.UserHandler.Create)
		})

		r.Route("/warehouses", func(r chi.Router) {
			r.With(mw.AuthMiddleware.RequireRoles("super_admin", "admin")).Group(func(r chi.Router) {
				r.Post("/", handler.WarehouseHandler.CreateWarehouse)

				r.With(mw.AuthMiddleware.RequireRoles("super_admin", "admin", "staff")).Group(func(r chi.Router) {
					r.Get("/", handler.WarehouseHandler.List)
					r.Route("/{warehouse_id}", func(r chi.Router) { //need id
						r.Get("/", handler.WarehouseHandler.GetById)

					})

				})

			})
		})
	})

	return r
}
