package app

import (
	"context"
	"log/slog"

	"github.com/Sanchir01/golang-avito/internal/config"
	grpcserver "github.com/Sanchir01/golang-avito/internal/server/servers/grpc"
)

type Env struct {
	Lg       *slog.Logger
	Cfg      *config.Config
	Database *Database
	Handlers *Handlers
}

func NewEnv() (*Env, error) {
	cfg := config.MustLoadConfig()
	logger := setupLogger(cfg.Env)
	ctx := context.Background()
	primarydb, err := NewDataBases(ctx, cfg.PrimaryDB.User, cfg.PrimaryDB.Host, cfg.PrimaryDB.Dbname, cfg.PrimaryDB.Port, cfg.PrimaryDB.MaxAttempts)
	if err != nil {
		return nil, err
	}
	pvzgrpc, err := grpcserver.NewGRPCClient(
		ctx, logger, cfg.Servers.GRPCServer.GRPCPVZ.Port,
		cfg.Servers.GRPCServer.GRPCPVZ.Host,
		cfg.Servers.GRPCServer.GRPCPVZ.Retries,
	)
	if err != nil {
		return nil, err
	}
	repositories := NewRepositories(primarydb)
	services := NewServices(repositories, primarydb)
	handlers := NewHandlers(services, logger, pvzgrpc)
	env := Env{
		Lg:       logger,
		Cfg:      cfg,
		Database: primarydb,
		Handlers: handlers,
	}

	return &env, nil
}
