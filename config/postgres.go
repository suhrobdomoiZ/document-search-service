package config

import (
	"fmt"
	"net"
)

type PostgresConfig struct {
	dbName     string
	dbPassword string
	dbUser     string
	dbHost     string
	dbPort     string
}

func NewPostgresConfig() *PostgresConfig {
	dbName := DatabaseName.Get("default")
	dbHost := DatabaseHost.MustGet()
	dbUser := DatabaseUser.MustGet()
	dbPassword := DatabasePassword.MustGet()
	dbPort := DatabasePort.MustGet()

	return &PostgresConfig{
		dbName:     dbName,
		dbPassword: dbPassword,
		dbUser:     dbUser,
		dbHost:     dbHost,
		dbPort:     dbPort,
	}
}

func (p *PostgresConfig) DBName() string {
	return p.dbName
}

func (p *PostgresConfig) DBPassword() string {
	return p.dbPassword
}

func (p *PostgresConfig) DBUser() string {
	return p.dbUser
}

func (p *PostgresConfig) DBHost() string {
	return p.dbHost
}

func (p *PostgresConfig) DBPort() string {
	return p.dbPort
}

func (p *PostgresConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		p.dbUser,
		p.dbPassword,
		net.JoinHostPort(p.dbHost, p.dbPort),
		p.dbName,
	)
}
