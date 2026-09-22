package application

type LoginTokenGenerator interface {
	GenerateLoginToken() (string, error)
}
