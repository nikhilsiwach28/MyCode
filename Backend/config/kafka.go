package config

type KafkaConfig struct {
	Brokers []string
}

func NewKafkaConfig() KafkaConfig {
	return KafkaConfig{
		Brokers: []string{GetEnvWithDefault("KAFKA_BROKERS", "localhost:9092")},
	}
}
