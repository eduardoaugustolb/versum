package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	redisCache "github.com/eduardoaugustolb/versum/api/internal/cache/redis"
	catalogpostgres "github.com/eduardoaugustolb/versum/api/internal/catalog/adapters/postgres"
	"github.com/eduardoaugustolb/versum/api/internal/catalog/application/queries"
	"github.com/eduardoaugustolb/versum/api/internal/clock"
	"github.com/eduardoaugustolb/versum/api/internal/config"
	keyring2 "github.com/eduardoaugustolb/versum/api/internal/cryptography/keyring"
	"github.com/eduardoaugustolb/versum/api/internal/database/postgres"
	"github.com/eduardoaugustolb/versum/api/internal/health"
	"github.com/eduardoaugustolb/versum/api/internal/id"
	identitycryptography "github.com/eduardoaugustolb/versum/api/internal/identityaccess/adapters/cryptography"
	identityaccessPg "github.com/eduardoaugustolb/versum/api/internal/identityaccess/adapters/postgres"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/application/commands"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/application/policy"
	outboxcryptography "github.com/eduardoaugustolb/versum/api/internal/outboxevent/adapters/cryptography"
	"github.com/eduardoaugustolb/versum/api/internal/transport/httpapi"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Error("error loading .env file")
		os.Exit(1)
	}
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

	dbExecutor := postgres.NewPgxExecutor(pool)

	rdb, err := newRedisClient(cfg.RedisURL)
	if err != nil {
		slog.Error("failed to create redis client", "error", err)
		os.Exit(1)
	}

	redisCacheClient := redisCache.NewRedisCache(rdb)

	catalogRepo := catalogpostgres.NewRepository(dbExecutor)

	keyring := keyring2.NewKeyRingStatic(cfg)
	emailProtector := identitycryptography.NewAESGCMEmailProtector(keyring)
	payloadProtector := outboxcryptography.NewAESGCMPayloadProtector(keyring)

	randomTokenGenerator := identitycryptography.RandomTokenGenerator{}

	defaultHasher := identitycryptography.SHA256Hasher{}
	uuidGenerator := id.UUIDGenerator{}

	defaultClock := clock.SystemClock{}

	magicLinkPolicy := policy.MagicLinkPolicy{}
	magicLinkPolicy.ApplyDefaults()

	indentityAccessUnitOfWork := identityaccessPg.NewUnitOfWork(dbExecutor, emailProtector, payloadProtector)

	requestMagicLink, err := commands.NewRequestMagicLink(
		indentityAccessUnitOfWork,
		randomTokenGenerator,
		defaultHasher,
		uuidGenerator,
		defaultClock,
		magicLinkPolicy.TTL,
	)

	if err != nil {
		slog.Error("failed to create request magic link", "error", err)
		os.Exit(1)
	}

	router := httpapi.NewHandler(httpapi.Dependencies{
		Health: health.NewCheckHealth(),
		Catalog: httpapi.CatalogDependencies{
			ListBooks:  queries.NewListBooks(catalogRepo),
			GetChapter: queries.NewGetChapter(catalogRepo),
		},
		IdentityAccess: httpapi.IdentityAccessDependencies{
			RequestMagicLink: requestMagicLink,
			Cache:            redisCacheClient,
		},
		Cache: redisCacheClient,
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

func newRedisClient(redisURL string) (*redis.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rdbOpts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse redis url: %w", err)
	}
	rdb := redis.NewClient(rdbOpts)
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}
	return rdb, nil
}

// newDatabasePool opens the pool and confirms the database is reachable
// before the server accepts traffic, instead of failing on the first request.
func newDatabasePool(databaseURL string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create database pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
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
