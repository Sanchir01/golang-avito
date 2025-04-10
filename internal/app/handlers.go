package app

import (
	"log/slog"

	"github.com/Sanchir01/golang-avito/internal/feature/user"
)

type Handlers struct {
	UserHandler *user.Handler
}

func NewHandlers(s *Services, log *slog.Logger) *Handlers {
	return &Handlers{
		UserHandler: user.NewHandler(s.UserService, log),
	}
}
