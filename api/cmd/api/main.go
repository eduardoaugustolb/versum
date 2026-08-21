package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/eduardoaugustolb/versum/api/internal/adapters/postgres"
	"github.com/eduardoaugustolb/versum/api/internal/catalog/application/queries"
	catalogpostgres "github.com/eduardoaugustolb/versum/api/internal/catalog/postgres"
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

	pool, err := newDatabasePool(cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	repo := catalogpostgres.NewRepository(postgres.NewPgxExecutor(pool))

	router := httpapi.NewRouter(httpapi.Dependencies{
		Health: health.NewCheckHealth(),
		Catalog: httpapi.CatalogDependencies{
			ListBooks:  queries.NewListBooks(repo),
			GetChapter: queries.NewGetChapter(repo),
		},
	})

	srv := &http.Server{
		Addr:              cfg.Address,
		Handler:           router,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}

	runServer(srv, cfg)
}

// newDatabasePool opens the pool and confirms the database is reachable
// before the server accepts traffic, instead of failing on the first request.
func newDatabasePool(databaseURL string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}

// runServer starts srv and blocks until it exits, either because
// ListenAndServe failed or because the process received a shutdown signal.
func runServer(srv *http.Server, cfg config.Config) {
	serverErr := make(chan error, 1)
	go func() {
		slog.Info("starting server", "address", cfg.Address, "environment", cfg.Environment)
		serverErr <- srv.ListenAndServe()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server stopped unexpectedly", "error", err)
			os.Exit(1)
		}
	case <-stop:
		slog.Info("shutting down")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			slog.Error("graceful shutdown failed", "error", err)
			os.Exit(1)
		}
	}
}
