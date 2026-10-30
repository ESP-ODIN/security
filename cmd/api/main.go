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

	"github.com/ESP-ODIN/authkit-go"

	"security/config"
	"security/db"
	"security/internal/handler"
	repositoryimpl "security/internal/repository/implementation"
	"security/internal/router"
	serviceimpl "security/internal/service/implementation"
)

//	@title			ODIN Security API
//	@version		1.0
//	@description	Security pipeline API for ODIN agent versions (normalize, quarantine, analyze, smoke test, sandbox).
//	@BasePath		/

//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				JWT issued by Auth, sent as "Bearer <token>".

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
	authenticator, err := authkit.New(authkit.Config{
		JWKSURL:  cfg.AuthJWKSURL,
		Issuer:   cfg.AuthIssuer,
		Audience: cfg.AuthAudience,
	})
	if err != nil {
		return fmt.Errorf("initialize authentication: %w", err)
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := db.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("connect to database: %w", err)
	}
	defer pool.Close()
	slog.Info("connected to database")

	versionRepo := repositoryimpl.NewVersionRepository()
	securityRepo := repositoryimpl.NewSecurityRepository()
	pipelineRepo := repositoryimpl.NewPipelineRepository()

	r := router.New(pool, router.Deps{
		Version:       handler.NewVersionHandler(serviceimpl.NewVersionService(versionRepo)),
		Security:      handler.NewSecurityHandler(serviceimpl.NewSecurityService(securityRepo)),
		Pipeline:      handler.NewPipelineHandler(serviceimpl.NewPipelineService(pipelineRepo), cfg.AdminUserIDs),
		Authenticator: authenticator,
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
