package httphandlers

import (
	"net/http"

	"github.com/Sanchir01/golang-avito/internal/app"
	"github.com/Sanchir01/golang-avito/internal/server/servers/http/custommiddleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func StartHTTTPHandlers(handlers *app.Handlers) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.RequestID, middleware.Recoverer)

	router.Route("/api", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/login", handlers.UserHandler.LoginHandler)
			r.Post("/register", handlers.UserHandler.RegistrationHandler)
			r.Post("/dummyLogin", handlers.UserHandler.DammyLogin)
		})
		r.Group(func(r chi.Router) {
			r.Use(custommiddleware.AuthMiddleware("admin"))
			r.Post("/pvz", handlers.PVZHandelr.Create)
		})
	})

	return router
}

func StartPrometheusHandlers() http.Handler {
	router := chi.NewRouter()
	router.Use(custommiddleware.PrometheusMiddleware)
	router.Handle("/metrics", promhttp.Handler())
	return router
}
