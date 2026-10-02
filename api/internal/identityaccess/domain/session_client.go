package domain

import (
	"net/netip"
	"strings"
)

// SessionClient is an immutable value object describing the observed client.
// Its metadata is a risk signal, not proof of possession of a credential.
// Session constructs it from primitive inputs.
type SessionClient struct {
	ipAddress string
	userAgent string
}

func newSessionClient(ipAddress, userAgent string) (SessionClient, error) {
	if ipAddress != "" {
		addr, err := netip.ParseAddr(ipAddress)
		if err != nil || addr.Zone() != "" {
			return SessionClient{}, ErrInvalidSessionIPAddress
		}
		ipAddress = addr.Unmap().String()
	}
	if len(userAgent) > 1024 || strings.ContainsAny(userAgent, "\r\n\x00") {
		return SessionClient{}, ErrInvalidSessionUserAgent
	}
	return SessionClient{ipAddress: ipAddress, userAgent: userAgent}, nil
}

func (c SessionClient) IPAddress() string { return c.ipAddress }
func (c SessionClient) UserAgent() string { return c.userAgent }
