package ports

type LoginTokenGenerator interface {
	GenerateLoginToken() (string, error)
}
