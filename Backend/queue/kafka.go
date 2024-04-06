package queue

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"gopkg.in/Shopify/sarama.v1"
)

type KafkaQueue struct {
	producer sarama.AsyncProducer
	consumer sarama.Consumer
}

func NewKafkaQueue(brokers []string) (*KafkaQueue, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	// Create new consumer
	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		return nil, err
	}

	// Create new producer
	producer, err := sarama.NewAsyncProducer(brokers, nil)
	if err != nil {
		return nil, err
	}

	return &KafkaQueue{
		producer: producer,
		consumer: consumer,
	}, nil
}

func (k *KafkaQueue) Subscribe(topic string) (<-chan *sarama.ConsumerMessage, <-chan error) {
	partitionConsumer, err := k.consumer.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		panic(err)
	}

	messages := make(chan *sarama.ConsumerMessage)
	errors := make(chan error)

	// Consumer loop
	go func() {
		for {
			select {
			case msg := <-partitionConsumer.Messages():
				messages <- msg
			case err := <-partitionConsumer.Errors():
				errors <- err
			}
		}
	}()

	return messages, errors
}

func (k *KafkaQueue) SendMessage(topic, message string) error {
	producerMessage := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	k.producer.Input() <- producerMessage

	return nil
}

func (k *KafkaQueue) Close() {
	if err := k.consumer.Close(); err != nil {
		fmt.Printf("Error closing consumer: %s\n", err)
	}

	if err := k.producer.Close(); err != nil {
		fmt.Printf("Error closing producer: %s\n", err)
	}
}

func InitQueue(brokers []string) *KafkaQueue {
	queue, err := NewKafkaQueue(brokers)
	if err != nil {
		panic(err)
	}

	// Handle signals for graceful shutdown
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	// Wait for shutdown signal
	<-sigterm

	return queue
}
