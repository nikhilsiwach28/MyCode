package queue

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/nikhilsiwach28/MyCode.git/config"
	"gopkg.in/Shopify/sarama.v1"
)

// KafkaQueue represents a Kafka-based message queue.
type KafkaQueue struct {
	producer      sarama.AsyncProducer
	consumer      sarama.Consumer
	subscribedMap map[string]subscriptionChannels // Map to track subscribed topics
	mutex         sync.Mutex
}

// subscriptionChannels represents channels for messages and errors for a subscription.
type subscriptionChannels struct {
	messages <-chan *sarama.ConsumerMessage
	errors   <-chan error
}

// NewKafkaQueue creates a new Kafka-based message queue.
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
		producer:      producer,
		consumer:      consumer,
		subscribedMap: make(map[string]subscriptionChannels), // Initialize map
	}, nil
}

// Subscribe subscribes to a Kafka topic and returns channels for receiving messages and errors.
func (k *KafkaQueue) Subscribe(topic string) (<-chan *sarama.ConsumerMessage, <-chan error) {
	k.mutex.Lock()
	defer k.mutex.Unlock()

	// Check if already subscribed
	if subscription, ok := k.subscribedMap[topic]; ok {
		return subscription.messages, subscription.errors
	}

	// Create new subscription
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

	// Update subscribed map
	k.subscribedMap[topic] = subscriptionChannels{
		messages: messages,
		errors:   errors,
	}

	return messages, errors
}

// SendMessage sends a message to a Kafka topic.
func (k *KafkaQueue) SendMessage(topic, message string) error {
	producerMessage := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	k.producer.Input() <- producerMessage

	return nil
}

// Close closes the Kafka queue.
func (k *KafkaQueue) Close() {
	if err := k.consumer.Close(); err != nil {
		fmt.Printf("Error closing consumer: %s\n", err)
	}

	if err := k.producer.Close(); err != nil {
		fmt.Printf("Error closing producer: %s\n", err)
	}
}

// InitKafkaQueue initializes a Kafka-based message queue.
func InitKafkaQueue(cfg config.KafkaConfig) *KafkaQueue {
	queue, err := NewKafkaQueue(cfg.Brokers)
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
