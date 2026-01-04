package router

import (
	"net/http"

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

	// checking
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// main route
	r.Mount("/api/v1", Apiv1(handler, service, mw))

	return r
}

func Apiv1(handler *handler.Handler, service *service.Service, mw *middlewareCustom.Middleware) *chi.Mux {
	r := chi.NewRouter()

	// superAdmin := mw.AuthMiddleware.RequireRoles("super_admin")
	adminOnly := mw.AuthMiddleware.RequireRoles("super_admin", "admin")
	allRoles := mw.AuthMiddleware.RequireRoles("super_admin", "admin", "staff")

	// auth routes
	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", handler.AuthHandler.Login)

		// protected
		r.Group(func(r chi.Router) {
			r.Use(mw.AuthMiddleware.SessionAuthMiddleware())
			r.Post("/logout", handler.AuthHandler.Logout)
		})
	})

	// public path USERS
	r.Post("/users/register", handler.UserHandler.Register)

	// protected routes
	r.Group(func(r chi.Router) {
		r.Use(mw.AuthMiddleware.SessionAuthMiddleware())

		//USERS
		r.Route("/users", func(r chi.Router) {
			// READ (staff + admin + super admin)
			r.Get("/me", handler.UserHandler.ShowMyData)

			// READ (admin + super admin)
			r.With(adminOnly).Get("/", handler.UserHandler.GetAllUsers)

			// WRITE (staff + admin + super admin)
			r.With(allRoles).Put("/update", handler.UserHandler.UpdateMyUserData)

			// WRITE (admin + super admin)
			r.With(adminOnly).Post("/create", handler.UserHandler.Create)
			r.With(adminOnly).Put("/{user_id}", handler.UserHandler.UpdateUser)
			r.With(adminOnly).Delete("/{user_id}", handler.UserHandler.DeleteUser)
			r.With(adminOnly).Patch("/{user_id}", handler.UserHandler.SuspendUser)

		})

		// WAREHOUSES
		r.Route("/warehouses", func(r chi.Router) {
			// READ (staff + admin + super_admin)
			r.With(allRoles).Get("/", handler.WarehouseHandler.List)
			r.With(allRoles).Get("/{warehouse_id}", handler.WarehouseHandler.GetById)

			// WRITE (admin + super_admin)
			r.With(adminOnly).Post("/", handler.WarehouseHandler.CreateWarehouse)
			r.With(adminOnly).Put("/{warehouse_id}", handler.WarehouseHandler.Update)
			r.With(adminOnly).Delete("/{warehouse_id}", handler.WarehouseHandler.Delete)
		})

		// RACKS
		r.Route("/racks", func(r chi.Router) {
			// READ (staff + admin + super_admin)
			r.With(allRoles).Get("/", handler.RackHandler.GetAllRacks)
			r.With(allRoles).Get("/{rack_id}", handler.RackHandler.GetRackById)

			// WRITE (admin + super_admin)
			r.With(adminOnly).Post("/", handler.RackHandler.Create)
			r.With(adminOnly).Put("/{rack_id}", handler.RackHandler.UpdateRack)
			r.With(adminOnly).Delete("/{rack_id}", handler.RackHandler.DeleteRack)

		})

		// CATEGORIES
		r.Route("/categories", func(r chi.Router) {
			// READ (staff + admin + super admin)
			r.With(allRoles).Get("/", handler.CategoryHandler.List)
			r.With(allRoles).Get("/{category_id}", handler.CategoryHandler.GetById)

			// WRITE (admin + super admin)
			r.With(adminOnly).Post("/", handler.CategoryHandler.CreateCategory)
			r.With(adminOnly).Put("/{category_id}", handler.CategoryHandler.Update)
			r.With(adminOnly).Delete("/{category_id}", handler.CategoryHandler.DeleteById)
		})

		// ITEMS
		r.Route("/items", func(r chi.Router) {
			// READ (staff + admin + super admin)
			r.With(allRoles).Get("/", handler.ItemHandler.GetAllItems)
			r.With(allRoles).Get("/{item_id}", handler.ItemHandler.GetItemById)

			// WRITE (admin + super admin)
			r.With(adminOnly).Post("/", handler.ItemHandler.InputNewItem)
			r.With(adminOnly).Put("/{item_id}", handler.ItemHandler.UpdateItem)
			r.With(adminOnly).Delete("/{item_id}", handler.ItemHandler.DeleteItem)

			// WRITE (staff + admin + super admin)
			r.With(allRoles).Patch("/stock-ajustment", handler.ItemHandler.StockAdjustment)

		})

		// SALES
		r.Route("/sales", func(r chi.Router) {
			// READ (staff + admin + super admin)
			r.With(allRoles).Get("/", handler.SaleHandler.GetAllSales) //staff only getAll by userId
			r.With(allRoles).Get("/{sale_id}", handler.SaleHandler.GetSaleInfoStaff)

			// READ (admin + super admin)
			r.With(adminOnly).Get("/{sale_id}", handler.SaleHandler.GetSaleInfo)

			// WRITE (admin + super admin)
			r.With(adminOnly).Delete("/{sale_id}", handler.SaleHandler.DeleteSale)

			// WRITE (staff + admin + super admin)
			r.With(allRoles).Post("/", handler.SaleHandler.InsertNewSale)
			r.With(allRoles).Patch("/{sale_id}/status", handler.SaleHandler.UpdateSaleStatus)

		})

		// REPORT
		r.Route("/reports", func(r chi.Router) {
			// READ (staff + admin + super admin)
			r.With(allRoles).Group(func(r chi.Router) {
				r.Get("/inventory", handler.ReportHandler.GetItemsReport)
				r.Get("/sales", handler.ReportHandler.GetSalesReport)
			})

			// READ (admin + super admin)
			r.With(adminOnly).Get("/revenue", handler.ReportHandler.GetRevenueReport)
		})
	})

	return r
}
