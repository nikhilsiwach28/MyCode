package consumer

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/nikhilsiwach28/MyCode.git/docker"
	"github.com/nikhilsiwach28/MyCode.git/errors"
	"github.com/nikhilsiwach28/MyCode.git/models"
	"github.com/nikhilsiwach28/MyCode.git/producer"
	"github.com/nikhilsiwach28/MyCode.git/utils"
	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	*kafka.Reader
}

func NewConsumerFromEnv() (*Consumer, error) {
	broker := utils.GetEnvValueWithDefault("KAFKA_BROKER", "localhost:9092")
	topic := utils.GetEnvValueWithDefault("KAFKA_TOPIC_IN", "toBeExecuted")
	if broker == "" || topic == "" {
		return nil, &errors.WorkerError{Message: "Kafka broker or topic not provided"}
	}
	return NewConsumer(broker, topic), nil
}

func NewConsumer(broker, topic string) *Consumer {
	return &Consumer{
		Reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:   []string{broker},
			Topic:     topic,
			Partition: 0,
			MinBytes:  10e3, // 10KB
			MaxBytes:  10e6, // 10MB
		}),
	}
}

func (c *Consumer) ConsumeMessages(sigterm chan os.Signal, producer *producer.Producer) {
	var wg sync.WaitGroup

consumeLoop:
	for {
		select {
		case <-sigterm:
			// Received SIGINT or SIGTERM, terminate loop
			fmt.Println("Received interrupt signal. Exiting...")
			break consumeLoop

		default:
			// Continue consuming messages
			msg, err := c.ReadMessage(context.Background())
			if err != nil {
				fmt.Printf("Error reading message: %v\n", err)
				continue consumeLoop
			}
			wg.Add(1)
			go processMessage(models.RequestMessage{Code: string(msg.Key), Language:string(msg.Value)}, &wg, producer)
		}
	}

	wg.Wait() // Wait for all processing to finish before returning
}


func processMessage(msg models.RequestMessage, wg *sync.WaitGroup, producer *producer.Producer) {
	defer wg.Done()
	dockerClient, err := docker.NewDockerClient()
	if err != nil {
		fmt.Printf("Error creating Docker client: %v\n", err)
		return
	}

	// Execute code in Docker container
	output, err := dockerClient.RunContainer(msg.Code, msg.Language)
	if err != nil {
		fmt.Printf("Error running container: %v\n", err)
		return
	}

	// Produce result message
	err = producer.ProduceMessage(models.ResponseMessage{Key: "randomId", Value: output})   // need to send codeId in Key
	if err != nil {
		fmt.Printf("Error producing message: %v\n", err)
		// Handle error...
	}
}
