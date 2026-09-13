package ports

import "context"

type ProtectedPayload struct {
	Ciphertext []byte
	KeyVersion int
}

type PayloadProtector interface {
	Protect(ctx context.Context, payload []byte) (ProtectedPayload, error)
	Unprotect(ctx context.Context, protected ProtectedPayload) ([]byte, error)
}
