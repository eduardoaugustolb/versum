package id

import (
	"uuid"
)

type UUIDGenerator struct{}

var _ Generator = UUIDGenerator{}

func (UUIDGenerator) Generate() string {
	return uuid.New().String()
}
