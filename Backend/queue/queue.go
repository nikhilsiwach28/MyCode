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

func (k *KafkaQueue) Produce(topic string, message string) error {
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

func InitQueue() {
	brokers := []string{"localhost:9092"}
	queue, err := NewKafkaQueue(brokers)
	if err != nil {
		panic(err)
	}
	defer queue.Close()

	topic := "executedCode"
	// produced_topic := "produced-topic"
	// Subscribe to topic
	messages, errors := queue.Subscribe(topic)

	// Handle signals for graceful shutdown
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	// Consumer loop
	go func() {
		for {
			select {
			case msg := <-messages:
				fmt.Printf("Received Executed Code: %s\n", msg.Value)
			case err := <-errors:
				fmt.Printf("Error in executed code: %s\n", err.Error())
			case <-sigterm:
				fmt.Println("Received shutdown signal. Shutting down consumer.")
				return
			}
		}
	}()

	// Producer loop
	// go func() {
	// 	fmt.Println("producing messages to Queue every 2 second")
	// 	for {
	// 		err := queue.Produce("toBeExecuted", "Hello, Kafka 2!")
	// 		if err != nil {
	// 			fmt.Printf("Error producing message: %s\n", err)
	// 		}
	// 		time.Sleep(5 * time.Second)
	// 	}
	// }()

	// Wait for shutdown signal
	<-sigterm
}
