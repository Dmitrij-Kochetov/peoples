package kafka

import "github.com/confluentinc/confluent-kafka-go/kafka"

type Consumer struct {
	Consumer  *kafka.Consumer
	TimeoutMs int
}

func NewKafkaConsumer(host, topic, groupID string, timeoutMs int) (*Consumer, error) {
	cfg := kafka.ConfigMap{
		"bootstrap.servers":  host,
		"group.id":           groupID,
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": false,
	}

	client, err := kafka.NewConsumer(&cfg)
	if err != nil {
		return nil, err
	}
	if err := client.Subscribe(topic, nil); err != nil {
		return nil, err
	}
	return &Consumer{Consumer: client, TimeoutMs: timeoutMs}, nil
}

type Producer struct {
	Producer *kafka.Producer
	Topic    string
}

func NewKafkaProducer(host, topic string) (*Producer, error) {
	cfg := kafka.ConfigMap{
		"bootstrap.servers": host,
	}
	client, err := kafka.NewProducer(&cfg)
	if err != nil {
		return nil, err
	}
	return &Producer{
		Producer: client,
		Topic:    topic,
	}, nil
}
