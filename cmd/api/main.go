package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"security/config"
	"security/db"
	"security/handler"
	"security/repository"
	"security/router"
	"security/service"

	_ "security/docs"
)

//	@title			ODIN Security API
//	@version		1.0
//	@description	Security pipeline API for ODIN agent versions (normalize, quarantine, analyze, smoke test, sandbox).
//	@BasePath		/

//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Bearer token. Example: "Bearer {token}"

//	@securityDefinitions.apikey	AdminAuth
//	@in							header
//	@name						X-Admin
//	@description				Admin flag placeholder. Use "true" for admin access.

func main() {
	if err := run(); err != nil {
		slog.Error("API stopped", "error", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := db.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("connect to database: %w", err)
	}
	defer pool.Close()
	slog.Info("connected to database")

	versionRepo := repository.NewVersionRepository()
	securityRepo := repository.NewSecurityRepository()
	pipelineRepo := repository.NewPipelineRepository()

	versionSvc := service.NewVersionService(versionRepo)
	securitySvc := service.NewSecurityService(securityRepo)
	pipelineSvc := service.NewPipelineService(pipelineRepo)

	r := router.New(router.Deps{
		Version:  handler.NewVersionHandler(versionSvc),
		Security: handler.NewSecurityHandler(securitySvc),
		Pipeline: handler.NewPipelineHandler(pipelineSvc),
	})

	server := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	listener, err := net.Listen("tcp", cfg.HTTPAddr)
	if err != nil {
		return err
	}
	serveErr := make(chan error, 1)
	go func() { serveErr <- server.Serve(listener) }()
	slog.Info("API listening", "address", listener.Addr().String())
	select {
	case err := <-serveErr:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-ctx.Done():
		stop()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			_ = server.Close()
			return err
		}
		slog.Info("API stopped")
		return nil
	}
}
