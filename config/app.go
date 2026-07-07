package config

type AppConfig struct {
	httpPort string
	envType  string

	DbConfig      *PostgresConfig
	TimeoutConfig *TimeoutConfig
	EsConfig      *EsConfig
}

func NewAppConfig() *AppConfig {
	httpPort := HTTPServerPort.MustGet()
	envType := EnvType.MustGet()

	return &AppConfig{
		httpPort:      httpPort,
		envType:       envType,
		DbConfig:      NewPostgresConfig(),
		TimeoutConfig: NewTimeoutConfig(),
		EsConfig:      NewESConfig(),
	}
}

func (c *AppConfig) HTTPPort() string {
	return c.httpPort
}

func (c *AppConfig) EnvType() string {
	return c.envType
}
