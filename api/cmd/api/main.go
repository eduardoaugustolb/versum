package main

import (
	"errors"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/eduardoaugustolb/versum/api/internal/config"
	"github.com/eduardoaugustolb/versum/api/internal/health"
	"github.com/eduardoaugustolb/versum/api/internal/transport/httpapi"
)

func main() {
	cfg, err := config.Load(os.Getenv)
	if err != nil {
		slog.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	router := httpapi.NewRouter(health.NewCheckHealth())

	srv := http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	slog.Info("starting server", "address", cfg.Address, "environment", cfg.Environment)

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("server stopped unexpectedly", "error", err)
		os.Exit(1)
	}
}
