package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/nikhilsiwach28/MyCode.git/docker"
	"github.com/nikhilsiwach28/MyCode.git/errors"
	"github.com/nikhilsiwach28/MyCode.git/models"
	"github.com/nikhilsiwach28/MyCode.git/producer"
	"github.com/nikhilsiwach28/MyCode.git/redis"
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

func (c *Consumer) ConsumeMessages(sigterm chan os.Signal, producer *producer.Producer, redis *redis.RedisService) {
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

			// Validate and parse the message
			requestMsg, err := parseRequestMessage(msg)
			if err != nil {
				fmt.Printf("Error parsing request message: %v\n", err)
				continue consumeLoop
			}

			// Process the parsed message
			wg.Add(1)
			go processMessage(requestMsg, &wg, producer, redis)

		}
	}

	wg.Wait() // Wait for all processing to finish before returning
}

func parseRequestMessage(msg kafka.Message) (models.RequestMessage, error) {
	var requestMsg models.RequestMessage
	err := json.Unmarshal(msg.Value, &requestMsg)
	if err != nil {
		return models.RequestMessage{}, err
	}
	return requestMsg, nil
}

func processMessage(msg models.RequestMessage, wg *sync.WaitGroup, producer *producer.Producer, redis *redis.RedisService) {
	defer wg.Done()
	dockerClient, err := docker.NewDockerClient()
	if err != nil {
		fmt.Printf("Error creating Docker client: %v\n", err)
		return
	}
	if file, err := redis.Get(msg.ID); err != nil {
		fmt.Println("No SUbmission File Found in Redis for SubmissionId = ", msg.ID)
		//TODO Retry or return
	} else {
		output, err := dockerClient.RunContainer(file, msg.Language)
		if err != nil {
			fmt.Printf("Error running container: %v\n", err)
			return
		}
		// update Redis With Output
		outputKey := msg.ID + "_output"
		if err := redis.Set(outputKey, output); err != nil {
			fmt.Println("Error inserting Output for submissionID = ", msg.ID)
		}

		// Produce result message
		err = producer.ProduceMessage(models.ResponseMessage{Key: msg.ID, Value: outputKey}) // need to send codeId in Key
		if err != nil {
			fmt.Printf("Error producing message: %v\n", err)
			// Handle error...
		}
	}

}
