package main

import (
	"context"
	"errors"
	"log/slog"
	"net/mail"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/eduardoaugustolb/versum/api/internal/clock"
	"github.com/eduardoaugustolb/versum/api/internal/config"
	"github.com/eduardoaugustolb/versum/api/internal/cryptography/keyring"
	databasepg "github.com/eduardoaugustolb/versum/api/internal/database/postgres"
	resendemail "github.com/eduardoaugustolb/versum/api/internal/email/resend"
	"github.com/eduardoaugustolb/versum/api/internal/id"
	identitycrypto "github.com/eduardoaugustolb/versum/api/internal/identityaccess/adapters/cryptography"
	identitypg "github.com/eduardoaugustolb/versum/api/internal/identityaccess/adapters/postgres"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/application/events"
	outboxcrypto "github.com/eduardoaugustolb/versum/api/internal/outboxevent/adapters/cryptography"
	"github.com/eduardoaugustolb/versum/api/internal/outboxevent/adapters/polling"
	outboxpg "github.com/eduardoaugustolb/versum/api/internal/outboxevent/adapters/postgres"
	"github.com/eduardoaugustolb/versum/api/internal/outboxevent/application"
	"github.com/eduardoaugustolb/versum/api/internal/outboxevent/domain"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := run(ctx); err != nil && ctx.Err() == nil {
		// Do not log connection URLs, provider errors, credentials or payloads.
		slog.Error("outbox worker stopped", "code", "worker_startup_failed")
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return errors.New("cannot load environment file")
	}
	cfg, err := config.Load(os.Getenv)
	if err != nil {
		return err
	}
	apiKey := strings.TrimSpace(os.Getenv("RESEND_API_KEY"))
	from := strings.TrimSpace(os.Getenv("EMAIL_FROM"))
	accessURL := strings.TrimSpace(os.Getenv("MAGIC_LINK_ACCESS_URL"))
	if apiKey == "" || from == "" || accessURL == "" {
		return errors.New("missing worker email configuration")
	}
	if _, err := mail.ParseAddress(from); err != nil {
		return errors.New("invalid worker sender address")
	}

	startupCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	poolConfig, err := pgxpool.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		return err
	}
	poolConfig.ConnConfig.ConnectTimeout = 5 * time.Second
	pool, err := pgxpool.NewWithConfig(startupCtx, poolConfig)
	if err != nil {
		return err
	}
	defer pool.Close()
	if err := pool.Ping(startupCtx); err != nil {
		return err
	}
	if ctx.Err() != nil {
		return nil
	}

	db := databasepg.NewPgxExecutor(pool)
	keys := keyring.NewKeyRingStatic(cfg)
	payloadProtector := outboxcrypto.NewAESGCMPayloadProtector(keys)
	emailProtector := identitycrypto.NewAESGCMEmailProtector(keys)
	identityDAO := identitypg.NewDAO(db)
	defaultClock := clock.SystemClock{}
	sender := resendemail.NewResendSender(apiKey, from)
	handler, err := events.NewMagicLinkRequestedHandler(
		sender,
		identitypg.NewLoginTokenRepository(identityDAO),
		identitypg.NewUserRepository(identityDAO, emailProtector),
		defaultClock,
		accessURL,
	)
	if err != nil {
		return err
	}
	deliveryRepository := outboxpg.NewEventDeliveryRepository(outboxpg.NewDAO(db), payloadProtector, id.UUIDGenerator{})
	processor, err := application.NewEventProcessor(deliveryRepository, defaultClock,
		map[domain.EventType]application.EventHandler{domain.EventType(events.MagicLinkRequested): handler},
		application.EventProcessorOptions{
			LeaseDuration: time.Minute,
			HandleTimeout: 20 * time.Second,
			RetryDelay:    30 * time.Second,
			MaxAttempts:   5,
		},
	)
	if err != nil {
		return err
	}
	slog.Info("outbox worker started")
	err = polling.Run(ctx, processor, time.Second)
	slog.Info("outbox worker stopped")
	return err
}
