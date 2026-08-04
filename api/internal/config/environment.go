package config

type Environment string

const (
	DevelopmentEnvironment Environment = "development"
	ProductionEnvironment  Environment = "production"
	DefaultEnvironment                 = DevelopmentEnvironment
)

func parseEnvironment(env string) (Environment, error) {
	switch env {
	case "":
		return DefaultEnvironment, nil
	case "development":
		return DevelopmentEnvironment, nil
	case "production":
		return ProductionEnvironment, nil
	default:
		return "", ErrInvalidEnvironment
	}
}
