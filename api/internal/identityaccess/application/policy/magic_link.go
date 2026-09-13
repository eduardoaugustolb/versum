package policy

import (
	"time"

	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/application"
)

type MagicLinkPolicy struct {
	TTL                     time.Duration
	InvalidatePreviousLinks bool
	TokenRetention          time.Duration

	MaxDeliveryAttempts    int
	InitialDeliveryBackoff time.Duration
	MaxDeliveryBackoff     time.Duration
}

var DefaultMagicLinkPolicy = MagicLinkPolicy{
	TTL:                     time.Minute * 15,
	InvalidatePreviousLinks: false,
	TokenRetention:          time.Hour * 24,
	MaxDeliveryAttempts:     2,
	InitialDeliveryBackoff:  time.Minute,
	MaxDeliveryBackoff:      time.Hour,
}

func ValidateMagicLinkPolicy(policy MagicLinkPolicy) error {
	if policy.TTL <= 0 {
		return application.ErrInvalidMagicLinkPolicy
	}
	if policy.MaxDeliveryAttempts < 1 {
		return application.ErrInvalidMagicLinkPolicy
	}
	if policy.InitialDeliveryBackoff <= 0 {
		return application.ErrInvalidMagicLinkPolicy
	}
	if policy.MaxDeliveryBackoff <= 0 {
		return application.ErrInvalidMagicLinkPolicy
	}

	if policy.InitialDeliveryBackoff > policy.MaxDeliveryBackoff {
		return application.ErrInvalidMagicLinkPolicy
	}

	if policy.TokenRetention <= 0 {
		return application.ErrInvalidMagicLinkPolicy
	}

	return nil
}

func (p *MagicLinkPolicy) ApplyDefaults() {
	if p.TTL == 0 {
		p.TTL = DefaultMagicLinkPolicy.TTL
	}
	if p.MaxDeliveryAttempts == 0 {
		p.MaxDeliveryAttempts = DefaultMagicLinkPolicy.MaxDeliveryAttempts
	}
	if p.InitialDeliveryBackoff == 0 {
		p.InitialDeliveryBackoff = DefaultMagicLinkPolicy.InitialDeliveryBackoff
	}
	if p.MaxDeliveryBackoff == 0 {
		p.MaxDeliveryBackoff = DefaultMagicLinkPolicy.MaxDeliveryBackoff
	}
	if p.TokenRetention == 0 {
		p.TokenRetention = DefaultMagicLinkPolicy.TokenRetention
	}
}
