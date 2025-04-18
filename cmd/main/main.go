package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Sanchir01/golang-avito/internal/app"
	grpcserver "github.com/Sanchir01/golang-avito/internal/server/servers/grpc"
	httpserver "github.com/Sanchir01/golang-avito/internal/server/servers/http"
	httphandlers "github.com/Sanchir01/golang-avito/internal/server/servers/http/handlers"
	"github.com/fatih/color"
)

// @title 🚀 Avito testovoe
// @version 1.0
// @description This is a sample server celler
// @termsOfService http://swagger.io/terms/

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @contact.name GitHub
// @contact.url https://github.com/Sanchir01
func main() {
	env, err := app.NewEnv()
	if err != nil {
		panic(err)
	}
	serve := httpserver.NewHTTPServer(env.Cfg.Servers.HTTPServer.Host, env.Cfg.Servers.HTTPServer.Port,
		env.Cfg.Servers.HTTPServer.Timeout, env.Cfg.Servers.HTTPServer.IdleTimeout)

	prometheusserver := httpserver.NewHTTPServer(
		env.Cfg.Servers.PrometheusServer.Host, env.Cfg.Servers.PrometheusServer.Port,
		env.Cfg.Servers.PrometheusServer.Timeout, env.Cfg.Servers.PrometheusServer.IdleTimeout)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	PVZGRPCClient, err := grpcserver.NewGRPCClient(ctx, env.Lg, env.Cfg.Servers.GRPCServer.GRPCPVZ.Port, env.Cfg.Servers.GRPCServer.GRPCPVZ.Host,
		env.Cfg.Servers.GRPCServer.GRPCPVZ.Retries)
	allpvzgrpc, err := PVZGRPCClient.AllPVZHandler(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(allpvzgrpc)
	green := color.New(color.FgGreen).SprintFunc()
	env.Lg.Info(green("🚀 Server started successfully!"),
		slog.String("time", time.Now().Format("2006-01-02 15:04:05")),
		slog.String("port", env.Cfg.Servers.HTTPServer.Port),
	)
	env.Lg.Info(green("🚀  Prometheus Server started successfully!"),
		slog.String("time", time.Now().Format("2006-01-02 15:04:05")),
		slog.String("port", env.Cfg.Servers.PrometheusServer.Port),
	)
	defer cancel()

	go func() {
		if err := serve.Run(httphandlers.StartHTTTPHandlers(env.Handlers)); err != nil {
			if !errors.Is(err, context.Canceled) {
				env.Lg.Error("Listen server error", slog.String("error", err.Error()))
				return
			}
		}
	}()

	go func() {
		if err := prometheusserver.Run(httphandlers.StartPrometheusHandlers()); err != nil {
			if !errors.Is(err, context.Canceled) {
				env.Lg.Error("Listen server error", slog.String("error", err.Error()))
				return
			}
		}
	}()
	<-ctx.Done()
	if err := serve.Gracefull(ctx); err != nil {
		env.Lg.Error("server gracefull")
	}
	if err := env.Database.Close(); err != nil {
		env.Lg.Error("Close database", slog.String("error", err.Error()))
	}
}
