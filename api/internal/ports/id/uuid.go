package id

type UUID string

type IDGenerator interface {
	Generate() UUID
}

func (id UUID) String() string {
	return string(id)
}
