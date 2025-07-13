package config

import "time"

const EnvPath = "local.env"

type AppConfig struct {
	LogLevel   string `envconfig:"LOG_LEVEL" required:"true"`
	Rest       Rest
	PostgreSQL PostgreSQL
}

type Rest struct {
	ListenAddress string        `envconfig:"PORT" required:"true"`
	WriteTimeout  time.Duration `envconfig:"WRITE_TIMEOUT" required:"true"`
	ServerName    string        `envconfig:"SERVER_NAME" required:"true"`
}

type PostgreSQL struct {
	Host     string `envconfig:"POSTGRES_HOST" required:"true"`
	Port     int    `envconfig:"POSTGRES_PORT" required:"true"`
	Name     string `envconfig:"POSTGRES_NAME" required:"true"`
	User     string `envconfig:"POSTGRES_USER" required:"true"`
	Password string `envconfig:"POSTGRES_PASSWORD" required:"true"`
	SSLMode  string `envconfig:"POSTGRES_SSL_MODE" required:"true"`

	PoolMaxConns        int           `envconfig:"POSTGRES_POOL_MAX_CONNS" default:"5"`
	PoolMaxConnLifetime time.Duration `envconfig:"POSTGRES_POOL_MAX_CONN_LIFETIME" default:"180s"`
	PoolMaxConnIdleTime time.Duration `envconfig:"POSTGRES_POOL_MAX_CONN_IDLE_TIME" default:"100s"`

	AutoMigrate bool `envconfig:"POSTGRES_AUTO_MIGRATE" default:"false"`
}
