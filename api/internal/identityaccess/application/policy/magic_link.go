package policy

import (
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/application"
	"time"
)

// MagicLinkPolicy configures issuance. Delivery and retention have no worker yet.
type MagicLinkPolicy struct{ TTL time.Duration }

var DefaultMagicLinkPolicy = MagicLinkPolicy{TTL: 15 * time.Minute}

func ValidateMagicLinkPolicy(p MagicLinkPolicy) error {
	if p.TTL <= 0 {
		return application.ErrInvalidMagicLinkTTL
	}
	return nil
}
func (p *MagicLinkPolicy) ApplyDefaults() {
	if p.TTL == 0 {
		p.TTL = DefaultMagicLinkPolicy.TTL
	}
}
