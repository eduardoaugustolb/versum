package config

import (
	"strconv"
	"strings"
)

const DefaultPort = "8080"

func parsePort(port string) (string, error) {
	port = strings.TrimPrefix(port, ":")
	if port == "" {
		return DefaultPort, nil
	}

	n, err := strconv.Atoi(port)
	if err != nil {
		return "", ErrPortNotNumeric
	}
	if n < 1 || n > 65535 {
		return "", ErrPortOutOfRange
	}

	return port, nil
}
