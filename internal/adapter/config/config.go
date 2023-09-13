package config

import (
	"github.com/Dmitrij-Kochetov/peoples/internal/adapter/config/graph_config"
	"github.com/Dmitrij-Kochetov/peoples/internal/adapter/config/kafka_config"
	"github.com/Dmitrij-Kochetov/peoples/internal/adapter/config/rest_config"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type IConfig interface {
	graph_config.Config | kafka_config.Config | rest_config.Config
}

func LoadConfig[C IConfig](cfg C) C {
	cfgPath := os.Getenv("CONFIG_PATH")
	if cfgPath == "" {
		log.Fatalf("[!Panic!] Config path variable is not set!")
	}

	if _, err := os.Stat(cfgPath); err != nil {
		log.Fatalf("[!Panic!] Config file is not exists!")
	}

	if err := cleanenv.ReadConfig(cfgPath, &cfg); err != nil {
		log.Fatalf("[!Panic!] %s", err)
	}
	return cfg
}
