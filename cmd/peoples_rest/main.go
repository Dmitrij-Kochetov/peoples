package main

import (
	"fmt"
	"github.com/Dmitrij-Kochetov/peoples/internal/adapter/config"
	"github.com/Dmitrij-Kochetov/peoples/internal/adapter/config/rest_config"
)

func main() {
	var cfg rest_config.Config
	cfg = config.LoadConfig(cfg)
	fmt.Printf("%v\n", cfg)
}
