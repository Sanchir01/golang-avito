package app

import (
	"github.com/Sanchir01/golang-avito/internal/feature/user"
	"log/slog"
)

type Handlers struct {
	UserHandler *user.Handler
}

func NewHandlers(s *Services, log *slog.Logger) *Handlers {
	return &Handlers{
		UserHandler: user.NewHandler(s.UserService, log),
	}
}
