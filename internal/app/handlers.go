package app

import (
	"log/slog"

	"github.com/Sanchir01/golang-avito/internal/feature/acceptance"
	"github.com/Sanchir01/golang-avito/internal/feature/product"
	grpcserver "github.com/Sanchir01/golang-avito/internal/server/servers/grpc"

	"github.com/Sanchir01/golang-avito/internal/feature/pvz"
	"github.com/Sanchir01/golang-avito/internal/feature/user"
)

type Handlers struct {
	UserHandler       *user.Handler
	PVZHandelr        *pvz.Handler
	AcceptanceHandler *acceptance.Handler
	ProductHandler    *product.Handler
}

func NewHandlers(s *Services, log *slog.Logger, pvzgrpc *grpcserver.GRPCClientPVZ) *Handlers {
	return &Handlers{
		UserHandler:       user.NewHandler(s.UserService, log),
		PVZHandelr:        pvz.NewHandler(s.PVZService, log, pvzgrpc),
		AcceptanceHandler: acceptance.NewHandler(s.AcceptanceService, log),
		ProductHandler:    product.NewHandler(s.ProductService, log),
	}
}
