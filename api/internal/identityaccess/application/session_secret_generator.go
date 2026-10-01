package application

// SessionSecretGenerator generates unpredictable bearer credentials.
type SessionSecretGenerator interface {
	GenerateSessionSecret() (string, error)
}
