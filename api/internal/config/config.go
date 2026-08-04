package config

type Config struct {
	Address     string
	Environment Environment
}

func Load(lookup func(string) string) (Config, error) {
	cfg := Config{}

	port, err := parsePort(lookup("PORT"))
	if err != nil {
		return Config{}, err
	}
	cfg.Address = ":" + port

	environment, err := parseEnvironment(lookup("ENVIRONMENT"))
	if err != nil {
		return Config{}, err
	}
	cfg.Environment = environment

	return cfg, nil
}
