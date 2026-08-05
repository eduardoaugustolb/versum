package config

import "errors"

var (
	ErrPortNotNumeric = errors.New("port must be a number")
	ErrPortOutOfRange = errors.New("port must be between 1 and 65535")

	ErrInvalidEnvironment = errors.New("invalid environment")

	ErrDatabaseURLNotSet = errors.New("database url not set")
)
