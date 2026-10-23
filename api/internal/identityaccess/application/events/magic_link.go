// Package events owns identity and access integration event contracts.
package events

const MagicLinkRequested = "identityaccess.magic_link_requested"

type MagicLinkRequestedPayload struct {
	LoginTokenID string `json:"login_token_id"`
	Token        string `json:"token"`
}
