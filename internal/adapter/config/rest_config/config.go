package rest_config

import "time"

type Config struct {
	Env    string `env:"ENV"`
	Db     DbConfig
	Server ServerConfig
	Redis  RedisConfig
}

type DbConfig struct {
	Driver string `env:"DB_DRIVER"`
	Url    string `env:"DB_URL"`
}

type ServerConfig struct {
	Address     string        `env:"SERVER_ADDRESS"`
	Timeout     time.Duration `env:"SERVER_TIMEOUT"`
	IdleTimeout time.Duration `env:"SERVER_IDLE_TIMEOUT"`
}

type RedisConfig struct {
	Address  string        `env:"REDIS_ADDRESS"`
	Password string        `env:"REDIS_PASSWORD"`
	DB       int           `env:"REDIS_DB"`
	Timeout  time.Duration `env:"REDIS_PING_TIMEOUT"`
}
