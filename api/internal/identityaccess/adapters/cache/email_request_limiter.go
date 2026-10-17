package cache

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	shared "github.com/eduardoaugustolb/versum/api/internal/cache"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/application"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/domain"
	"time"
)

type EmailRequestLimiter struct{ store shared.Cache }

var _ application.EmailRequestLimiter = (*EmailRequestLimiter)(nil)

func NewEmailRequestLimiter(store shared.Cache) *EmailRequestLimiter {
	return &EmailRequestLimiter{store}
}
func (l *EmailRequestLimiter) Allow(ctx context.Context, email domain.Email) (bool, error) {
	digest := sha256.Sum256([]byte(email.String()))
	return shared.AllowWithinWindow(ctx, l.store, "identityaccess:magic-link:email:"+hex.EncodeToString(digest[:]), time.Minute, 10)
}
