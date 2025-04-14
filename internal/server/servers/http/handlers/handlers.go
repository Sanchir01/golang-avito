package httphandlers

import (
	"net/http"

	_ "github.com/Sanchir01/golang-avito/docs"
	"github.com/Sanchir01/golang-avito/internal/app"
	"github.com/Sanchir01/golang-avito/internal/server/servers/http/custommiddleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func StartHTTTPHandlers(handlers *app.Handlers) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.RequestID, middleware.Recoverer)

	router.Route("/api", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/login", handlers.UserHandler.LoginHandler)
			r.Post("/register", handlers.UserHandler.RegistrationHandler)
			r.Post("/dummyLogin", handlers.UserHandler.DummyLoginHandler)
		})
		r.Group(func(r chi.Router) {
			r.Use(custommiddleware.AuthMiddleware("moderator"))
			r.Post("/pvz", handlers.PVZHandelr.CreatePVZHandler)
		})
		r.Group(func(r chi.Router) {
			r.Use(custommiddleware.AuthMiddleware("employee"))
			r.Post("/receptions", handlers.AcceptanceHandler.CreateAcceptanceHandler)
			r.Post("/pvz/{pvzId}/close_last_reception", handlers.AcceptanceHandler.CloseLastAcceptanceHandler)
			r.Post("/products", handlers.ProductHandler.CreateProductHandler)
			r.Post("/products/{acceptanceID}/delete_last_product", handlers.ProductHandler.DeleteProductHandler)
		})
		r.Group(func(r chi.Router) {
			r.Use(custommiddleware.AuthMiddleware("moderator", "employee"))
			r.Get("/pvz", handlers.PVZHandelr.GetAllPVZHandler)
		})
		r.Get("/pvz_grpc", handlers.PVZHandelr.GetAllGRPCPVZHandler)
	})
	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))
	return router
}

func StartPrometheusHandlers() http.Handler {
	router := chi.NewRouter()
	router.Use(custommiddleware.PrometheusMiddleware)
	router.Handle("/metrics", promhttp.Handler())
	return router
}
