package ports

type LoginTokenHasher interface {
	Hash(token string) ([]byte, error)
}
