package app

import (
	"log/slog"
)

type Handlers struct {
}

func NewHandlers(services *Services, log *slog.Logger) *Handlers {
	return &Handlers{}
}
