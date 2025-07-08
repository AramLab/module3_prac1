package config

import "time"

const EnvPath = "local.env"

type AppConfig struct {
	LogLevel         string
	Rest             Rest
	RepositoryConfig RepositoryConfig
}

type Rest struct {
	ListenAddress string        `envconfig:"PORT" required:"true"`
	WriteTimeout  time.Duration `envconfig:"WRITE_TIMEOUT" required:"true"`
	ServerName    string        `envconfig:"SERVER_NAME" required:"true"`
}

type RepositoryConfig struct {
	Capacity int `envconfig:"CAPACITY" required:"true"`
}
