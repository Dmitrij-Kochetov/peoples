package kafka_config

type Config struct {
	Env   string `env:"ENV"`
	Kafka KafkaConfig
	DB    DbConfig
}

type KafkaConfig struct {
	KafkaURL      string `env:"KAFKA_URL"`
	ConsumerTopic string `env:"KAFKA_CONSUMER_TOPIC"`
	ConsumerGroup string `env:"KAFKA_CONSUMER_GROUP"`
	ProducerTopic string `env:"KAFKA_PRODUCER_TOPIC"`
	KafkaTimeout  int    `env:"KAFKA_TIMEOUT"`
}

type DbConfig struct {
	Driver string `env:"DB_DRIVER"`
	DbURL  string `env:"DB_URL"`
}
