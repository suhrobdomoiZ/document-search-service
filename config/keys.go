package config

import "github.com/suhrobdomoiZ/document-search-service/pkg/configkey"

const (
	HTTPServerPort configkey.Key = "HTTP_PORT"
	EnvType        configkey.Key = "ENV_TYPE"

	DatabaseName     configkey.Key = "DATABASE_NAME"
	DatabaseHost     configkey.Key = "DATABASE_HOST"
	DatabasePort     configkey.Key = "DATABASE_PORT"
	DatabaseUser     configkey.Key = "DATABASE_USER"
	DatabasePassword configkey.Key = "DATABASE_PASSWORD"

	ServerReadTimeout  configkey.Key = "READ_TIMEOUT"
	ServerWriteTimeout configkey.Key = "WRITE_TIMEOUT"
	ServerIdleTimeout  configkey.Key = "IDLE_TIMEOUT"
	ShutdownCtxTimeout configkey.Key = "SHUTDOWN_TIMEOUT"

	EsPort       configkey.Key = "ES_PORT"
	EsSearchSize configkey.Key = "ES_SEARCH_SIZE"

	DataPath configkey.Key = "DATA_PATH"
)
