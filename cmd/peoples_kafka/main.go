package main

import (
	"fmt"
	"github.com/Dmitrij-Kochetov/peoples/internal/adapter/config"
	"github.com/Dmitrij-Kochetov/peoples/internal/adapter/config/kafka_config"
	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	consumer *kafka.Reader
}

func NewConsumer() (*Consumer, error) {
	return nil, nil
}

type Producer struct {
	producer *kafka.Writer
}

func NewProducer() (*Producer, error) {
	return nil, nil
}

func main() {
	var cfg kafka_config.Config
	cfg = config.LoadConfig(cfg)
	fmt.Printf("%v\n", cfg)
}
