package producer

import (
	"context"
	"fmt"

	"github.com/nikhilsiwach28/MyCode.git/errors"
	"github.com/nikhilsiwach28/MyCode.git/models"
	"github.com/nikhilsiwach28/MyCode.git/utils"
	"github.com/segmentio/kafka-go"
)

type Producer struct {
	*kafka.Writer
}

func NewProducerFromEnv() (*Producer, error) {
	broker := utils.GetEnvValueWithDefault("KAFKA_BROKER", "localhost:9092")
	topic := utils.GetEnvValueWithDefault("KAFKA_TOPIC_OUT", "executedCode")
	if broker == "" || topic == "" {
		return nil, &errors.WorkerError{Message: "Kafka broker or topic not provided"}
	}
	return NewProducer(broker, topic), nil
}

func NewProducer(broker, topic string) *Producer {
	return &Producer{
		Writer: kafka.NewWriter(kafka.WriterConfig{
			Brokers:  []string{broker},
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		}),
	}
}

func (p *Producer) ProduceMessage(msg models.ResponseMessage) error {
	fmt.Println("Producing Executed Code")
	return p.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(msg.Key),
		Value: []byte(msg.Value),
	})
}
