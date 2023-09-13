package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Db    Db
	Redis Redis
	Kafka Kafka
}

type Db struct {
	DbUrl string `env:"DB_URL"`
}

type Redis struct {
	RedisUrl string        `env:"REDIS_URL"`
	RedisExp time.Duration `env:"REDIS_EXP"`
}

type Kafka struct {
	KafkaUrl     string `env:"KAFKA_URL"`
	KafkaTopic   string `env:"KAFKA_TOPIC"`
	KafkaGroupID string `env:"KAFKA_GROUP_ID"`
}

func LoadConfig() *Config {
	cfgPath := os.Getenv("CONFIG_PATH")
	if cfgPath == "" {
		log.Fatalf("[!PANIC!] Config path environment variable is not set!")
	}

	if _, err := os.Stat(cfgPath); err != nil {
		log.Fatalf("[!Panic!] Config file is not exists!")
	}

	var cfg Config

	if err := cleanenv.ReadConfig(cfgPath, &cfg); err != nil {
		log.Fatalf("[!Panic!] %s", err)
	}

	return &cfg
}
