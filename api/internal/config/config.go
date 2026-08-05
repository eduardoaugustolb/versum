package config

type Config struct {
	Address     string
	Environment Environment
	DatabaseURL string
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

	databaseURL := lookup(DefaultDatabaseURLKey)
	if databaseURL == "" {
		databaseURL = DefaultDatabaseURL
	}
	cfg.DatabaseURL = databaseURL

	return cfg, nil
}
