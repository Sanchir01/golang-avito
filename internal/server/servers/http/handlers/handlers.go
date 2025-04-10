package httphandlers

import (
	"net/http"

	"github.com/Sanchir01/golang-avito/internal/app"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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
	})

	return router
}
