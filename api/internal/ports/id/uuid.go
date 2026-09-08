package id

type UUID string

type IDGenerator interface {
	Generate() UUID
}
