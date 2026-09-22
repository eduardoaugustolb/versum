package id

// Generator produces identifiers for application entities.
type Generator interface {
	Generate() string
}
