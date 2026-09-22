package application

type LoginTokenHasher interface {
	Hash(token string) ([]byte, error)
}
