package main

import (
	"fmt"
	"github.com/Dmitrij-Kochetov/peoples/internal/adapter/config"
	"github.com/Dmitrij-Kochetov/peoples/internal/adapter/config/graph_config"
)

func main() {
	var cfg graph_config.Config
	cfg = config.LoadConfig(cfg)
	fmt.Printf("%v\n", cfg)
}
